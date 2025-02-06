package types

import "context"

type BankDataCore struct {
	Address       string `json:"address" validate:"required"`
	BankName      string `json:"bankName" validate:"required"`
	CountryIso2   string `json:"countryISO2" validate:"required,countryISO2"`
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
	DoesSwiftCodeExist(ctx context.Context, swiftCode string) (int64, error)
	SaveBankData(ctx context.Context, data BankDataDetails) error
	DeleteBankData(ctx context.Context, swiftCode string) error
	FindBanksDataByCountryCode(ctx context.Context, countryCode string) ([]BankDataCore, error)
	FindBranchesDataByHqSwiftCode(ctx context.Context, swiftCode string) ([]BankDataCore, error)
	FindBankDetailsBySwiftCode(ctx context.Context, swiftCode string) (*BankDataDetails, error)
	Ping(ctx context.Context) error
}
