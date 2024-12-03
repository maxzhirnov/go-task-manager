package config

import (
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	Database struct {
		Host     string
		Port     int
		User     string
		Password string
		DBName   string
	}
	SMTP struct {
		Host     string
		Port     int
		Username string
		Password string
	}
	Server struct {
		Port string
	}
	JWT struct {
		Secret string
	}
}

func LoadConfig() (*Config, error) {
	// Load .env file if it exists
	godotenv.Load()

	config := &Config{}

	// Database configuration
	config.Database.Host = getEnv("DB_HOST", "localhost")
	config.Database.Port = getEnvAsInt("DB_PORT", 5432)
	config.Database.User = getEnv("DB_USER", "postgres")
	config.Database.Password = getEnv("DB_PASSWORD", "")
	config.Database.DBName = getEnv("DB_NAME", "taskmanager")

	// SMTP configuration
	config.SMTP.Host = getEnv("SMTP_HOST", "smtp.gmail.com")
	config.SMTP.Port = getEnvAsInt("SMTP_PORT", 587)
	config.SMTP.Username = getEnv("SMTP_USERNAME", "")
	config.SMTP.Password = getEnv("SMTP_PASSWORD", "")

	// Server configuration
	config.Server.Port = getEnv("SERVER_PORT", "8080")

	// JWT configuration
	config.JWT.Secret = getEnv("JWT_SECRET", "your-secret-key")

	return config, nil
}

// Helper function to get environment variable with a default value
func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}

// Helper function to get environment variable as integer
func getEnvAsInt(key string, defaultValue int) int {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	intValue, err := strconv.Atoi(value)
	if err != nil {
		return defaultValue
	}
	return intValue
}
