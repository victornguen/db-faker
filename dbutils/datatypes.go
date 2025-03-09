package dbutils

import (
	"fmt"
	"github.com/go-faker/faker/v4"
	"github.com/go-faker/faker/v4/pkg/options"
	"github.com/samber/mo"
	funcutil "github.com/victornguen/db-faker/common"
	"github.com/victornguen/db-faker/datagen"
	"math/rand"
	"regexp"
	"strings"
)

type DataType interface {
	DefaultGenerator() func() string
}

// BigInt - signed eight-byte integer
// aliases:int8
type BigInt struct{}

func (b BigInt) DefaultGenerator() func() string {
	fun, err := datagen.RuleToGeneratorFunc("int")
	if err != nil {
		panic(err)
	}
	return fun
}

// BigSerial - autoincrementing eight-byte integer
// aliases:serial8
type BigSerial struct{}

func (b BigSerial) DefaultGenerator() func() string {
	fun, err := datagen.RuleToGeneratorFunc("int")
	if err != nil {
		panic(err)
	}
	return fun
}

// Bit - fixed-length bit string
type Bit struct {
	Len mo.Option[int]
}

func (b Bit) DefaultGenerator() func() string {
	val := rand.Intn(2)
	return func() string {
		return fmt.Sprintf("%d", val)
	}
}

// VarBit - variable-length bit string
// aliases:bit varying, varbit
type VarBit struct {
	Len mo.Option[int]
}

func (v VarBit) DefaultGenerator() func() string {
	// generate random variable-length bit string
	length, present := v.Len.Get()
	if !present {
		length = 8
		return func() string {
			return fmt.Sprintf("%08b", rand.Intn(256))
		}
	}
	return func() string {
		return fmt.Sprintf("%0*b", length, rand.Intn(1<<length))
	}
}

// Boolean - logical Boolean (true/false)
// aliases:bool
type Boolean struct{}

func (b Boolean) DefaultGenerator() func() string {
	return func() string {
		val := rand.Intn(2)
		return fmt.Sprintf("%t", val == 1)
	}
}

// Box - rectangular box on a plane
type Box struct{}

func (b Box) DefaultGenerator() func() string {
	return func() string {
		return fmt.Sprintf("(%d,%d),(%d,%d)", rand.Intn(100), rand.Intn(100), rand.Intn(100), rand.Intn(100))
	}
}

// ByteA - binary data ("byte array")
type ByteA struct{}

func (b ByteA) DefaultGenerator() func() string {
	return func() string {
		return "\\x012345"
	}
}

// Char - fixed-length character string
// aliases:char, character
type Char struct {
	Len mo.Option[int]
}

func (c Char) DefaultGenerator() func() string {
	length, present := c.Len.Get()
	if !present || length < 1 {
		length = 1
	}
	return func() string {
		return faker.Word(
			options.WithRandomStringLength(uint(length)),
		)
	}
}

// VarChar - variable-length character string
// aliases:character varying, varchar
type VarChar struct {
	MaxLen mo.Option[int]
}

func (v VarChar) DefaultGenerator() func() string {
	maxLen, present := v.MaxLen.Get()
	if !present || maxLen < 1 {
		maxLen = 255
	}
	return func() string {
		return faker.Sentence(
			options.WithRandomStringLength(uint(maxLen)),
		)
	}
}

// Cidr - IPv4 or IPv6 network address
type Cidr struct{}

func (c Cidr) DefaultGenerator() func() string {
	return func() string {
		return faker.IPv4()
	}
}

// Circle - circle on a plane
type Circle struct{}

func (c Circle) DefaultGenerator() func() string {
	return func() string {
		return fmt.Sprintf("<(%d,%d),%d>", rand.Intn(100), rand.Intn(100), rand.Intn(100))
	}
}

// Date - calendar date (year, month, day)
type Date struct{}

func (d Date) DefaultGenerator() func() string {
	return func() string {
		return faker.Date()
	}
}

// Float8 - double precision floating-point number (8 bytes)
// aliases:float8, double precision
type Float8 struct{}

