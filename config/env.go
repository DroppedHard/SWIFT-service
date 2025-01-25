package config

import (
	"fmt"
	"os"

	"github.com/lpernett/godotenv"
)

type Config struct {
	PublicHost string
	Port       string
	DBPassword string
	DBAddress  string
	DBNum      int
}

var Envs = initConfig()

func initConfig() Config {
	godotenv.Load()
	// if dbNum, err := strconv.Atoi(getEnv("DB_NUMBER", "0")); err != nil {
	return Config{
		PublicHost: getEnv("PUBLIC_HOST", "http://localhost"),
		Port:       getEnv("PORT", ":8080"),
		DBPassword: getEnv("DB_PASSWORD", ""),
		DBAddress:  fmt.Sprintf("%s:%s", getEnv("DB_HOST", "localhost"), getEnv("DB_PORT", "6379")),
		DBNum:      0,
	}
	// } else {
	// 	log.Fatal(err)
	// 	return Config{
	// 		PublicHost: getEnv("PUBLIC_HOST", "http://localhost"),
	// 		Port:       getEnv("PORT", "8080"),
	// 		DBPassword: getEnv("DB_PASSWORD", ""),
	// 		DBAddress:  fmt.Sprintf("%s:%s", getEnv("DB_HOST", "localhost"), getEnv("DB_PORT", "6379")),
	// 	}
	// }
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}
