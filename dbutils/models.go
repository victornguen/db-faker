package dbutils

type Table struct {
	Name        string
	Columns     map[string]Column
	DependsOn   []string
	PrimaryKeys map[string]bool
	RowNum      int
	Rules       map[string]func() string // key contains column name and value contains function to generate data
}

type Column struct {
	Name         string
	DataType     DataType
	IsForeignKey bool
	RefTable     string
	DataGen      func() string
}

//type TableDependency struct {
//	TableName    string
//	Dependencies []string
//}
