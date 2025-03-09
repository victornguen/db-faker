package datagen

import (
	"strconv"
	"testing"
	"time"
)

func Test_ruleToGenerator_int(t *testing.T) {
	rule := "int(1, 10)"

	r, err := RuleToGeneratorFunc(rule)

	if err != nil {
		t.Errorf("Error: %v", err)
	}

	val, err := strconv.Atoi(r())
	if err != nil {
		t.Errorf("Error: %v", err)
	}

	if val < 1 || val > 10 {
		t.Errorf("Error: %v", val)
	}

	t.Logf("Value: %v", val)

}

func Test_ruleToGenerator_timestamp(t *testing.T) {
	rule := "timestamp"

	r, err := RuleToGeneratorFunc(rule)

	if err != nil {
		t.Errorf("Error: %v", err)
	}

	val := r()

	// parse timestamp
	ts, err := time.Parse("2006-01-02 15:04:05", val)
	if err != nil {
		t.Errorf("Error: %v", err)
	}

	t.Logf("Value: %v", ts)

}
