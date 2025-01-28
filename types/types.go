package types

import "context"

type BankDataCore struct {
	Address       string `json:"address" validate:"required"`
	BankName      string `json:"bankName" validate:"required"`
	CountryISO2   string `json:"countryISO2" validate:"required,countryISO2"` // potentially add a custom validator that verifies if it is a country code.
	IsHeadquarter bool   `json:"isHeadquarter" validate:"boolRequired"`
	SwiftCode     string `json:"swiftCode" validate:"required,len=11,swiftCode"`
}

type BankDataDetails struct {
	BankDataCore
	CountryName string `json:"countryName" validate:"required"`
}

type BankHeadquatersResponse struct {
	BankDataDetails
	Branches []BankDataCore `json:"branches"`
}

type CountrySwiftCodesResponse struct {
	CountryIso2 string         `json:"countryISO2"`
	CountryName string         `json:"countryName"`
	SwiftCodes  []BankDataCore `json:"swiftCodes"`
}

type BankDataStore interface {
	GetBankDetailsBySwiftCode(ctx context.Context, swiftCode string) (*BankDataDetails, error)
	GetBranchesDataByHqSwiftCode(ctx context.Context, swiftCode string) ([]BankDataCore, error)
	GetBanksDataByCountryCode(ctx context.Context, countryCode string) ([]BankDataCore, error)
	AddBankData(ctx context.Context, data BankDataDetails) error
	DeleteBankData(ctx context.Context, swiftCode string) error
	DoesSwiftCodeExist(ctx context.Context, swiftCode string) (int64, error)
}
