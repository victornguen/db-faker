package dbutils

import (
	"database/sql"
	"fmt"
	"github.com/victornguen/db-faker/datagen"
	"strings"
)

const (
	getTableNamesQuery = `
		SELECT table_name 
		FROM information_schema.tables 
		WHERE table_schema = 'public'
	`

	getPrimaryKeyColumnsQuery = `
		SELECT kcu.column_name
		FROM information_schema.key_column_usage kcu
		JOIN information_schema.table_constraints tc
			ON tc.constraint_name = kcu.constraint_name
		WHERE tc.table_name = $1
			AND tc.constraint_type = 'PRIMARY KEY'
	`

	getTableDependenciesQuery = `
		SELECT DISTINCT kcu2.table_name
		FROM information_schema.referential_constraints rc
		JOIN information_schema.key_column_usage kcu
			ON rc.constraint_name = kcu.constraint_name
		JOIN information_schema.key_column_usage kcu2
			ON rc.unique_constraint_name = kcu2.constraint_name
		WHERE kcu.table_name = $1
	`

	getPrimaryKeyColumnQuery = `
		SELECT kcu.column_name
		FROM information_schema.key_column_usage kcu
		JOIN information_schema.table_constraints tc
			ON tc.constraint_name = kcu.constraint_name
		WHERE tc.table_name = $1
			AND tc.constraint_type = 'PRIMARY KEY'
		LIMIT 1
	`
)

func ApplyRulesToTables(tables *[]Table, rules datagen.TablesRules) error {
	for i, table := range *tables {
		if rule, ok := rules.Rules[table.Name]; ok {
			table.RowNum = rule.RowNum
			for colName, rule := range rule.Rules {
				genFunc, err := datagen.RuleToGeneratorFunc(rule)
				if err != nil {
					return fmt.Errorf("error generating function for rule %s: %v", rule, err)
				}
				col := table.Columns[colName]
				col.DataGen = genFunc
				_, present := table.Columns[colName]
				if present {
					table.Columns[colName] = col
				}
			}
			(*tables)[i] = table
		}
	}
	return nil
}

func GetTablesWithDependencies(db *sql.DB) ([]Table, error) {
	tables := make([]Table, 0)

	// Get all tables
	rows, err := db.Query(getTableNamesQuery)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var tableName string
		if err := rows.Scan(&tableName); err != nil {
			return nil, err
		}

		// Get primary keys
		pkCols, err := getPrimaryKeyColumns(db, tableName)
		if err != nil {
			return nil, err
		}

		// Get columns
		columns, err := GetColumns(db, tableName)
		if err != nil {
			return nil, err
		}

		columnsMap := make(map[string]Column)
		for _, col := range columns {
			columnsMap[col.Name] = col
		}

		// Get dependencies
		deps, err := getTableDependencies(db, tableName)
		if err != nil {
			return nil, err
		}

		tables = append(tables, Table{
			Name:        tableName,
			Columns:     columnsMap,
			DependsOn:   deps,
			PrimaryKeys: pkCols,
			RowNum:      0,
			Rules:       make(map[string]func() string),
		})
	}

	return tables, nil
}

func getPrimaryKeyColumns(db *sql.DB, tableName string) (map[string]bool, error) {
	rows, err := db.Query(getPrimaryKeyColumnsQuery, tableName)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	pkCols := make(map[string]bool)
	for rows.Next() {
		var colName string
		if err := rows.Scan(&colName); err != nil {
			return nil, err
		}
		pkCols[colName] = true
	}
	return pkCols, nil
}

func getTableDependencies(db *sql.DB, tableName string) ([]string, error) {
	rows, err := db.Query(getTableDependenciesQuery, tableName)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	deps := make([]string, 0)
	for rows.Next() {
		var dep string
		if err := rows.Scan(&dep); err != nil {
			return nil, err
		}
		deps = append(deps, dep)
	}

	return deps, nil
}

func TopologicalSort(tables []Table) []Table {
	var sorted []Table
	visited := make(map[string]bool)

	var visit func(table Table)
	visit = func(table Table) {
		if !visited[table.Name] {
			visited[table.Name] = true
			for _, dep := range table.DependsOn {
				for _, t := range tables {
					if t.Name == dep {
						visit(t)
					}
				}
			}
			sorted = append(sorted, table)
		}
	}

	for _, table := range tables {
		visit(table)
	}

	return sorted
}

func getPrimaryKeyColumn(db *sql.DB, tableName string) (string, error) {
	var colName string
	err := db.QueryRow(getPrimaryKeyColumnQuery, tableName).Scan(&colName)
	if err != nil {
		return "", fmt.Errorf("error getting primary key for %s: %v", tableName, err)
	}
	return colName, nil
}

func GenerateAndInsertData(db *sql.DB, table Table) error {

	// Filter out primary key columns
	var filteredColumns []Column
	for _, col := range table.Columns {
		switch col.DataType.(type) {
		case Serial, BigSerial, SmallSerial, TsVector, TsQuery:
			continue
		}
		if !table.PrimaryKeys[col.Name] {
			filteredColumns = append(filteredColumns, col)
		}
	}

	// Prepare column names and placeholders for non-PK columns
	columns := make([]string, 0, len(filteredColumns))
	placeholders := make([]string, 0, len(filteredColumns))
	for i, col := range filteredColumns {
		columns = append(columns, col.Name)
		placeholders = append(placeholders, fmt.Sprintf("$%d", i+1))
	}

	query := fmt.Sprintf("INSERT INTO %s (%s) VALUES (%s)",
		table.Name,
		strings.Join(columns, ", "),
		strings.Join(placeholders, ", "),
	)

	stmt, err := db.Prepare(query)
	if err != nil {
		return err
	}
	defer stmt.Close()

	for i := 0; i < table.RowNum; i++ {
		values := make([]interface{}, len(filteredColumns))
		for j, col := range filteredColumns {
			if col.IsForeignKey && col.RefTable == "" {
				return fmt.Errorf("no reference table for %s.%s", table.Name, col.Name)
			} else if col.IsForeignKey {
				pkCol, err := getPrimaryKeyColumn(db, col.RefTable)
				if err != nil {
					return fmt.Errorf("table %s foreign key error: %v", table.Name, err)
				}

				var refID interface{}
				// Get random ID from referenced table
				err = db.QueryRow(
					fmt.Sprintf("SELECT %s FROM %s ORDER BY RANDOM() LIMIT 1", pkCol, col.RefTable),
				).Scan(&refID)

				if err != nil {
					return fmt.Errorf("no reference data found in %s for %s.%s: %v",
						col.RefTable, table.Name, col.Name, err)
				}
				values[j] = refID
			} else {
				values[j] = col.DataGen()
			}
		}

		_, err := stmt.Exec(values...)
		if err != nil {
			fmt.Printf("Error inserting row %d into %s: %v\n", i, table.Name, err)
		}
	}

	fmt.Printf("Inserted %d rows into %s\n", table.RowNum, table.Name)
	return nil
}
