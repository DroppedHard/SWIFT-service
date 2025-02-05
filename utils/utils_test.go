package utils_test

import (
	"testing"

	"github.com/DroppedHard/SWIFT-service/utils"
	"github.com/stretchr/testify/assert"
)

func TestBranchRegex(t *testing.T) {
	result := utils.BranchRegex("ALBPPLPWXXX")
	assert.Equal(t, "ALBPPLPW???", result)
}

func TestCountryCodeRegex(t *testing.T) {
	result := utils.CountryCodeRegex("PL")
	assert.Equal(t, "????PL?????", result)
}

func TestGetCountryNameFromCountryCode(t *testing.T) {
	result := utils.GetCountryNameFromCountryCode("PL")
	assert.Equal(t, "Poland", result)

	result = utils.GetCountryNameFromCountryCode("XX")
	assert.Equal(t, "", result)
}

func TestGetCountryCodeFromSwiftCode(t *testing.T) {
	result, err := utils.GetCountryCodeFromSwiftCode("ALBPPLPWXXX")
	assert.NoError(t, err)
	assert.Equal(t, "PL", result)

	result, err = utils.GetCountryCodeFromSwiftCode("____")
	assert.Error(t, err)
	assert.Empty(t, result)
}

func TestXor(t *testing.T) {
	assert.True(t, utils.Xor(true, false))
	assert.False(t, utils.Xor(false, false))
	assert.True(t, utils.Xor(false, true))
	assert.False(t, utils.Xor(true, true))
}

func TestGetFunctionName(t *testing.T) {
	fn := func() {}
	result := utils.GetFunctionName(fn)
	assert.Equal(t, "func1", result)

	invalidFn := "invalid function"
	result = utils.GetFunctionName(invalidFn)
	assert.Equal(t, "", result)
}
