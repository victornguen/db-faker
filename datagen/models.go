package datagen

type TableRule struct {
	TableName string
	RowNum    int               `yaml:"num"`
	Rules     map[string]string `yaml:"columns"`
}

type TablesRules struct {
	Rules map[string]TableRule // key contains table name and value contains rules for that table
}
