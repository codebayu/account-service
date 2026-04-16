package config

import (
	"fmt"
	"os"
)

type Config struct {
	AppPort     string
	DBHost      string
	DBPort      string
	DBUser      string
	DBPassword  string
	DBName      string
	DBSSLMode   string
	JWTSecret   string
}

func LoadConfig() (*Config, error) {
	config := &Config{
		AppPort:    getEnv("APP_PORT", "8080"),
		DBHost:     os.Getenv("DB_HOST"),
		DBPort:     getEnv("DB_PORT", "5432"),
		DBUser:     os.Getenv("DB_USER"),
		DBPassword: os.Getenv("DB_PASSWORD"),
		DBName:     os.Getenv("DB_NAME"),
		DBSSLMode:  getEnv("DB_SSLMODE", "disable"),
		JWTSecret:  getEnv("JWT_SECRET", "secret"),
	}

	if config.DBHost == "" || config.DBUser == "" || config.DBName == "" {
		return nil, fmt.Errorf("❌ database environment variables not set properly")
	}

	return config, nil
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}
