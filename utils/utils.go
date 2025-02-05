package utils

import (
	"fmt"

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
		return result.NameStr()
	}
	return ""
}

func GetCountryCodeFromSwiftCode(swiftCode string) (string, error) {
	parsed, err := swift.Parse(swiftCode)
	if err != nil {
		return "", fmt.Errorf("failed to parse SWIFT code: %v", err)
	}
	return parsed.CountryCode(), nil
}

func Xor(a bool, b bool) bool {
	return (a || b) && !(a && b)
}
