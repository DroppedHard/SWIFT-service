package config

import (
	"os"
	"strconv"

	"github.com/lpernett/godotenv"
)

type Config struct {
	PublicHost        string
	Port              string
	DBPassword        string
	DBHost            string
	DBPort            string
	DBNum             int
	DBPoolSize        int
	DBMinIdleConns    int
	MigrationFilePath string
}

var defaultConfig = Config{
	PublicHost:        "localhost",
	Port:              ":8080",
	DBPassword:        "",
	DBHost:            "localhost",
	DBPort:            "6379",
	DBNum:             0,
	DBPoolSize:        20,
	DBMinIdleConns:    1,
	MigrationFilePath: "./cmd/migrate/migrations/initial_data.csv",
}

var Envs = initConfig()

func initConfig() Config {
	godotenv.Load()
	return Config{
		PublicHost:        getEnv("PUBLIC_HOST", defaultConfig.PublicHost),
		Port:              getEnv("PORT", defaultConfig.Port),
		DBPassword:        getEnv("DB_PASSWORD", defaultConfig.DBPassword),
		DBHost:            getEnv("DB_HOST", defaultConfig.DBHost),
		DBPort:            getEnv("DB_PORT", defaultConfig.DBPort),
		DBNum:             getEnvInt("DB_NUM", defaultConfig.DBNum),
		DBPoolSize:        getEnvInt("DB_POOL_SIZE", defaultConfig.DBPoolSize),
		DBMinIdleConns:    getEnvInt("DB_MIN_IDLE_CONNS", defaultConfig.DBMinIdleConns),
		MigrationFilePath: getEnv("MIGRATION_FILE", defaultConfig.MigrationFilePath),
	}
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}

func getEnvInt(key string, fallback int) int {
	value, ok := os.LookupEnv(key)
	if !ok {
		return fallback
	}
	intValue, err := strconv.Atoi(value)
	if err != nil {
		return fallback
	}
	return intValue
}
