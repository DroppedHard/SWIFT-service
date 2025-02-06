package main

import (
	"context"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/DroppedHard/SWIFT-service/config"
	"github.com/DroppedHard/SWIFT-service/utils"
	"github.com/redis/go-redis/v9"
)

type RedisData struct {
	Key    string            `json:"key"`
	Fields map[string]string `json:"fields"`
}

func parseCSV(file *os.File) ([]RedisData, error) {
	var data []RedisData
	reader := csv.NewReader(file)
	reader.Comma = ';'

	records, err := reader.ReadAll()
	if err != nil {
		return nil, err
	}

	for _, record := range records[1:] {
		if len(record) < 5 {
			continue
		}
		swiftCode := record[1]
		address := record[3]
		isHeadquarter := isHeadquarterParser(swiftCode)
		countryIso2 := strings.ToUpper(record[0])
		bankName := record[2]
		countryName := strings.ToUpper(record[4])

		redisEntry := RedisData{
			Key: swiftCode,
			Fields: map[string]string{
				utils.RedisHashSwiftCode:     swiftCode,
				utils.RedisHashAddress:       address,
				utils.RedisHashIsHeadquarter: isHeadquarter,
				utils.RedisHashCountryISO2:   countryIso2,
				utils.RedisHashBankName:      bankName,
				utils.RedisHashCountryName:   countryName,
			},
		}
		data = append(data, redisEntry)
	}
	return data, nil
}

func parseJSON(file *os.File) ([]RedisData, error) {
	var data []RedisData
	if err := json.NewDecoder(file).Decode(&data); err != nil {
		return nil, err
	}
	return data, nil
}

func isHeadquarterParser(swiftCode string) string {
	if strings.HasSuffix(swiftCode, utils.BranchSuffix) {
		return utils.RedisStoreTrue
	}
	return utils.RedisStoreFalse
}

func connectToRedis() *redis.Client {
	var rdb *redis.Client
	retryCount := 10
	delay := 2 * time.Second

	for i := 0; i < retryCount; i++ {
		fmt.Printf("Redis connection attempt %d...\n", i)

		rdb = redis.NewClient(&redis.Options{
			Addr:     fmt.Sprintf("%s:%s", config.Envs.DBHost, config.Envs.DBPort),
			Password: config.Envs.DBPassword,
			DB:       config.Envs.DBNum,
		})

		ctx, cancel := context.WithTimeout(context.Background(), delay)
		defer cancel()

		_, err := rdb.Ping(ctx).Result()
		if err == nil {
			fmt.Println("Connected to Redis!")
			return rdb
		}

		fmt.Printf("Redis connection failed: %v. Retrying in %v...\n", err, delay)
		time.Sleep(delay)
	}

	fmt.Println("Failed to connect to Redis after multiple attempts. Exiting.")
	os.Exit(1)
	return nil
}

func startMigration(data []RedisData, rdb *redis.Client) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var wg sync.WaitGroup
	for _, entry := range data {
		wg.Add(1)
		go func(entry RedisData) {
			defer wg.Done()

			if _, err := rdb.HSet(ctx, entry.Key, entry.Fields).Result(); err != nil {
				fmt.Printf("Failed to populate key %s: %v\n", entry.Key, err)
				return
			}
			fmt.Printf("Successfully populated key: %s\n", entry.Key)
		}(entry)
	}
	wg.Wait()
	fmt.Println("Migration completed successfully.")
}
