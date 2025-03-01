package dbutils

type Table struct {
	Name        string
	Columns     []Column
	DependsOn   []string
	PrimaryKeys map[string]bool
}

type Column struct {
	Name         string
	DataType     DataType
	IsForeignKey bool
	RefTable     string
}

//type TableDependency struct {
//	TableName    string
//	Dependencies []string
//}
