package datagen

import (
	"fmt"
	"gopkg.in/yaml.v2"
	"os"
)

func LoadRulesFromYAMLFile(path string) (TablesRules, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return TablesRules{}, fmt.Errorf("error reading YAML file: %v", err)
	}

	// Parse the YAML into TablesRules struct
	var tablesRules TablesRules
	err = yaml.Unmarshal(data, &tablesRules)
	if err != nil {
		return TablesRules{}, fmt.Errorf("error unmarshalling YAML: %v", err)
	}

	return tablesRules, nil
}
