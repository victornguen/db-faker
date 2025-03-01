package dbutils

import (
	"testing"
)

var dataTypes = []string{
	"numeric (10, 2)",
	"time (6) with time zone",
	"time with time zone",
	"time(4)",
	"integer",
}

func TestStringToDataType(t *testing.T) {
	for _, dtStr := range dataTypes {
		var dt, err = StringToDataType(dtStr)
		if err != nil {
			t.Error(err)
		}
		if dt == nil {
			t.Error("DataType is nil")
		}
	}
}
