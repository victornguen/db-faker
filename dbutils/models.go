package dbutils

type Table struct {
	Name        string
	Columns     []Column
	DependsOn   []string
	PrimaryKeys map[string]bool
	RowNum      int32
	Rules       map[string]func() string // key contains column name and value contains function to generate data
}

type Column struct {
	Name         string
	DataType     DataType
	IsForeignKey bool
	RefTable     string
	DataGen      func() string
}

type TableRule struct {
	TableName string
	RowNum    int32             `yaml:"num"`
	Rules     map[string]string `yaml:"columns"`
}

type TablesRules struct {
	Rules map[string]TableRule // key contains table name and value contains rules for that table
}

//type TableDependency struct {
//	TableName    string
//	Dependencies []string
//}
