// Package database provides database connectivity and management for the application.
// It implements a retry mechanism for database connections and provides
// a common interface for database operations.
package database

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	_ "github.com/lib/pq"
	"github.com/maxzhirnov/go-task-manager/pkg/config"
)

// DB interface defines the required database operations.
// This interface allows for easy mocking in tests and provides
// a consistent API for database interactions.
type DB interface {
	// Query executes a query that returns rows
	Query(query string, args ...interface{}) (*sql.Rows, error)

	// QueryRow executes a query that returns at most one row
	QueryRow(query string, args ...interface{}) *sql.Row

	// Exec executes a query without returning any rows
	Exec(query string, args ...interface{}) (sql.Result, error)

	// Close closes the database connection
	Close() error

	// Begin starts a new transaction
	Begin() (*sql.Tx, error)
}

// DBConfig holds the configuration parameters for database connection.
type DBConfig struct {
	Host     string // Database server hostname
	Port     int    // Database server port
	User     string // Database username
	Password string // Database password
	DBName   string // Database name
}

// InitDB initializes a database connection using configuration from the environment.
//
// It loads the configuration and attempts to establish a connection with retry mechanism.
// The connection pool is configured with optimal settings for production use.
//
// Returns:
//   - DB: Database interface if connection is successful
//   - error: Any error encountered during initialization
//
// Example Usage:
//
//	db, err := InitDB()
//	if err != nil {
//	    log.Fatalf("Failed to initialize database: %v", err)
//	}
//	defer db.Close()
func InitDB() (DB, error) {
	// Load configuration
	cfg, err := config.LoadConfig()
	if err != nil {
		return nil, fmt.Errorf("failed to load config: %v", err)
	}

	// Create database config from loaded configuration
	dbConfig := DBConfig{
		Host:     cfg.Database.Host,
		Port:     cfg.Database.Port,
		User:     cfg.Database.User,
		Password: cfg.Database.Password,
		DBName:   cfg.Database.DBName,
	}

	return ConnectWithRetry(dbConfig)
}

// ConnectWithRetry attempts to establish a database connection with retry mechanism.
//
// It will attempt to connect up to maxRetries times with a 5-second delay between attempts.
// Once connected, it configures the connection pool with optimal settings.
//
// Parameters:
//   - cfg: Database configuration parameters
//
// Returns:
//   - DB: Database interface if connection is successful
//   - error: Any error encountered during connection attempts
//
// Connection Pool Settings:
//   - MaxOpenConns: 25
//   - MaxIdleConns: 25
//   - ConnMaxLifetime: 5 minutes
func ConnectWithRetry(cfg DBConfig) (DB, error) {
	// Create connection string
	connStr := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		cfg.Host, cfg.Port, cfg.User, cfg.Password, cfg.DBName)

	// Try to connect to the database with retries
	var db *sql.DB
	var err error
	maxRetries := 5

	// Attempt connection with retry
	for i := 0; i < maxRetries; i++ {
		db, err = sql.Open("postgres", connStr)
		if err != nil {
			log.Printf("Failed to open database connection: %v", err)
			time.Sleep(time.Second * 5)
			continue
		}

		// Verify connection is alive
		err = db.Ping()
		if err == nil {
			break
		}
		log.Printf("Failed to ping database (attempt %d/%d): %v", i+1, maxRetries, err)
		time.Sleep(time.Second * 5)
	}

	if err != nil {
		return nil, fmt.Errorf("failed to connect to database after %d attempts: %v", maxRetries, err)
	}

	// Configure connection pool
	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(25)
	db.SetConnMaxLifetime(5 * time.Minute)

	return db, nil
}

// InitTestDB initializes a database connection for testing purposes.
// It uses the provided configuration instead of loading from environment.
//
// Parameters:
//   - cfg: Database configuration parameters
//
// Returns:
//   - DB: Database interface if connection is successful
//   - error: Any error encountered during initialization
func InitTestDB(cfg DBConfig) (DB, error) {
	return ConnectWithRetry(cfg)
}
