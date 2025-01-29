package main

import (
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/DroppedHard/SWIFT-service/config"
)

func main() {
	var filePath string
	flag.StringVar(&filePath, "source", config.Envs.MigrationFilePath, "Path to the JSON file containing migration data")
	flag.Parse()

	rdb := connectToRedis()

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

	startMigration(data, rdb)
}
