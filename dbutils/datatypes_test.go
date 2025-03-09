package dbutils

import (
	"github.com/stretchr/testify/assert"
	"strconv"
	"testing"
	"time"
)

func TestStringToDataType(t *testing.T) {
	var dataTypes = []string{
		"numeric (10, 2)",
		"time (6) with time zone",
		"time with time zone",
		"time(4)",
		"integer",
	}
	var results = make([]DataType, 0)
	for _, dtStr := range dataTypes {
		var dt, err = StringToDataType(dtStr)
		results = append(results, dt)
		if err != nil {
			t.Error(err)
		}
		if dt == nil {
			t.Error("DataType is nil")
		}
	}
	assert.Equal(t, 10, results[0].(Numeric).Precision.OrElse(0))
	assert.Equal(t, 2, results[0].(Numeric).Scale.OrElse(0))
	assert.Equal(t, 6, results[1].(Time).Precision.OrElse(0))
	assert.Equal(t, 0, results[2].(Time).Precision.OrElse(0))
	assert.Equal(t, 4, results[3].(Time).Precision.OrElse(0))
	switch results[4].(type) {
	case Int:
	default:
		t.Error("Invalid type")
	}

}

func TestBoolean_DefaultGenerator(t *testing.T) {
	var b Boolean
	var gen = b.DefaultGenerator()
	if gen == nil {
		t.Error("Generator is nil")
	}
	assert.Contains(t, []string{"true", "false"}, gen())
}

func TestBigInt_DefaultGenerator(t *testing.T) {
	var b BigInt
	var gen = b.DefaultGenerator()
	if gen == nil {
		t.Error("Generator is nil")
	}
	assert.NotEmpty(t, gen())
}

func TestInt_DefaultGenerator(t *testing.T) {
	var i Int
	var gen = i.DefaultGenerator()
	if gen == nil {
		t.Error("Generator is nil")
	}
	assert.NotEmpty(t, gen())
}

func TestFloat8_DefaultGenerator(t *testing.T) {
	var f Float8
	var gen = f.DefaultGenerator()
	if gen == nil {
		t.Error("Generator is nil")
	}
	result := gen()
	_, err := strconv.ParseFloat(result, 64)
	if err != nil {
		t.Error(err)
	}
	assert.NotEmpty(t, result)
}

func TestVarChar_DefaultGenerator(t *testing.T) {
	var s VarChar
	var gen = s.DefaultGenerator()
	if gen == nil {
		t.Error("Generator is nil")
	}
	assert.NotEmpty(t, gen())
}

func TestDate_DefaultGenerator(t *testing.T) {
	var d Date
	var gen = d.DefaultGenerator()
	if gen == nil {
		t.Error("Generator is nil")
	}
	result := gen()

	_, err := time.Parse("2006-01-02", result)
	if err != nil {
		t.Error(err)
	}
	assert.NotEmpty(t, result)
}

func TestTimeStamp_DefaultGenerator(t *testing.T) {
	var ts TimeStamp
	var gen = ts.DefaultGenerator()
	if gen == nil {
		t.Error("Generator is nil")
	}
	result := gen()

	_, err := time.Parse("2006-01-02 15:04:05", result)
	if err != nil {
		t.Error(err)
	}
	assert.NotEmpty(t, result)
}

func TestPolygon_DefaultGenerator(t *testing.T) {
	var p Polygon
	var gen = p.DefaultGenerator()
	if gen == nil {
		t.Error("Generator is nil")
	}
	assert.NotEmpty(t, gen())
}
