package utils

import (
	"fmt"
	"regexp"

	"github.com/go-playground/validator/v10"
)

var Validate = validator.New()

func init() {
	errSwift := Validate.RegisterValidation("swiftCode", swiftCodeValidation)
	errIso2 := Validate.RegisterValidation("countryISO2", countryIso2Validation)
	if errSwift != nil || errIso2 != nil {
		fmt.Println("failed to register validation - swiftCode: ", errSwift, " countryISO2: ", errIso2)
		return
	}
}

func swiftCodeValidation(fl validator.FieldLevel) bool {
	swiftCode := fl.Field().String()
	swiftRegex := `^[A-Z0-9]{8}([A-Z0-9]{3})?$` // SWIFT code pattern.

	// Compile regex and validate.
	re := regexp.MustCompile(swiftRegex)
	return re.MatchString(swiftCode)
}

func countryIso2Validation(fl validator.FieldLevel) bool {
	swiftCode := fl.Field().String()
	swiftRegex := `^[A-Z0-9]{2}$` // SWIFT code pattern.

	// Compile regex and validate.
	re := regexp.MustCompile(swiftRegex)
	return re.MatchString(swiftCode)
}
