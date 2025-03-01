package dbutils

import (
	funcutil "db_faker/common"
	"fmt"
	"github.com/samber/mo"
	"regexp"
	"strings"
)

type DataType interface {
}

// BigInt - signed eight-byte integer
// aliases:int8
type BigInt struct{}

// BigSerial - autoincrementing eight-byte integer
// aliases:serial8
type BigSerial struct{}

// Bit - fixed-length bit string
type Bit struct {
	Len mo.Option[int]
}

// VarBit - variable-length bit string
// aliases:bit varying, varbit
type VarBit struct {
	Len mo.Option[int]
}

// Boolean - logical Boolean (true/false)
// aliases:bool
type Boolean struct{}

// Box - rectangular box on a plane
type Box struct{}

// ByteA - binary data ("byte array")
type ByteA struct{}

// Char - fixed-length character string
// aliases:char, character
type Char struct {
	Len mo.Option[int]
}

// VarChar - variable-length character string
// aliases:character varying, varchar
type VarChar struct {
	MaxLen mo.Option[int]
}

// Cidr - IPv4 or IPv6 network address
type Cidr struct{}

// Circle - circle on a plane
type Circle struct{}

// Date - calendar date (year, month, day)
type Date struct{}

// Float8 - double precision floating-point number (8 bytes)
// aliases:float8, double precision
type Float8 struct{}

// Inet - IPv4 or IPv6 host address
type Inet struct{}

// Int - signed four-byte integer
// aliases:integer, int, int4
type Int struct{}

// Interval - time span
type Interval struct{}

// JSON - textual JSON data
type JSON struct{}

// JsonB - binary JSON data, decomposed
type JsonB struct{}

// Line - infinite line on a plane
type Line struct{}

// LSeg - line segment on a plane
type LSeg struct{}

// MacAddr - MAC (Media Access Control) address
type MacAddr struct{}

// MacAddr8 - MAC (Media Access Control) address (EUI-64 format)
type MacAddr8 struct{}

// Money - currency amount
type Money struct{}

// Numeric - exact numeric of selectable precision
// aliases:decimal, numeric
type Numeric struct {
	Precision mo.Option[int]
	Scale     mo.Option[int]
}

// Path - geometric path on a plane
type Path struct{}

// PgLsn - PostgreSQL Log Sequence Number
type PgLsn struct{}

// PgSnapshot - user-level transaction ID snapshot
type PgSnapshot struct{}

// Point - geometric point on a plane
type Point struct{}

// Polygon - closed geometric path on a plane
type Polygon struct{}

// Real - single precision floating-point number (4 bytes)
// aliases: float4, real
type Real struct{}

// SmallInt - signed two-byte integer
// aliases: int2, smallint
type SmallInt struct{}

// SmallSerial - autoincrementing two-byte integer
// aliases: smallserial, serial2
type SmallSerial struct{}

// Serial - autoincrementing four-byte integer
// aliases: serial, serial4
type Serial struct{}

// Text - variable-length character string
// aliases: text
type Text struct{}

type Time struct {
	WithTimeZone bool
	Precision    mo.Option[int]
}

// TimeStamp - date and time (no time zone)
// aliases: timestamp
type TimeStamp struct {
	WithTimeZone bool
	Precision    mo.Option[int]
}

// TsQuery - text search query
type TsQuery struct{}

// TsVector - text search document
type TsVector struct{}

// TxIDSnapshot - user-level transaction ID snapshot
type TxIDSnapshot struct{}

// UUID - universally unique identifier
type UUID struct{}

// XML - XML data
type XML struct{}

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
	case "txid_snapshot":
		return TxIDSnapshot{}, nil
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
