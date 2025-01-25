package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	APP_ID  string
	API_KEY string
	TOKEN   string
}

func New() *Config {
	err := godotenv.Load()
	if err != nil {
		log.Default().Printf("[config] error loading .env file: %v\n", err)
	}

	return &Config{
		APP_ID:  getEnv("APP_ID", ""),
		API_KEY: getEnv("API_KEY", ""),
		TOKEN:   getEnv("TOKEN", ""),
	}
}

func getEnv(key, defValue string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}

	return defValue
}
