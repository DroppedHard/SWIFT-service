package types

import "context"

type BankDataCore struct {
	Address       string `json:"address" validate:"required"`
	BankName      string `json:"bankName" validate:"required"`
	CountryISO2   string `json:"countryISO2" validate:"required,len=2,countryISO2"` // potentially add a custom validator that verifies if it is a country code.
	IsHeadquarter bool   `json:"isHeadquarters" validate:"required"`
	SwiftCode     string `json:"swiftCode" validate:"required,len=11,swiftCode"`
}

type BankDataDetails struct {
	BankDataCore
	CountryName string `json:"countryName" validate:"required"`
}

type BankDataStore interface {
	GetBankDetailsBySwiftCode(ctx context.Context, swiftCode string) (*BankDataDetails, error)
	GetBranchesDataByHqSwiftCode(ctx context.Context, swiftCode string) ([]BankDataCore, error)
	GetBanksDataByCountryCode(ctx context.Context, countryCode string) ([]BankDataCore, error)
	AddBankData(ctx context.Context, data BankDataDetails) error
	DeleteBankData(ctx context.Context, swiftCode string) error // after I get info about how it should work.
	DoesSwiftCodeExist(ctx context.Context, swiftCode string) (int64, error)
}
