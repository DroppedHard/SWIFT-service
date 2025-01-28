package config

import (
	"fmt"
	"os"
	"strconv"

	"github.com/lpernett/godotenv"
)

type Config struct {
	PublicHost     string
	Port           string
	DBPassword     string
	DBAddress      string
	DBNum          int
	DBPoolSize     int
	DBMinIdleConns int
}

var Envs = initConfig()

func initConfig() Config {
	godotenv.Load()
	return Config{
		PublicHost:     getEnv("PUBLIC_HOST", "localhost"),
		Port:           getEnv("PORT", ":8080"),
		DBPassword:     getEnv("DB_PASSWORD", ""),
		DBAddress:      fmt.Sprintf("%s:%s", getEnv("DB_HOST", "localhost"), getEnv("DB_PORT", "6379")),
		DBNum:          getEnvInt("DB_NUM", 0),
		DBPoolSize:     getEnvInt("DB_POOL_SIZE", 20),
		DBMinIdleConns: getEnvInt("DB_MIN_IDLE_CONNS", 1),
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
