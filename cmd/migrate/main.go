package main

import (
	"context"
	"encoding/csv"
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/DroppedHard/SWIFT-service/config"
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
		isHeadquarter := strings.HasSuffix(swiftCode, "XXX")
		countryISO2 := record[0]
		bankName := record[2]
		countryName := record[4]

		redisEntry := RedisData{
			Key: swiftCode,
			Fields: map[string]string{
				"swiftCode":     swiftCode,
				"address":       address,
				"isHeadquarter": fmt.Sprintf("%t", isHeadquarter),
				"countryISO2":   countryISO2,
				"bankName":      bankName,
				"countryName":   countryName,
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

func main() {

	var filePath string
	flag.StringVar(&filePath, "source", "./cmd/migrate/migrations/initial_data.csv", "Path to the JSON file containing migration data")
	flag.Parse()

	file, err := os.Open(filePath)
	if err != nil {
		fmt.Printf("Failed to open the file: %v\n", err)
		return
	}
	defer file.Close()

	var data []RedisData
	switch {
	case strings.HasSuffix(filePath, ".json"):
		data, err = parseJSON(file)
	case strings.HasSuffix(filePath, ".csv"):
		data, err = parseCSV(file)
	default:
		fmt.Println("Unsupported file format. Please provide a JSON or CSV file.")
		return
	}

	if err != nil {
		fmt.Printf("Failed to decode file: %v\n", err)
		return
	}

	rdb := redis.NewClient(&redis.Options{
		Addr:     config.Envs.DBAddress,
		Password: config.Envs.DBPassword,
		DB:       config.Envs.DBNum,
	})

	startMigration(data, rdb)
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
