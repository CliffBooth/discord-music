package config

import (
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	APP_ID    string
	API_KEY   string
	TOKEN     string
	DEV       bool
	LOG_LEVEL string
}

func Load() *Config {
	err := godotenv.Load()
	if err != nil {
		log.Default().Printf("[config] ERROR LOADING .env FILE: %v\n", err)
	}

	return &Config{
		APP_ID:    getEnv("APP_ID", ""),
		API_KEY:   getEnv("API_KEY", ""),
		TOKEN:     getEnv("TOKEN", ""),
		DEV:       getBool("DEV", false),
		LOG_LEVEL: getEnv("LOG_LEVEL", ""),
	}
}

func getEnv(key, defValue string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}

	return defValue
}

func getBool(key string, defValue bool) bool {
	value, ok := os.LookupEnv(key)
	if !ok {
		return defValue
	}

	res, err := strconv.ParseBool(value)
	if err != nil {
		log.Default().Printf("[config] CANNOT CONVERT %s=%s TO BOOL, USING DEFAULT %v\n", key, value, defValue)
		return defValue
	}

	return res
}
