package dbutils

import (
	"database/sql"
	"fmt"
	"github.com/go-faker/faker/v4"
	"math/rand"
	"strings"
	"time"
)

const (
	getTableNamesQuery = `
		SELECT table_name 
		FROM information_schema.tables 
		WHERE table_schema = 'public'
	`

	getColumnsQuery = `
		SELECT 
			c.column_name,
			c.data_type,
			COALESCE((
				SELECT TRUE 
				FROM information_schema.key_column_usage kcu
				JOIN information_schema.referential_constraints rc
					ON kcu.constraint_name = rc.constraint_name
				WHERE kcu.table_name = c.table_name
					AND kcu.column_name = c.column_name
				LIMIT 1
			), FALSE) AS is_foreign_key,
			COALESCE((
				SELECT kcu2.table_name
				FROM information_schema.referential_constraints rc
				JOIN information_schema.key_column_usage kcu
					ON rc.constraint_name = kcu.constraint_name
				JOIN information_schema.key_column_usage kcu2
					ON rc.unique_constraint_name = kcu2.constraint_name
				WHERE kcu.table_name = c.table_name
					AND kcu.column_name = c.column_name
				LIMIT 1
			), '') AS ref_table
		FROM information_schema.columns c
		WHERE table_name = $1
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
		columns, err := getColumns(db, tableName)
		if err != nil {
			return nil, err
		}

		// Get dependencies
		deps, err := getTableDependencies(db, tableName)
		if err != nil {
			return nil, err
		}

		tables = append(tables, Table{
			Name:        tableName,
			Columns:     columns,
			DependsOn:   deps,
			PrimaryKeys: pkCols,
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

func getColumns(db *sql.DB, tableName string) ([]Column, error) {
	rows, err := db.Query(getColumnsQuery, tableName)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	columns := make([]Column, 0)
	for rows.Next() {
		var col Column
		var dataTypeStr string
		if err := rows.Scan(&col.Name, &dataTypeStr, &col.IsForeignKey, &col.RefTable); err != nil {
			return nil, err
		}
		dataType, err := StringToDataType(dataTypeStr)
		if err != nil {
			return nil, err
		}
		col.DataType = dataType
		columns = append(columns, col)
	}

	return columns, nil
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

func GenerateAndInsertData(db *sql.DB, table Table, count int) error {

	// Filter out primary key columns
	var filteredColumns []Column
	for _, col := range table.Columns {
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

	for i := 0; i < count; i++ {
		values := make([]interface{}, len(filteredColumns))
		for j, col := range filteredColumns {
			if col.IsForeignKey && col.RefTable == "" {
				return fmt.Errorf("no reference table for %s.%s", table.Name, col.Name)
			}
			if col.IsForeignKey {
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
				switch col.DataType.(type) {
				case Int:
					values[j] = rand.Intn(1000)
				case Numeric:
					values[j] = rand.Float64() * 1000
				case VarChar:
					if strings.Contains(col.Name, "email") {
						values[j] = faker.Email()
					} else if strings.Contains(col.Name, "username") {
						values[j] = faker.Username()
					} else {
						values[j] = faker.Word()
					}
				case Text:
					values[j] = faker.Sentence()
				case Boolean:
					values[j] = rand.Intn(2) == 1
				case TimeStamp:
					values[j] = time.Now().Add(time.Duration(rand.Intn(100000)) * time.Hour)
				case Date:
					values[j] = time.Now().AddDate(0, 0, rand.Intn(100))
				default:
					return fmt.Errorf("unsupported data type: %T", col.DataType)
				}
			}
		}

		_, err := stmt.Exec(values...)
		if err != nil {
			fmt.Printf("Error inserting row %d into %s: %v\n", i, table.Name, err)
		}
	}

	fmt.Printf("Inserted %d rows into %s\n", count, table.Name)
	return nil
}
