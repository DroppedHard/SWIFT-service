package swiftCode

import (
	"context"
	"fmt"
	"log"

	"github.com/DroppedHard/SWIFT-service/types"
	"github.com/redis/go-redis/v9"
)

type Store struct {
	client redis.Client
}

func NewStore(client *redis.Client) *Store {
	return &Store{client: *client}
}

func (s *Store) DoesSwiftCodeExist(ctx context.Context, swiftCode string) (int64, error) {
	out, err := s.client.Exists(ctx, swiftCode).Result()
	if err != nil {
		return -1, err
	}
	return out, nil
}

func (s *Store) AddBankData(ctx context.Context, data types.BankDataDetails) error {
	hashData := map[string]interface{}{
		"address":       data.Address,
		"bankName":      data.BankName,
		"countryISO2":   data.CountryISO2,
		"countryName":   data.CountryName,
		"isHeadquarter": data.IsHeadquarter,
		"swiftCode":     data.SwiftCode,
	}
	if _, err := s.client.HSet(ctx, data.SwiftCode, hashData).Result(); err != nil {
		return fmt.Errorf("failed to store data for key %s: %w", data.SwiftCode, err)
	}
	return nil
}

// DeleteBankData implements types.BankDataStore.
func (s *Store) DeleteBankData(ctx context.Context, swiftCode string) error {
	panic("unimplemented")
}

// GetBanksDataByCountryCode implements types.BankDataStore.
func (s *Store) GetBanksDataByCountryCode(ctx context.Context, countryCode string) ([]types.BankDataCore, error) {
	panic("unimplemented")
}

// GetBranchesDataByHqSwiftCode implements types.BankDataStore.
func (s *Store) GetBranchesDataByHqSwiftCode(ctx context.Context, swiftCode string) ([]types.BankDataCore, error) {
	panic("unimplemented")
}

func (s *Store) GetBankDetailsBySwiftCode(ctx context.Context, swiftCode string) (*types.BankDataDetails, error) {
	rows, err := s.client.HGetAll(ctx, swiftCode).Result()
	if err != nil {
		return nil, err
	}
	log.Println("TODO", fmt.Sprint(rows))
	return nil, nil
}
