package datagen

import (
	"fmt"
	"github.com/go-faker/faker/v4"
	"github.com/go-faker/faker/v4/pkg/options"
	funcutil "github.com/victornguen/db-faker/common"
	"math/rand"
	"os"
	"reflect"
	"regexp"
	"strconv"
	"strings"
)

type DataGenerator[T any] struct {
	gen func() T
}

func RuleToGeneratorFunc(rule string) (func() string, error) {
	var customRuleFunc, err = customGenerator(rule)
	if err == nil {
		return customRuleFunc, nil
	}
	var regex = `([a-zA-Z]+)\s*(\((\d+)(,\s*(\d+))?\))?\s*`
	var pattern = regexp.MustCompile(regex)
	if rule == "" {
		return nil, fmt.Errorf("empty rule")
	}
	rule = strings.ToLower(rule)
	var matches = pattern.FindStringSubmatch(rule)
	if len(matches) < 5 {
		return nil, fmt.Errorf("invalid rule, rule must match pattern: %s", regex)
	}
	_ = matches[5]
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
		generator = genF(faker.FirstName)
	case "lastname":
		generator = genF(faker.LastName)
	case "email":
		generator = genF(faker.Email)
	case "username":
		generator = genF(faker.Username)
	case "currency":
		generator = genF(faker.Currency)
	case "ccnumber":
		generator = genF(faker.CCNumber)
	case "cctype":
		generator = genF(faker.CCType)
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
		generator = genF(faker.Phonenumber)
	case "date":
		generator = genF(faker.Date)
	case "dayofweek":
		generator = genF(faker.DayOfWeek)
	case "month":
		generator = genF(faker.MonthName)
	case "year":
		generator = genF(faker.YearString)
	case "time":
		generator = genF(faker.TimeString)
	case "datetime", "timestamp":
		generator = genF(faker.Timestamp)
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
		generator = genF(faker.Paragraph)
	case "ipv4":
		generator = genF(faker.IPv4)
	case "ipv6":
		generator = genF(faker.IPv6)
	case "mac":
		generator = genF(faker.MacAddress)
	case "url":
		generator = genF(faker.URL)
	case "useragent":
		generator = useragentGenerator()
	default:
		return nil, fmt.Errorf("unknown rule: %s", typeName)
	}

	return generator, nil
}

// for rule looks like: "oneof[one%10, two%50, three%40]" or "oneof[one%40, two%50]" or "constant[one]"
// one%40 means 'one' generates with 40% probability
// constant[one] means 'one' generates always
func customGenerator(rule string) (func() string, error) {
	var regex = `([a-zA-Z]+)\s*\[(.+)\]`
	var pattern = regexp.MustCompile(regex)
	if rule == "" {
		return nil, fmt.Errorf("empty rule: %s", rule)
	}
	var matches = pattern.FindStringSubmatch(rule)
	if len(matches) < 3 {
		return nil, fmt.Errorf("invalid rule, rule must match pattern: %s", regex)
	}
	var typeName = matches[1]
	var values = strings.Split(matches[2], ",")
	var generator func() string
	typeName = strings.ToLower(typeName)
	switch typeName {
	case "oneof":
		var probs = make([]int, 0)
		var names = make([]string, 0)
		for _, v := range values {
			var parts = strings.Split(v, "%")
			if len(parts) != 2 {
				return nil, fmt.Errorf("invalid rule, oneof rule must contain name and probability")
			}
			name := parts[0]
			prob, err := strconv.Atoi(parts[1])
			if err != nil {
				return nil, fmt.Errorf("invalid rule, probability must be an integer")
			}
			names = append(names, name)
			probs = append(probs, prob)
		}
		generator = func() string {
			// generate names[i] with probability prob[i]
			var r = rand.Intn(100)
			var sum = 0
			for i, prob := range probs {
				sum += prob
				if r < sum {
					return names[i]
				}
			}
			return names[len(names)-1]
		}
	case "constant":
		if len(values) != 1 {
			return nil, fmt.Errorf("invalid rule, constant rule must contain one value")
		}
		generator = func() string {
			return values[0]
		}
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

func genF(f func(...options.OptionFunc) string) func() string {
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
			_, _ = fmt.Fprintf(os.Stderr, "error generating user agent: %v\n", err)
			return "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/58.0.3029.110 Safari/537.36"
		}
		return ua.(string)
	}
}
