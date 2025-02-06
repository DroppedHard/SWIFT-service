package utils

import (
	"fmt"
	"reflect"
	"runtime"
	"strings"

	"github.com/jbub/banking/swift"
	country "github.com/mikekonan/go-countries"
)

func BranchRegex(swiftCode string) string {
	return swiftCode[:SwiftCodeLength] + "???"
}

func CountryCodeRegex(countryCode string) string {
	return "????" + countryCode + "?????"
}

func GetCountryNameFromCountryCode(countryCode string) string {
	result, ok := country.ByAlpha2Code(country.Alpha2Code(countryCode))
	if ok {
		return strings.ToUpper(result.NameStr())
	}
	return ""
}

func GetCountryCodeFromSwiftCode(swiftCode string) (string, error) {
	parsed, err := swift.Parse(swiftCode)
	if err != nil {
		return "", fmt.Errorf("failed to parse SWIFT code: %v", err)
	}
	return strings.ToUpper(parsed.CountryCode()), nil
}

func Xor(a bool, b bool) bool {
	return (a || b) && !(a && b)
}

func GetFunctionName(i interface{}) string {
	funcType := reflect.TypeOf(i)
	if funcType.Kind() != reflect.Func {
		return ""
	}
	fullName := runtime.FuncForPC(reflect.ValueOf(i).Pointer()).Name()
	return fullName[strings.LastIndex(fullName, ".")+1:]
}
