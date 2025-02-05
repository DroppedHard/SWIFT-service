package api_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/DroppedHard/SWIFT-service/service/api"
	"github.com/DroppedHard/SWIFT-service/types"
	"github.com/DroppedHard/SWIFT-service/utils"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
)

type ValidatorsTestCase struct {
	Description string
	Input       interface{}
	Tags        string
	ExpectedErr string
}

func TestValidateInput(t *testing.T) {
	tests := []ValidatorsTestCase{
		{
			Description: "Valid input",
			Input: struct {
				Field1 string `validate:"required"`
			}{Field1: "valid value"},
			Tags:        "required",
			ExpectedErr: "",
		},
		{
			Description: "Missing required field",
			Input: struct {
				Field1 string `validate:"required"`
			}{},
			Tags:        "required",
			ExpectedErr: "validation failed on 'required' tag",
		},
		{
			Description: "Invalid email format",
			Input: struct {
				Email string `validate:"email"`
			}{Email: "invalid-email"},
			Tags:        "email",
			ExpectedErr: "validation failed on 'email' tag",
		},
	}

	for _, test := range tests {
		t.Run(test.Description, func(t *testing.T) {
			err := api.ValidateInput(test.Input, test.Tags)

			if test.ExpectedErr != "" {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), test.ExpectedErr)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestValidateSwiftCode(t *testing.T) {
	tests := []struct {
		Description string
		SwiftCode   string
		ExpectedErr string
	}{
		{
			Description: "Valid Swift Code",
			SwiftCode:   "ALBPPLXAXXX",
			ExpectedErr: "",
		},
		{
			Description: "Missing Swift Code",
			SwiftCode:   "",
			ExpectedErr: "validation failed on 'required' tag",
		},
		{
			Description: "Invalid Swift Code",
			SwiftCode:   "ALB)__XAXXX",
			ExpectedErr: "validation failed on 'swiftCode' tag",
		},
	}

	for _, test := range tests {
		t.Run(test.Description, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, "/swift-codes/"+test.SwiftCode, nil)
			vars := map[string]string{
				utils.PathParamSwiftCode: test.SwiftCode,
			}
			req = mux.SetURLVars(req, vars)

			err := api.ValidateSwiftCode(req)

			if test.ExpectedErr != "" {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), test.ExpectedErr)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestValidateCountryCode(t *testing.T) {
	tests := []struct {
		Description string
		CountryCode string
		ExpectedErr string
	}{
		{
			Description: "Valid Country Code",
			CountryCode: "US",
			ExpectedErr: "",
		},
		{
			Description: "Missing Country Code",
			CountryCode: "",
			ExpectedErr: "validation failed on 'required' tag",
		},
		{
			Description: "Invalid Country Code",
			CountryCode: "XX",
			ExpectedErr: "validation failed on 'countryISO2' tag",
		},
	}

	for _, test := range tests {
		t.Run(test.Description, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, "/swift-codes/country/"+test.CountryCode, nil)
			vars := map[string]string{
				utils.PathParamCountryIso2: test.CountryCode,
			}
			req = mux.SetURLVars(req, vars)

			err := api.ValidateCountryCode(req)

			if test.ExpectedErr != "" {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), test.ExpectedErr)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestValidatePostSwiftCodePayload(t *testing.T) {
	tests := []struct {
		Description string
		Payload     *types.BankDataDetails
		ExpectedErr string
	}{
		{
			Description: "Valid Payload",
			Payload: &types.BankDataDetails{
				BankDataCore: types.BankDataCore{
					SwiftCode:     "ALBPPLPWXXX",
					Address:       "Valid address",
					BankName:      "Valid bank name",
					CountryIso2:   "PL",
					IsHeadquarter: true,
				},
				CountryName: "Poland",
			},
			ExpectedErr: "",
		},
		{
			Description: "Invalid Country Code",
			Payload: &types.BankDataDetails{
				BankDataCore: types.BankDataCore{
					SwiftCode:     "ALBPPLPWXXX",
					Address:       "Valid address",
					BankName:      "Valid bank name",
					CountryIso2:   "CA",
					IsHeadquarter: true,
				},
				CountryName: "Canada",
			},
			ExpectedErr: "countryISO2 'CA' does not match the country derived from SWIFT code 'PL'",
		},
		{
			Description: "Invalid Headquarter Flag",
			Payload: &types.BankDataDetails{
				BankDataCore: types.BankDataCore{
					SwiftCode:     "ALBPPLPWXXX",
					Address:       "Valid address",
					BankName:      "Valid bank name",
					CountryIso2:   "PL",
					IsHeadquarter: false,
				},
				CountryName: "Poland",
			},
			ExpectedErr: "isHeadquarter value 'false' does not match the swiftCode value 'ALBPPLPWXXX'",
		},
		{
			Description: "Invalid SWIFT code",
			Payload: &types.BankDataDetails{
				BankDataCore: types.BankDataCore{
					SwiftCode:     "INVALIDSWIFTCODE",
					Address:       "Valid address",
					BankName:      "Valid bank name",
					CountryIso2:   "US",
					IsHeadquarter: true,
				},
				CountryName: "United States",
			},
			ExpectedErr: "Field validation for 'SwiftCode' failed on the 'len' tag",
		},
	}

	for _, test := range tests {
		t.Run(test.Description, func(t *testing.T) {
			err := api.ValidatePostSwiftCodePayload(context.Background(), test.Payload)

			if test.ExpectedErr != "" {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), test.ExpectedErr)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
