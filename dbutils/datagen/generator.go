package datagen

import (
	funcutil "db_faker/common"
	"fmt"
	"github.com/go-faker/faker/v4"
	"github.com/go-faker/faker/v4/pkg/options"
	"math/rand"
	"os"
	"reflect"
	"regexp"
	"strings"
)

type DataGenerator[T any] struct {
	gen func() T
}

func RuleToGeneratorFunc(rule string) (func() string, error) {
	var pattern = regexp.MustCompile(`([a-zA-Z]+)\s*(\((\d+)(,\s*(\d+))?\))?\s*`)
	if rule == "" {
		return nil, fmt.Errorf("empty rule")
	}
	rule = strings.ToLower(rule)
	var matches = pattern.FindStringSubmatch(rule)
	if len(matches) < 5 {
		return nil, fmt.Errorf("0 matches found")
	}
	var typeName = matches[1]
	var nStr = matches[3]
	var mStr = matches[5]
	var generator func() string

	switch typeName {
	case "int", "integer":
		gen, _ := intGenerator(nStr, mStr)
		generator = func() string {
			return fmt.Sprintf("%d", gen())
		}
	case "sentence", "text":
		generator, _ = sentenceGenerator(nStr)
	case "firstname", "name":
		generator = gen(faker.FirstName)
	case "lastname":
		generator = gen(faker.LastName)
	case "email":
		generator = gen(faker.Email)
	case "username":
		generator = gen(faker.Username)
	case "currency":
		generator = gen(faker.Currency)
	case "ccnumber":
		generator = gen(faker.CCNumber)
	case "cctype":
		generator = gen(faker.CCType)
	case "country":
		generator = genMap(faker.GetCountryInfo, func(info faker.CountryInfo) string {
			return info.Name
		})
	case "city":
		generator = genMap(faker.GetRealAddress, func(info faker.RealAddress) string {
			return info.City
		})
	case "address":
		generator = genMap(faker.GetRealAddress, func(info faker.RealAddress) string {
			return info.Address
		})
	case "state":
		generator = genMap(faker.GetRealAddress, func(info faker.RealAddress) string {
			return info.State
		})
	case "postalcode":
		generator = genMap(faker.GetRealAddress, func(info faker.RealAddress) string {
			return info.PostalCode
		})
	case "latitude", "lat":
		generator = genMap(faker.GetRealAddress, func(info faker.RealAddress) string {
			return fmt.Sprintf("%f", info.Coordinates.Latitude)
		})
	case "longitude", "lon":
		generator = genMap(faker.GetRealAddress, func(info faker.RealAddress) string {
			return fmt.Sprintf("%f", info.Coordinates.Longitude)
		})
	case "phone":
		generator = gen(faker.Phonenumber)
	case "date":
		generator = gen(faker.Date)
	case "dayofweek":
		generator = gen(faker.DayOfWeek)
	case "month":
		generator = gen(faker.MonthName)
	case "year":
		generator = gen(faker.YearString)
	case "time":
		generator = gen(faker.TimeString)
	case "datetime", "timestamp":
		generator = gen(faker.Timestamp)
	case "bloodtype":
		generator = genMap(faker.GetBlood, func(b faker.Blooder) string {
			bt, _ := b.BloodType(reflect.Value{})
			return bt.(string)
		})
	case "bloodrhfactor":
		generator = genMap(faker.GetBlood, func(b faker.Blooder) string {
			bf, _ := b.BloodRHFactor(reflect.Value{})
			return bf.(string)
		})
	case "bloodgroup":
		generator = genMap(faker.GetBlood, func(b faker.Blooder) string {
			bg, _ := b.BloodGroup(reflect.Value{})
			return bg.(string)
		})
	case "paragraph":
		generator = gen(faker.Paragraph)
	case "ipv4":
		generator = gen(faker.IPv4)
	case "ipv6":
		generator = gen(faker.IPv6)
	case "mac":
		generator = gen(faker.MacAddress)
	case "url":
		generator = gen(faker.URL)
	case "useragent":
		generator = useragentGenerator()
	default:
		return nil, fmt.Errorf("unsupported type: %s", typeName)
	}

	return generator, nil
}

func intGenerator(nStr, mStr string) (func() int, error) {
	var generator func() int
	if nStr != "" && mStr != "" {
		n, nPresent := funcutil.ParseInt(nStr)
		m, mPresent := funcutil.ParseInt(mStr)
		if !nPresent || !mPresent {
			return nil, fmt.Errorf("invalid rule")
		}
		generator = func() int {
			i, _ := faker.RandomInt(n, m)
			return i[0]
		}
	} else if nStr != "" {
		n, nPresent := funcutil.ParseInt(nStr)
		if !nPresent {
			return nil, fmt.Errorf("invalid rule")
		}
		generator = func() int {
			i, _ := faker.RandomInt(n)
			return i[0]
		}
	} else {
		generator = func() int {
			return rand.Int()
		}
	}
	return generator, nil
}

func gen(f func(...options.OptionFunc) string) func() string {
	return func() string {
		return f()
	}
}

func genMap[A any, B any](f func(...options.OptionFunc) A, transform func(A) B) func() B {
	return func() B {
		return transform(f())
	}
}

func sentenceGenerator(nStr string) (func() string, error) {
	var generator func() string
	if nStr != "" {
		n, nPresent := funcutil.ParseUint(nStr)
		if !nPresent {
			return nil, fmt.Errorf("invalid rule")
		}
		generator = func() string {
			return faker.Sentence(
				options.WithRandomStringLength(n))
		}
	} else {
		generator = func() string {
			return faker.Sentence()
		}
	}
	return generator, nil
}

func useragentGenerator() func() string {
	return func() string {
		ua, err := faker.GetUserAgent().UserAgent(reflect.Value{})
		if err != nil {
			fmt.Fprintf(os.Stderr, "error generating user agent: %v\n", err)
			return ""
		}
		return ua.(string)
	}
}
