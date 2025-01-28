package swiftCode

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/DroppedHard/SWIFT-service/types"
	"github.com/DroppedHard/SWIFT-service/utils"
	"github.com/go-playground/validator/v10"
	"github.com/gorilla/mux"
)

func validateInput(input interface{}, tags string) error {
	var err error

	if tags == "" {
		err = utils.Validate.Struct(input)
	} else {
		err = utils.Validate.Var(input, tags)
	}

	if err == nil {
		return nil
	}

	errors := make(map[string]string)
	if validationErrs, ok := err.(validator.ValidationErrors); ok {
		for _, fieldErr := range validationErrs {
			errors[fieldErr.Field()] = fmt.Sprintf("validation failed on '%s' tag", fieldErr.Tag())
		}
	} else {
		errors["error"] = err.Error()
	}
	return utils.ValidationError{Errors: errors}
}

func validateSwiftCode(r *http.Request) error {
	swiftCode := mux.Vars(r)["swift-code"]
	return validateInput(swiftCode, "required,swiftCode")
}

func validateCountryCode(r *http.Request) error {
	countryCode := mux.Vars(r)["countryISO2"]
	return validateInput(countryCode, "required,countryISO2")
}

func validateAddSwiftCode(ctx context.Context, payload *types.BankDataDetails) error {
	if err := utils.Validate.Struct(payload); err != nil {
		return fmt.Errorf("invalid payload structure: %w", err)
	}

	expectedCountryCode, err := utils.GetCountryCodeFromSwiftCode(payload.SwiftCode)
	if err != nil {
		return fmt.Errorf("invalid SWIFT code: %w", err)
	}
	if payload.CountryISO2 != expectedCountryCode {
		return fmt.Errorf("countryISO2 '%s' does not match the country derived from SWIFT code '%s'", payload.CountryISO2, expectedCountryCode)
	}
	expectedCountryName := utils.GetCountryNameFromCountryCode(payload.CountryISO2)
	if !strings.EqualFold(payload.CountryName, expectedCountryName) {
		return fmt.Errorf("countryName '%s' does not match the country derived from countryISO2 '%s'", payload.CountryName, expectedCountryName)
	}
	if !utils.Xor(payload.IsHeadquarter, strings.HasSuffix(payload.SwiftCode, "XXX")) {
		return fmt.Errorf("isHeadquarter value '%v' does not match the swiftCode value '%s'", payload.IsHeadquarter, payload.SwiftCode)
	}

	return nil
}
