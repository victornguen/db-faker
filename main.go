package main

import (
	"database/sql"
	"fmt"
	"github.com/go-faker/faker/v4"
	_ "github.com/lib/pq"
	"log"
	"math/rand"
	_ "sort"
	"strings"
	"time"
)

type Table struct {
	Name      string
	Columns   []Column
	DependsOn []string
}

type Column struct {
	Name         string
	DataType     string
	IsForeignKey bool
	RefTable     string
}

type TableDependency struct {
	TableName    string
	Dependencies []string
}

const (
	host     = "localhost"
	port     = 5433
	user     = "postgres"
	password = "postgres"
	dbname   = "portal"
)

func main() {
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)

	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	tables, err := getTablesWithDependencies(db)
	if err != nil {
		log.Fatal(err)
	}

	sortedTables := topologicalSort(tables)

	println("Tables in topological order:")
	for _, table := range sortedTables {
		fmt.Printf("%s\n", table)
	}

	for _, table := range sortedTables {
		err := generateAndInsertData(db, table, 10) // Generate 10 rows per table
		if err != nil {
			log.Printf("Error inserting data into %s: %v", table.Name, err)
		}
	}
}

func getTablesWithDependencies(db *sql.DB) ([]Table, error) {
	tables := make([]Table, 0)

	// Get all tables
	rows, err := db.Query(`
		SELECT table_name 
		FROM information_schema.tables 
		WHERE table_schema = 'public'
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var tableName string
		if err := rows.Scan(&tableName); err != nil {
			return nil, err
		}

		// Get columns for table
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
			Name:      tableName,
			Columns:   columns,
			DependsOn: deps,
		})
	}

	return tables, nil
}

func getColumns(db *sql.DB, tableName string) ([]Column, error) {
	query := `
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

	rows, err := db.Query(query, tableName)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	columns := make([]Column, 0)
	for rows.Next() {
		var col Column
		if err := rows.Scan(&col.Name, &col.DataType, &col.IsForeignKey, &col.RefTable); err != nil {
			return nil, err
		}
		columns = append(columns, col)
	}

	return columns, nil
}

func getTableDependencies(db *sql.DB, tableName string) ([]string, error) {
	query := `
		SELECT DISTINCT kcu2.table_name
		FROM information_schema.referential_constraints rc
		JOIN information_schema.key_column_usage kcu
			ON rc.constraint_name = kcu.constraint_name
		JOIN information_schema.key_column_usage kcu2
			ON rc.unique_constraint_name = kcu2.constraint_name
		WHERE kcu.table_name = $1
	`

	rows, err := db.Query(query, tableName)
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

func topologicalSort(tables []Table) []Table {
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
	query := `
		SELECT kcu.column_name
		FROM information_schema.key_column_usage kcu
		JOIN information_schema.table_constraints tc
			ON tc.constraint_name = kcu.constraint_name
		WHERE tc.table_name = $1
			AND tc.constraint_type = 'PRIMARY KEY'
		LIMIT 1
	`
	var colName string
	err := db.QueryRow(query, tableName).Scan(&colName)
	if err != nil {
		return "", fmt.Errorf("error getting primary key for %s: %v", tableName, err)
	}
	return colName, nil
}

func generateAndInsertData(db *sql.DB, table Table, count int) error {
	columns := make([]string, 0, len(table.Columns))
	placeholders := make([]string, 0, len(table.Columns))
	for i, col := range table.Columns {
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
		values := make([]interface{}, len(table.Columns))
		for j, col := range table.Columns {
			switch {
			case col.IsForeignKey:
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
			case col.DataType == "integer":
				values[j] = rand.Intn(1000)
			case col.DataType == "numeric":
				values[j] = rand.Float64() * 1000
			case col.DataType == "character varying":
				values[j] = faker.Word()
			case col.DataType == "text":
				values[j] = faker.Sentence()
			case col.DataType == "boolean":
				values[j] = rand.Intn(2) == 1
			case col.DataType == "timestamp with time zone":
				values[j] = time.Now().Add(time.Duration(rand.Intn(100000)) * time.Hour)
			case col.DataType == "timestamp without time zone":
				values[j] = time.Now().Add(time.Duration(rand.Intn(100000)) * time.Hour)
			case col.DataType == "date":
				values[j] = time.Now().AddDate(0, 0, rand.Intn(100))
			default:
				values[j] = faker.Word()
			}
		}

		_, err := stmt.Exec(values...)
		if err != nil {
			return fmt.Errorf("execution error: %v", err)
		}
	}

	fmt.Printf("Inserted %d rows into %s\n", count, table.Name)
	return nil
}
