package store_test

import (
	"github.com/DroppedHard/SWIFT-service/types"
	"github.com/DroppedHard/SWIFT-service/utils"
)

type RedisData struct {
	Key    string            `json:"key"`
	Fields map[string]string `json:"fields"`
}

var TestRedisData = []RedisData{
	{
		Key: "ALBPPLPWXXX",
		Fields: map[string]string{
			utils.RedisHashAddress:     "testBank HQ Address",
			utils.RedisHashCountryISO2: "PL",
			utils.RedisHashSwiftCode:   "ALBPPLPWXXX",
			utils.RedisHashBankName:    "testBank",
		},
	},
	{
		Key: "ALBPPLPWABC",
		Fields: map[string]string{
			utils.RedisHashAddress:     "testBank Branch 1 Address",
			utils.RedisHashCountryISO2: "PL",
			utils.RedisHashSwiftCode:   "ALBPPLPWABC",
			utils.RedisHashBankName:    "testBank",
		},
	},
	{
		Key: "ALBPPLPWDEF",
		Fields: map[string]string{
			utils.RedisHashAddress:     "testBank Branch 2 Address",
			utils.RedisHashCountryISO2: "PL",
			utils.RedisHashSwiftCode:   "ALBPPLPWDEF",
			utils.RedisHashBankName:    "testBank",
		},
	},
}

var NonexistentSwiftCodes = []string{
	"ALBPPLPWGHI",
	"ALBPPLPWJKL",
	"ALBPPLPW123",
}

var NewBankData = []types.BankDataDetails{
	{
		BankDataCore: types.BankDataCore{
			Address:       "testBank Branch 3 Address",
			BankName:      "testBank",
			CountryIso2:   "PL",
			IsHeadquarter: false,
			SwiftCode:     "ALBPPLPWMNO",
		},
		CountryName: "Poland",
	},
	{
		BankDataCore: types.BankDataCore{
			Address:       "newBank HQ Address",
			BankName:      "newBank",
			CountryIso2:   "PL",
			IsHeadquarter: false,
			SwiftCode:     "ABCDPLPWMNO",
		},
		CountryName: "Poland",
	},
}
