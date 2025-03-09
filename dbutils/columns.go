package dbutils

import "database/sql"

const getColumnsQuery = `
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

func GetColumns(db *sql.DB, tableName string) ([]Column, error) {
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
		col.DataGen = dataType.DefaultGenerator()
		columns = append(columns, col)
	}

	return columns, nil
}
