package main

import (
	"context"
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"github.com/urfave/cli/v3"
	"github.com/victornguen/db-faker/datagen"
	"github.com/victornguen/db-faker/dbutils"
	"log"
	"os"
	_ "sort"
)

func main() {
	app := &cli.Command{
		Name:  "db_faker",
		Usage: "Generate and insert fake data into a PostgreSQL database",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:     "host",
				Usage:    "PostgreSQL host",
				Value:    "localhost",
				Required: false,
			},
			&cli.IntFlag{
				Name:     "port",
				Aliases:  []string{"p"},
				Usage:    "PostgreSQL port",
				Value:    5432,
				Required: false,
			},
			&cli.StringFlag{
				Name:     "user",
				Usage:    "PostgreSQL user",
				Required: true,
			},
			&cli.StringFlag{
				Name:     "password",
				Usage:    "PostgreSQL password",
				Required: true,
			},
			&cli.StringFlag{
				Name:     "dbname",
				Aliases:  []string{"db"},
				Usage:    "PostgreSQL database name",
				Required: true,
			},
			&cli.StringFlag{
				Name:     "rules",
				Aliases:  []string{"r"},
				Usage:    "Path to YAML file containing rules",
				Value:    "./gen_settings.yaml",
				Required: false,
			},
		},
		Commands: []*cli.Command{
			{
				Name:   "generate",
				Usage:  "Generate and insert fake data",
				Action: generateData,
			},
		},
	}

	if err := app.Run(context.Background(), os.Args); err != nil {
		log.Fatal(err)
	}

}

func generateData(c context.Context, command *cli.Command) error {
	host := command.String("host")
	port := command.Int("port")
	user := command.String("user")
	password := command.String("password")
	dbname := command.String("dbname")

	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)

	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		return err
	}
	defer db.Close()

	rules, err := datagen.LoadRulesFromYAMLFile("./gen_settings.yaml")
	if err != nil {
		return err
	}

	tables, err := dbutils.GetTablesWithDependencies(db)
	if err != nil {
		return err
	}

	sortedTables := dbutils.TopologicalSort(tables)

	err = dbutils.ApplyRulesToTables(&sortedTables, rules)
	if err != nil {
		return err
	}

	for _, table := range sortedTables {
		err := dbutils.GenerateAndInsertData(db, table)
		if err != nil {
			log.Printf("Error inserting data into %s: %v", table.Name, err)
		}
	}

	return nil
}
