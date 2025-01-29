package swiftCode

import (
	"context"
	"fmt"
	"strings"
	"sync"

	"github.com/DroppedHard/SWIFT-service/types"
	"github.com/DroppedHard/SWIFT-service/utils"
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
	hashData := map[string]interface{}{ // TODO - maybe reduce this?
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

func (s *Store) DeleteBankData(ctx context.Context, swiftCode string) error {
	_, err := s.client.Del(ctx, swiftCode).Result()
	if err != nil {
		return fmt.Errorf("failed to delete data for SWIFT code %s: %w", swiftCode, err)
	}

	return nil
}

func (s *Store) GetBanksDataByCountryCode(ctx context.Context, countryCode string) ([]types.BankDataCore, error) {
	keys, err := s.client.Keys(ctx, utils.CountryCodeRegex(countryCode)).Result()
	if err != nil {
		return nil, fmt.Errorf("failed to fetch keys for country code %s: %w", countryCode, err)
	}

	return s.getBankDetailsByCodesConcurrently(ctx, keys, "")
}

func (s *Store) GetBranchesDataByHqSwiftCode(ctx context.Context, swiftCode string) ([]types.BankDataCore, error) {
	branchKeys, err := s.client.Keys(ctx, utils.BranchRegex(swiftCode)).Result()
	if err != nil {
		return nil, fmt.Errorf("failed to fetch branches for SWIFT code %s: %w", swiftCode, err)
	}

	return s.getBankDetailsByCodesConcurrently(ctx, branchKeys, swiftCode)
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
			IsHeadquarter: rows["isHeadquarter"] == "true", // TODO - check if this is necessary -> boolean mapping string incident
			SwiftCode:     rows["swiftCode"],
		},
		CountryName: rows["countryName"],
	}

	return bankDetails, nil
}

type bankDataChanResult struct {
	bankData types.BankDataCore
	err      error
}

func (s *Store) getBankDetailsByCodesConcurrently(ctx context.Context, branchKeys []string, currentSwiftCode string) (branches []types.BankDataCore, aggregatedErr error) {
	resultsChan := make(chan bankDataChanResult, len(branchKeys))
	var (
		wg   sync.WaitGroup
		errs []string
	)

	for _, branchKey := range branchKeys {
		if branchKey == currentSwiftCode {
			continue
		}
		wg.Add(1)
		go s.fetchBankDetails(ctx, branchKey, resultsChan, &wg)
	}

	go func() {
		wg.Wait()
		close(resultsChan)
	}()

	for result := range resultsChan {
		if result.err != nil {
			errs = append(errs, result.err.Error())
		} else {
			branches = append(branches, result.bankData)
		}
	}

	if len(errs) > 0 {
		aggregatedErr = fmt.Errorf("encountered errors: %s", strings.Join(errs, "; "))
	}
	return
}

func (s *Store) fetchBankDetails(ctx context.Context, branchKey string, resultsChan chan<- bankDataChanResult, wg *sync.WaitGroup) {
	defer wg.Done()

	branchFields, err := s.GetBankDetailsBySwiftCode(ctx, branchKey)
	if err != nil {
		resultsChan <- bankDataChanResult{err: fmt.Errorf("failed to fetch branch data for key %s: %w", branchKey, err)}
	} else {
		resultsChan <- bankDataChanResult{
			bankData: types.BankDataCore{
				Address:       branchFields.Address,
				BankName:      branchFields.BankName,
				CountryISO2:   branchFields.CountryISO2,
				IsHeadquarter: branchFields.IsHeadquarter,
				SwiftCode:     branchFields.SwiftCode,
			},
		}
	}
}
