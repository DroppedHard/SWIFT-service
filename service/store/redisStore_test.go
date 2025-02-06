package store_test

import (
	"context"
	"fmt"
	"strings"
	"testing"

	"github.com/DroppedHard/SWIFT-service/config"
	"github.com/DroppedHard/SWIFT-service/db"
	"github.com/DroppedHard/SWIFT-service/service/store"
	"github.com/DroppedHard/SWIFT-service/utils"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/suite"
)

type RedisStoreTestSuite struct {
	suite.Suite
	client redis.Client
	store  *store.RedisStore
}

func TestRoutesSuite(t *testing.T) {
	suite.Run(t, &RedisStoreTestSuite{})
}

func (suite *RedisStoreTestSuite) prepareData() {
	ctx := context.Background()
	for _, entry := range TestRedisData {
		if _, err := suite.client.HSet(ctx, entry.Key, entry.Fields).Result(); err != nil {
			fmt.Printf("Failed to populate key %s: %v\n", entry.Key, err)
			return
		}
		fmt.Printf("Successfully populated key: %s\n", entry.Key)
	}
}
func (suite *RedisStoreTestSuite) deleteData() {
	ctx := context.Background()
	for _, entry := range TestRedisData {
		if _, err := suite.client.Del(ctx, entry.Key).Result(); err != nil {
			fmt.Printf("Failed to delete key %s: %v\n", entry.Key, err)
			return
		}
		fmt.Printf("Successfully deleted key: %s\n", entry.Key)
	}
}

func (suite *RedisStoreTestSuite) countBranches(swiftCode string) int {
	counter := 0
	for _, entry := range TestRedisData {
		if !strings.HasSuffix(entry.Key, utils.BranchSuffix) && entry.Key[:8] == swiftCode[:8] {
			counter++
		}
	}
	return counter
}

func (suite *RedisStoreTestSuite) countCountryMatching(countryIso2Code string) int {
	counter := 0
	for _, entry := range TestRedisData {
		if entry.Fields[utils.RedisHashCountryISO2] == countryIso2Code {
			counter++
		}
	}
	return counter
}

func (suite *RedisStoreTestSuite) SetupSuite() {
	rdb := db.NewRedisStorage(&redis.Options{
		Addr:         fmt.Sprintf("%s:%s", config.Envs.DBHost, config.Envs.DBPort),
		Password:     config.Envs.DBPassword,
		DB:           config.Envs.DBTestNum,
		PoolSize:     config.Envs.DBPoolSize,
		MinIdleConns: config.Envs.DBMinIdleConns,
	})
	suite.client = *rdb
	suite.store = store.NewStore(&suite.client)
	suite.prepareData()
}

func (suite *RedisStoreTestSuite) TearDownSuite() {
	suite.deleteData()
	suite.client.Close()
}

func (suite *RedisStoreTestSuite) TestDoesSwiftCodeExist() {
	ctx := context.Background()
	suite.Run("Positive Cases", func() {
		for _, entry := range TestRedisData {
			suite.Run(entry.Key+" should exist", func() {
				count, err := suite.store.DoesSwiftCodeExist(ctx, entry.Key)
				suite.NoError(err)
				suite.Equal(count, int64(1))
			})
		}
	})
	for _, code := range NonexistentSwiftCodes {
		suite.Run("Negative case - "+code+"should not exist", func() {
			count, err := suite.store.DoesSwiftCodeExist(ctx, code)
			suite.NoError(err)
			suite.Equal(count, int64(0))
		})
	}
	suite.Run("Negative Cases - Invalid Contexts", func() {
		swiftCode := TestRedisData[0].Key
		suite.Run("Canceled Context", func() {
			canceledCtx, cancel := context.WithCancel(context.Background())
			cancel()

			_, err := suite.store.DoesSwiftCodeExist(canceledCtx, swiftCode)
			suite.Error(err)
			suite.Contains(err.Error(), "context canceled")
		})

		suite.Run("Expired Context", func() {
			expiredCtx, cancel := context.WithTimeout(context.Background(), 0)
			defer cancel()

			_, err := suite.store.DoesSwiftCodeExist(expiredCtx, swiftCode)
			suite.Error(err)
			suite.Contains(err.Error(), "context deadline exceeded")
		})
	})
}

