package datagen

import (
	"database/sql"
	"fmt"
)

const (
	getRulesQuery = `
		SELECT table_name, col_name, rule
		FROM dbfaker_col_rules
	`
)

func loadRules(db *sql.DB) (map[string]map[string]string, error) {
	rows, err := db.Query(getRulesQuery)
	if err != nil {
		return nil, fmt.Errorf("error loading rules: %v", err)
	}
	defer rows.Close()

	rules := make(map[string]map[string]string)
	for rows.Next() {
		var table, column, rule string
		if err := rows.Scan(&table, &column, &rule); err != nil {
			return nil, err
		}

		if rules[table] == nil {
			rules[table] = make(map[string]string)
		}
		rules[table][column] = rule
	}
	return rules, nil
}
