package config

import (
	"fmt"
	"os"
	"strings"
)

type Config struct {
	AppPort        string
	DBHost         string
	DBPort         string
	DBUser         string
	DBPassword     string
	DBName         string
	DBSSLMode      string
	JWTSecret      string
	APIKey         string
	APISecret      string
	ChannelID      string
	AllowedOrigins []string
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
		JWTSecret:  getEnv("JWT_SECRET", ""),
		APIKey:     getEnv("API_KEY", ""),
		APISecret:  getEnv("API_SECRET", ""),
		ChannelID:  getEnv("CHANNEL_ID", ""),
	}

	if config.DBHost == "" || config.DBUser == "" || config.DBName == "" {
		return nil, fmt.Errorf("❌ database environment variables not set properly")
	}

	if config.APIKey == "" || config.APISecret == "" {
		return nil, fmt.Errorf("❌ API_KEY or API_SECRET not set properly")
	}

	allowedOrigins := os.Getenv("ALLOWED_ORIGINS")
	if allowedOrigins == "" {
		return nil, fmt.Errorf("❌ ALLOWED_ORIGINS not set properly")
	}
	config.AllowedOrigins = strings.Split(allowedOrigins, ",")

	return config, nil
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}
