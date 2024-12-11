// Package config provides configuration management for the task manager application.
// It handles loading and parsing of configuration values from environment variables
// with fallback to default values.
package config

import (
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

// Package config provides configuration management for the task manager application.
// It handles loading and parsing of configuration values from environment variables
// with fallback to default values.
type Config struct {
	// Database contains PostgreSQL connection settings
	Database struct {
		Host     string // Database server hostname
		Port     int    // Database server port
		User     string // Database user name
		Password string // Database user password
		DBName   string // Database name
	}

	// SMTP contains email server configuration and templates
	SMTP struct {
		Host     string // SMTP server hostname
		Port     int    // SMTP server port
		Username string // SMTP authentication username
		From     string // Sender email address
		Password string // SMTP authentication password
		FromName string // Sender name for emails
		BaseURL  string // Base URL for email links

		// Templates contains email template configuration
		Templates struct {
			Path string // Path to email template files
		}
	}

	Mixpanel struct {
		Token string // Mixpanel API token
	}

	// Server contains HTTP server settings
	Server struct {
		Port string // HTTP server port
	}

	// JWT contains JSON Web Token settings
	JWT struct {
		Secret string // Secret key for signing JWTs
	}
}

// LoadConfig reads configuration from environment variables and returns a Config instance.
// It first attempts to load variables from a .env file if present, then reads from
// the environment, falling back to default values when necessary.
//
// Environment Variables:
//
//	Database:
//	  - DB_HOST: Database server hostname (default: "localhost")
//	  - DB_PORT: Database server port (default: 5432)
//	  - DB_USER: Database username (default: "postgres")
//	  - DB_PASSWORD: Database password (default: "")
//	  - DB_NAME: Database name (default: "taskmanager")
//
//	SMTP:
//	  - SMTP_HOST: SMTP server hostname (default: "smtp.gmail.com")
//	  - SMTP_PORT: SMTP server port (default: 587)
//	  - SMTP_USERNAME: SMTP username
//	  - SMTP_PASSWORD: SMTP password
//	  - SMTP_FROM_NAME: Email sender name (default: "Task Manager")
//	  - SMTP_BASE_URL: Base URL for links (default: "http://localhost:8080")
//	  - SMTP_TEMPLATES_PATH: Email templates path (default: "templates/email")
//
//	Server:
//	  - SERVER_PORT: HTTP server port (default: "8080")
//
//	JWT:
//	  - JWT_SECRET: JWT signing key (default: "your-secret-key")
//
// Returns:
//   - *Config: Populated configuration struct
//   - error: Any error encountered during loading
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
	config.SMTP.From = getEnv("SMTP_FROM", "")
	config.SMTP.Password = getEnv("SMTP_PASSWORD", "")
	config.SMTP.FromName = getEnv("SMTP_FROM_NAME", "Task Manager")
	config.SMTP.BaseURL = getEnv("SMTP_BASE_URL", "http://localhost:8080")
	config.SMTP.Templates.Path = getEnv("SMTP_TEMPLATES_PATH", "templates/email")

	// Mixpanel configuration
	config.Mixpanel.Token = getEnv("MIXPANEL_TOKEN", "")

	// Server configuration
	config.Server.Port = getEnv("SERVER_PORT", "8080")

	// JWT configuration
	config.JWT.Secret = getEnv("JWT_SECRET", "your-secret-key")

	return config, nil
}

// getEnv retrieves an environment variable value or returns a default value if not set.
//
// Parameters:
//   - key: The environment variable name
//   - defaultValue: Value to return if environment variable is not set
//
// Returns:
//   - string: The environment variable value or default
func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}

// getEnvAsInt retrieves an environment variable and converts it to an integer.
// Returns the default value if the variable is not set or cannot be converted.
//
// Parameters:
//   - key: The environment variable name
//   - defaultValue: Value to return if environment variable is not set or invalid
//
// Returns:
//   - int: The parsed integer value or default
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
