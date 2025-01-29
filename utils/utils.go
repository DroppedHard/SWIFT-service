package utils

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/jbub/banking/swift"
	country "github.com/mikekonan/go-countries"
)

func ParseJson(r *http.Request, payload any) error {
	if r.Body == nil {
		return fmt.Errorf("missing request body")
	}

	return json.NewDecoder(r.Body).Decode(payload)
}

func WriteJson(w http.ResponseWriter, status int, v any) error {
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(status)

	return json.NewEncoder(w).Encode(v)
}

func WriteMessage(w http.ResponseWriter, status int, mess string) {
	WriteJson(w, status, map[string]string{"message": mess})
}

func WriteError(w http.ResponseWriter, status int, err error) {
	WriteJson(w, status, map[string]string{"message": err.Error()})
}

func BranchRegex(swiftCode string) string {
	return swiftCode[:8] + "???"
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