func (f Float8) DefaultGenerator() func() string {
	return func() string {
		return fmt.Sprintf("%f", rand.Float64())
	}
}

// Inet - IPv4 or IPv6 host address
type Inet struct{}

func (i Inet) DefaultGenerator() func() string {
	return func() string {
		return faker.IPv4()
	}
}

// Int - signed four-byte integer
// aliases:integer, int, int4
type Int struct{}

func (i Int) DefaultGenerator() func() string {
	return func() string {
		return fmt.Sprintf("%d", rand.Intn(1000))
	}
}

// Interval - time span
type Interval struct{}

func (i Interval) DefaultGenerator() func() string {
	return func() string {
		intervals := []string{
			"1 year 3 hours 20 minutes",
			"2 weeks ago",
			"5 days 4 hours",
			"3 months 2 days",
			"1 year 3 hours 20 minutes",
			"5 minutes",
			"1 hour",
			"1 day",
			"1 week",
			"1 month",
			"4 weeks",
		}
		return intervals[rand.Intn(len(intervals))]
	}
}

// JSON - textual JSON data
type JSON struct{}

func (j JSON) DefaultGenerator() func() string {
	return func() string {
		return `{"key": "value"}`
	}
}

// JsonB - binary JSON data, decomposed
type JsonB struct{}

func (j JsonB) DefaultGenerator() func() string {
	return func() string {
		return `{"key": "value"}`
	}
}

// Line - infinite line on a plane
type Line struct{}

func (l Line) DefaultGenerator() func() string {
	return func() string {
		return fmt.Sprintf("{%d,%d,%d,%d}", rand.Intn(100), rand.Intn(100), rand.Intn(100), rand.Intn(100))
	}
}

// LSeg - line segment on a plane
type LSeg struct{}

func (l LSeg) DefaultGenerator() func() string {
	return func() string {
		return fmt.Sprintf("[(%d,%d),(%d,%d)]", rand.Intn(100), rand.Intn(100), rand.Intn(100), rand.Intn(100))
	}
}

// MacAddr - MAC (Media Access Control) address
type MacAddr struct{}

func (m MacAddr) DefaultGenerator() func() string {
	return func() string {
		return faker.MacAddress()
	}
}

// MacAddr8 - MAC (Media Access Control) address (EUI-64 format)
type MacAddr8 struct{}

func (m MacAddr8) DefaultGenerator() func() string {
	return func() string {
		return faker.MacAddress()
	}
}

// Money - currency amount
type Money struct{}

func (m Money) DefaultGenerator() func() string {
	return func() string {
		return fmt.Sprintf("%f", rand.Float64()*1000)
	}
}

// Numeric - exact numeric of selectable precision
// aliases:decimal, numeric
type Numeric struct {
	Precision mo.Option[int]
	Scale     mo.Option[int]
}

func (n Numeric) DefaultGenerator() func() string {
	return func() string {
		return fmt.Sprintf("%f", rand.Float64()*1000)
	}
}

// Path - geometric path on a plane
type Path struct{}

func (p Path) DefaultGenerator() func() string {
	return func() string {
		return fmt.Sprintf("((%d,%d),(%d,%d))", rand.Intn(100), rand.Intn(100), rand.Intn(100), rand.Intn(100))
	}
}

// PgLsn - PostgreSQL Log Sequence Number
type PgLsn struct{}

func (p PgLsn) DefaultGenerator() func() string {
	return func() string {
		return fmt.Sprintf("%d", rand.Intn(100))
	}
}

// PgSnapshot - user-level transaction ID snapshot
type PgSnapshot struct{}

func (p PgSnapshot) DefaultGenerator() func() string {
	return func() string {
		return fmt.Sprintf("%d", rand.Intn(100))
	}
}

// Point - geometric point on a plane
type Point struct{}

func (p Point) DefaultGenerator() func() string {
	return func() string {
		return fmt.Sprintf("(%d,%d)", rand.Intn(100), rand.Intn(100))
	}
}

