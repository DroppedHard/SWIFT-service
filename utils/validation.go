package utils

import (
	"fmt"

	"github.com/go-playground/validator/v10"
	"github.com/jbub/banking/swift"
	country "github.com/mikekonan/go-countries"
)

var Validate = validator.New()

type ValidationError struct {
	Errors map[string]string
}

func (v ValidationError) Error() string {
	return fmt.Sprintf("validation failed: %v", v.Errors)
}

func init() {
	errSwift := Validate.RegisterValidation(ValidatorSwiftCode, swiftCodeValidation)
	errIso2 := Validate.RegisterValidation(ValidatorCountryIso2, countryIso2Validation)
	errBool := Validate.RegisterValidation("boolRequired", boolValidation)
	if errSwift != nil {
		fmt.Println("failed to register swiftCode validation: ", errSwift)
	}
	if errIso2 != nil {
		fmt.Println("failed to register countryISO2 validation:", errIso2)
	}
	if errBool != nil {
		fmt.Println("failed to register bool validation:", errBool)
	}
}

func swiftCodeValidation(fl validator.FieldLevel) bool {
	swiftCode := fl.Field().String()
	if err := swift.Validate(swiftCode); err != nil {
		return false
	}
	return true
}

func countryIso2Validation(fl validator.FieldLevel) bool {
	countryCode := fl.Field().String()
	_, ok := country.ByAlpha2Code(country.Alpha2Code(countryCode))
	return ok
}

func boolValidation(fl validator.FieldLevel) bool { // TODO - verify if this is necessary
	return fl.Field().String() == "true" || fl.Field().String() == "false"
}
