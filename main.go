package main

import (
	"database/sql"
	"db_faker/dbutils"
	"fmt"
	_ "github.com/lib/pq"
	"log"
	_ "sort"
)

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

	tables, err := dbutils.GetTablesWithDependencies(db)
	if err != nil {
		log.Fatal(err)
	}

	sortedTables := dbutils.TopologicalSort(tables)

	println("Tables in topological order:")
	for _, table := range sortedTables {
		fmt.Printf("%s\n", table)
	}

	for _, table := range sortedTables {
		err := dbutils.GenerateAndInsertData(db, table, 10) // Generate 10 rows per table
		if err != nil {
			log.Printf("Error inserting data into %s: %v", table.Name, err)
		}
	}
}