// Polygon - closed geometric path on a plane
type Polygon struct{}

func (p Polygon) DefaultGenerator() func() string {
	return func() string {
		return fmt.Sprintf("((%d,%d),(%d,%d),(%d,%d))", rand.Intn(100), rand.Intn(100), rand.Intn(100), rand.Intn(100), rand.Intn(100), rand.Intn(100))
	}
}

// Real - single precision floating-point number (4 bytes)
// aliases: float4, real
type Real struct{}

func (r Real) DefaultGenerator() func() string {
	return func() string {
		return fmt.Sprintf("%f", rand.Float32())
	}
}

// SmallInt - signed two-byte integer
// aliases: int2, smallint
type SmallInt struct{}

func (s SmallInt) DefaultGenerator() func() string {
	return func() string {
		return fmt.Sprintf("%d", rand.Intn(100))
	}
}

// SmallSerial - autoincrementing two-byte integer
// aliases: smallserial, serial2
type SmallSerial struct{}

func (s SmallSerial) DefaultGenerator() func() string {
	return func() string {
		return fmt.Sprintf("%d", rand.Intn(100))
	}
}

// Serial - autoincrementing four-byte integer
// aliases: serial, serial4
type Serial struct{}

func (s Serial) DefaultGenerator() func() string {
	return func() string {
		return fmt.Sprintf("%d", rand.Intn(100))
	}
}

// Text - variable-length character string
// aliases: text
type Text struct{}

func (t Text) DefaultGenerator() func() string {
	return func() string {
		return faker.Sentence()
	}
}

type Time struct {
	WithTimeZone bool
	Precision    mo.Option[int]
}

func (t Time) DefaultGenerator() func() string {
	return func() string {
		return faker.TimeString()
	}
}

// TimeStamp - date and time (no time zone)
// aliases: timestamp
type TimeStamp struct {
	WithTimeZone bool
	Precision    mo.Option[int]
}

func (t TimeStamp) DefaultGenerator() func() string {
	return func() string {
		return faker.Timestamp()
	}
}

// TsQuery - text search query
type TsQuery struct{}

func (t TsQuery) DefaultGenerator() func() string {
	return func() string {
		return ""
	}
}

// TsVector - text search document
type TsVector struct{}

func (t TsVector) DefaultGenerator() func() string {
	return func() string {
		return ""
	}
}

// TxIDSnapshot - user-level transaction ID snapshot
type TxIDSnapshot struct{}

// UUID - universally unique identifier
type UUID struct{}

func (u UUID) DefaultGenerator() func() string {
	return func() string {
		return faker.UUIDHyphenated()
	}
}

// XML - XML data
type XML struct{}

func (x XML) DefaultGenerator() func() string {
	return func() string {
		return "<xml></xml>"
	}
}

func StringToDataType(s string) (DataType, error) {
	normalizedType := NormalizeType(s)
	matches := extractMatches(normalizedType)
	if len(matches) == 0 {
		return nil, fmt.Errorf("type name not found in %s", s)
	}

	typeName, n, m, additive, err := parseMatches(matches)
	if err != nil {
		return nil, err
	}

	return createDataType(typeName, n, m, additive)
}

func extractMatches(s string) []string {
	typePattern := regexp.MustCompile(`([a-zA-Z]+)\s*(\((\d+)(,\s*(\d+))?\))?\s*([a-zA-Z\s]+)?`)
	return typePattern.FindStringSubmatch(s)
}

func parseMatches(matches []string) (string, mo.Option[int], mo.Option[int], mo.Option[string], error) {
	typeName, isSome := funcutil.GetOpt(matches, 1).Get()
	if !isSome {
		return "", mo.None[int](), mo.None[int](), mo.None[string](), fmt.Errorf("type name not found")
	}

	nStr, nPresent := funcutil.GetOpt(matches, 3).Get()
	n := parseOptionalInt(nStr, nPresent)

	mStr, mPresent := funcutil.GetOpt(matches, 5).Get()
	m := parseOptionalInt(mStr, mPresent)

	additive := funcutil.GetOpt(matches, 6)

	return typeName, n, m, additive, nil
}