func (suite *RedisStoreTestSuite) TestSaveBankDataAndFindBankDetailsBySwiftCode() {
	ctx := context.Background()
	suite.Run("Positive Cases", func() {
		pipe := suite.client.Pipeline()
		for _, entry := range NewBankData {
			suite.Run("Save "+entry.SwiftCode+" bank data", func() {
				err := suite.store.SaveBankData(ctx, entry)
				defer pipe.Del(ctx, entry.SwiftCode).Result()
				suite.NoError(err)

				redisData, err := pipe.HGetAll(ctx, entry.SwiftCode).Result()
				suite.NoError(err)
				fmt.Printf("Redis Data for %s: %+v\n", entry.SwiftCode, redisData)

				data, err := suite.store.FindBankDetailsBySwiftCode(ctx, entry.SwiftCode)
				suite.NoError(err)
				suite.Equal(&entry, data)
			})
		}
	})
	suite.Run("Negative Cases - Invalid Contexts", func() {
		exampleData := NewBankData[0]
		defer suite.client.Del(ctx, exampleData.SwiftCode)
		suite.Run("Canceled Context", func() {
			canceledCtx, cancel := context.WithCancel(context.Background())
			cancel()

			err := suite.store.SaveBankData(canceledCtx, exampleData)
			suite.Error(err)
			suite.Contains(err.Error(), "context canceled")
		})

		suite.Run("Expired Context", func() {
			expiredCtx, cancel := context.WithTimeout(context.Background(), 0)
			defer cancel()

			err := suite.store.SaveBankData(expiredCtx, exampleData)
			suite.Error(err)
			suite.Contains(err.Error(), "context deadline exceeded")
		})
	})
}

func (suite *RedisStoreTestSuite) TestDeleteBankData() {
	ctx := context.Background()

	suite.Run("Delete bank data", func() {
		for _, entry := range NewBankData {
			err := suite.store.SaveBankData(ctx, entry)
			suite.NoError(err)
		}

		swiftCodes := []string{NewBankData[0].SwiftCode, NewBankData[1].SwiftCode}
		for _, code := range swiftCodes {
			err := suite.store.DeleteBankData(ctx, code)
			suite.NoError(err)
		}

		for _, swiftCode := range swiftCodes {
			data, err := suite.store.FindBankDetailsBySwiftCode(ctx, swiftCode)
			suite.NoError(err)
			suite.Nil(data)
		}
	})

	suite.Run("Negative Cases - Invalid Contexts", func() {
		exampleData := NewBankData[0]
		suite.Run("Canceled Context", func() {
			canceledCtx, cancel := context.WithCancel(context.Background())
			cancel()

			err := suite.store.DeleteBankData(canceledCtx, exampleData.SwiftCode)
			suite.Error(err)
			suite.Contains(err.Error(), "context canceled")
		})

		suite.Run("Expired Context", func() {
			expiredCtx, cancel := context.WithTimeout(context.Background(), 0)
			defer cancel()

			err := suite.store.DeleteBankData(expiredCtx, exampleData.SwiftCode)
			suite.Error(err)
			suite.Contains(err.Error(), "context deadline exceeded")
		})
	})
}

func (suite *RedisStoreTestSuite) TestFindBanksDataByCountryCode() {
	ctx := context.Background()
	countryCode := "PL"

	suite.Run("Positive Case", func() {
		banksData, err := suite.store.FindBanksDataByCountryCode(ctx, countryCode)
		suite.NoError(err)
		suite.Len(banksData, suite.countCountryMatching(countryCode))

		testDataMap := make(map[string]RedisData)
		for _, data := range TestRedisData {
			testDataMap[data.Key] = data
		}

		for _, bank := range banksData {
			expectedData, exists := testDataMap[bank.SwiftCode]
			suite.True(exists, "Expected data for SwiftCode %s does not exist", bank.SwiftCode)

			suite.Equal(expectedData.Key, bank.SwiftCode)
			suite.Equal(expectedData.Fields[utils.RedisHashAddress], bank.Address)
			suite.Equal(expectedData.Fields[utils.RedisHashBankName], bank.BankName)
			suite.Equal(expectedData.Fields[utils.RedisHashCountryISO2], bank.CountryIso2)
			suite.Equal(expectedData.Fields[utils.RedisHashSwiftCode], bank.SwiftCode)
		}
	})

	suite.Run("Cancelled context", func() {
		ctx, cancel := context.WithCancel(context.Background())
		cancel()

		banksData, err := suite.store.FindBanksDataByCountryCode(ctx, countryCode)
		suite.Error(err)
		suite.Nil(banksData)
	})

	suite.Run("Timeout Context", func() {
		ctx, cancel := context.WithTimeout(context.Background(), 0)
		defer cancel()

		banksData, err := suite.store.FindBanksDataByCountryCode(ctx, countryCode)
		suite.Error(err)
		suite.Nil(banksData)
	})
}

