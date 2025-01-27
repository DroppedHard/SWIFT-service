package swiftCode

import (
	"context"
	"fmt"
	"strings"
	"sync"

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

func (s *Store) GetBranchesDataByHqSwiftCode(ctx context.Context, swiftCode string) ([]types.BankDataCore, error) {
	branchPrefix := swiftCode[:8] + "???"
	branchKeys, err := s.client.Keys(ctx, branchPrefix).Result()
	if err != nil {
		return nil, fmt.Errorf("failed to fetch branches for SWIFT code prefix %s: %w", branchPrefix, err)
	}

	var (
		branches []types.BankDataCore
		mu       sync.Mutex       // To safely append to branches slice
		wg       sync.WaitGroup   // To wait for all goroutines to complete
		errs     []string         // To collect errors from goroutines
		errMu    sync.Mutex       // To safely append to the errs slice
	)

	for _, branchKey := range branchKeys {
		if branchKey == swiftCode {
			continue
		}
		wg.Add(1)
		go func(branchKey string) {
			defer wg.Done()
			
			branchFields, err := s.GetBankDetailsBySwiftCode(ctx, branchKey)
			if err != nil {
				errMu.Lock()
				errs = append(errs, fmt.Sprintf("failed to fetch branch data for key %s: %v", branchKey, err))
				errMu.Unlock()
				return
			}

			mu.Lock()
			branches = append(branches, types.BankDataCore{
				Address:       branchFields.Address,
				BankName:      branchFields.BankName,
				CountryISO2:   branchFields.CountryISO2,
				IsHeadquarter: branchFields.IsHeadquarter,
				SwiftCode:     branchFields.SwiftCode,
			})
			mu.Unlock()
		}(branchKey)

	}

	wg.Wait()

	if len(errs) > 0 {
		return branches, fmt.Errorf("encountered errors: %s", strings.Join(errs, "; "))
	}
	
	return branches, nil
}

func (s *Store) GetBankDetailsBySwiftCode(ctx context.Context, swiftCode string) (*types.BankDataDetails, error) {
	rows, err := s.client.HGetAll(ctx, swiftCode).Result()
	if err != nil {
		return nil, fmt.Errorf("failed to fetch data from Redis for key %s: %v", swiftCode, err)
	}
	if len(rows) == 0 {
		return nil, nil
	}

	bankDetails := &types.BankDataDetails{
		BankDataCore: types.BankDataCore{
			Address:       rows["address"],
			BankName:      rows["bankName"],
			CountryISO2:   rows["countryISO2"],
			IsHeadquarter: rows["isHeadquarter"] == "1",
			SwiftCode:     rows["swiftCode"],
		},
		CountryName: rows["countryName"],
	}

	return bankDetails, nil
}