func parseOptionalInt(s string, present bool) mo.Option[int] {
	if !present || s == "" {
		return mo.None[int]()
	}
	val, ok := funcutil.ParseInt(s)
	if !ok {
		return mo.None[int]()
	}
	return mo.Some(val)
}

func createDataType(typeName string, n, m mo.Option[int], additive mo.Option[string]) (DataType, error) {
	switch typeName {
	case "bigint", "int8":
		return BigInt{}, nil
	case "bigserial", "serial8":
		return BigSerial{}, nil
	case "bit":
		return Bit{}, nil
	case "varbit":
		return VarBit{}, nil
	case "boolean", "bool":
		return Boolean{}, nil
	case "box":
		return Box{}, nil
	case "bytea":
		return ByteA{}, nil
	case "char", "character":
		return Char{Len: n}, nil
	case "varchar":
		return VarChar{MaxLen: n}, nil
	case "cidr":
		return Cidr{}, nil
	case "circle":
		return Circle{}, nil
	case "date":
		return Date{}, nil
	case "float8":
		return Float8{}, nil
	case "inet":
		return Inet{}, nil
	case "integer", "int", "int4":
		return Int{}, nil
	case "interval":
		return Interval{}, nil
	case "json":
		return JSON{}, nil
	case "jsonb":
		return JsonB{}, nil
	case "line":
		return Line{}, nil
	case "lseg":
		return LSeg{}, nil
	case "macaddr":
		return MacAddr{}, nil
	case "macaddr8":
		return MacAddr8{}, nil
	case "money":
		return Money{}, nil
	case "numeric", "decimal":
		return Numeric{Precision: n, Scale: m}, nil
	case "path":
		return Path{}, nil
	case "pg_lsn":
		return PgLsn{}, nil
	case "pg_snapshot":
		return PgSnapshot{}, nil
	case "point":
		return Point{}, nil
	case "polygon":
		return Polygon{}, nil
	case "real", "float4":
		return Real{}, nil
	case "smallint", "int2":
		return SmallInt{}, nil
	case "smallserial", "serial2":
		return SmallSerial{}, nil
	case "serial", "serial4":
		return Serial{}, nil
	case "text":
		return Text{}, nil
	case "time":
		return createTimeType(n, additive), nil
	case "timetz":
		return Time{WithTimeZone: true}, nil
	case "timestamp":
		return createTimeStampType(n, additive), nil
	case "timestamptz":
		return TimeStamp{WithTimeZone: true}, nil
	case "tsquery":
		return TsQuery{}, nil
	case "tsvector":
		return TsVector{}, nil
	// shoud not be column type
	//case "txid_snapshot":
	//	return TxIDSnapshot{}, nil
	case "uuid":
		return UUID{}, nil
	case "xml":
		return XML{}, nil
	default:
		return nil, fmt.Errorf("unknown data type: %s", typeName)
	}
}

func createTimeType(n mo.Option[int], additive mo.Option[string]) Time {
	if additive.IsPresent() {
		switch additive.MustGet() {
		case "with time zone":
			return Time{WithTimeZone: true, Precision: n}
		case "without time zone":
			return Time{WithTimeZone: false, Precision: n}
		}
	}
	return Time{WithTimeZone: false, Precision: n}
}

func createTimeStampType(n mo.Option[int], additive mo.Option[string]) TimeStamp {
	if additive.IsPresent() {
		switch additive.MustGet() {
		case "with time zone":
			return TimeStamp{WithTimeZone: true, Precision: n}
		case "without time zone":
			return TimeStamp{WithTimeZone: false, Precision: n}
		}
	}
	return TimeStamp{WithTimeZone: false, Precision: n}
}

func NormalizeType(s string) string {
	var res = strings.ReplaceAll(s, "bit varying", "varbit")
	res = strings.ReplaceAll(res, "character varying", "varchar")
	res = strings.ReplaceAll(res, "double precision", "float8")
	return res
}