func (suite *RedisStoreTestSuite) TestFindBranchesDataByHqSwiftCode() {
	ctx := context.Background()
	hqSwiftCode := TestRedisData[0].Key

	suite.Run("Positive Case", func() {
		banksData, err := suite.store.FindBranchesDataByHqSwiftCode(ctx, hqSwiftCode)
		suite.NoError(err)
		suite.Len(banksData, suite.countBranches(hqSwiftCode))

		testDataMap := make(map[string]RedisData)
		for _, data := range TestRedisData {
			testDataMap[data.Key] = data
		}

		for _, bank := range banksData {
			expectedData, exists := testDataMap[bank.SwiftCode]
			suite.True(exists, "Expected data for SwiftCode %s does not exist", bank.SwiftCode)

			suite.Equal(expectedData.Key, bank.SwiftCode)
			suite.Equal(expectedData.Fields[utils.RedisHashAddress], bank.Address)
			suite.Equal(expectedData.Fields[utils.RedisHashBankName], bank.BankName)
			suite.Equal(expectedData.Fields[utils.RedisHashCountryISO2], bank.CountryIso2)
			suite.Equal(expectedData.Fields[utils.RedisHashSwiftCode], bank.SwiftCode)
		}
	})

	suite.Run("Cancelled context", func() {
		ctx, cancel := context.WithCancel(context.Background())
		cancel()

		banksData, err := suite.store.FindBranchesDataByHqSwiftCode(ctx, hqSwiftCode)
		suite.Error(err)
		suite.Nil(banksData)
	})

	suite.Run("Timeout Context", func() {
		ctx, cancel := context.WithTimeout(context.Background(), 0)
		defer cancel()

		banksData, err := suite.store.FindBranchesDataByHqSwiftCode(ctx, hqSwiftCode)
		suite.Error(err)
		suite.Nil(banksData)
	})
}

func (suite *RedisStoreTestSuite) TestFindBankDetailsBySwiftCode() {
	ctx := context.Background()

	suite.Run("Positive Cases", func() {
		expectedData := TestRedisData[0]

		bankData, err := suite.store.FindBankDetailsBySwiftCode(ctx, expectedData.Key)
		suite.NoError(err)
		suite.NotNil(bankData)

		suite.Equal(expectedData.Key, bankData.SwiftCode)
		suite.Equal(expectedData.Fields[utils.RedisHashAddress], bankData.Address)
		suite.Equal(expectedData.Fields[utils.RedisHashBankName], bankData.BankName)
		suite.Equal(expectedData.Fields[utils.RedisHashCountryISO2], bankData.CountryIso2)
		suite.Equal(expectedData.Fields[utils.RedisHashSwiftCode], bankData.SwiftCode)
		suite.Equal(expectedData.Fields[utils.RedisHashCountryName], bankData.CountryName)
		suite.Equal(expectedData.Fields[utils.RedisHashIsHeadquarter] == utils.RedisStoreTrue, bankData.IsHeadquarter)
	})

	suite.Run("Negative Case - No Data Found", func() {
		swiftCode := "NONEXISTENTCODE"

		bankData, err := suite.store.FindBankDetailsBySwiftCode(ctx, swiftCode)
		suite.NoError(err)
		suite.Nil(bankData)
	})

	suite.Run("Cancelled context", func() {
		ctx, cancel := context.WithCancel(context.Background())
		cancel()

		banksData, err := suite.store.FindBankDetailsBySwiftCode(ctx, TestRedisData[0].Key)
		suite.Error(err)
		suite.Nil(banksData)
	})

	suite.Run("Timeout Context", func() {
		ctx, cancel := context.WithTimeout(context.Background(), 0)
		defer cancel()

		banksData, err := suite.store.FindBankDetailsBySwiftCode(ctx, TestRedisData[0].Key)
		suite.Error(err)
		suite.Nil(banksData)
	})
}
