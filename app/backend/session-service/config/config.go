package config

import (
	"fmt"
	"os"
)

type Config struct {
	ServerPort string
	DBHost     string
	DBPort     string
	DBUser     string
	DBPass     string
	DBName     string
	DBSSLMode  string
}

func LoadConfig() *Config {
	return &Config{
		ServerPort: getEnv("SERVER_PORT", "8084"),
		DBHost:     getEnv("DB_HOST", "localhost"),
		DBPort:     getEnv("DB_PORT", "5432"),
		DBUser:     getEnv("DB_USER", "postgres"),
		DBPass:     getEnv("DB_PASS", "postgrespassword"),
		DBName:     getEnv("DB_NAME", "session_db"),
		DBSSLMode:  getEnv("DB_SSLMODE", "disable"),
	}
}

func (c *Config) DSN() string {
	return fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		c.DBHost, c.DBPort, c.DBUser, c.DBPass, c.DBName, c.DBSSLMode)
}

func getEnv(key, fallback string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	return fallback
}
