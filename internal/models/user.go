// Package models provides data structures and database operations
// for the task management application.
package models

import (
	"crypto/rand"
	"database/sql"
	"encoding/hex"
	"fmt"
	"strings"
	"time"

	"github.com/lib/pq"
	"github.com/maxzhirnov/go-task-manager/pkg/database"
	"golang.org/x/crypto/bcrypt"
)

// User represents a registered user in the system.
// It contains all user-related information including authentication
// and verification status. The password field is omitted from JSON
// responses for security.
type User struct {
	// ID uniquely identifies the user
	ID int `json:"id"`

	// Username is derived from the user's email
	Username string `json:"username"`

	// Email is the user's unique identifier for authentication
	Email string `json:"email"`

	// Password is the hashed password, omitted from JSON responses
	Password string `json:"password,omitempty"`

	// IsVerified indicates whether the email has been verified
	IsVerified bool `json:"is_verified"`

	// CreatedAt stores the timestamp of user registration
	CreatedAt time.Time `json:"created_at"`

	// UpdatedAt stores the timestamp of the last modification
	UpdatedAt time.Time `json:"updated_at"`
}

// VerificationToken represents an email verification token.
// It is used to verify user email addresses and includes expiration
// and usage tracking.
type VerificationToken struct {
	// ID uniquely identifies the verification token
	ID int `json:"id"`

	// UserID associates the token with a specific user
	UserID int `json:"user_id"`

	// Token is the actual verification token string
	Token string `json:"token"`

	// ExpiresAt indicates when the token becomes invalid
	ExpiresAt time.Time `json:"expires_at"`

	// CreatedAt stores when the token was generated
	CreatedAt time.Time `json:"created_at"`

	// UsedAt stores when the token was used (if it has been)
	UsedAt *time.Time `json:"used_at,omitempty"`
}

// UserStatistics represents aggregated statistics about a user's tasks.
// It provides various metrics about task status and activity.
type UserStatistics struct {
	// UserID identifies the user these statistics belong to
	UserID int `json:"user_id"`

	// Username of the statistics owner
	Username string `json:"username"`

	// TotalTasks is the total number of tasks created by the user
	TotalTasks int `json:"total_tasks"`

	// CompletedTasks is the number of tasks marked as completed
	CompletedTasks int `json:"completed_tasks"`

	// PendingTasks is the number of tasks not yet started
	PendingTasks int `json:"pending_tasks"`

	// InProgressTasks is the number of tasks currently being worked on
	InProgressTasks int `json:"in_progress_tasks"`

	// DeletedTasks is the number of soft-deleted tasks
	DeletedTasks int `json:"deleted_tasks"`

	// TasksCreatedToday is the number of tasks created in the last 24 hours
	TasksCreatedToday int `json:"tasks_created_today"`
}

// GenerateVerificationToken creates a secure random token for email verification.
//
// It generates a cryptographically secure random token using crypto/rand
// and encodes it as a hexadecimal string. The resulting token is 64 characters
// long (32 bytes of random data encoded as hex).
//
// Returns:
//   - string: A 64-character hexadecimal token
//   - error: If random number generation fails
//
// Security Features:
//   - Uses crypto/rand for cryptographic security
//   - Produces 256 bits of entropy
//   - Returns hex-encoded string for safe transmission
//
// Example Usage:
//
//	token, err := GenerateVerificationToken()
//	if err != nil {
//	    return fmt.Errorf("failed to generate verification token: %w", err)
//	}
//	// token example: "8f7d3b2e1a..."
//
// Note: This function is suitable for email verification as it provides
// sufficient entropy to prevent token guessing attacks.
func GenerateVerificationToken() (string, error) {
	// Create byte slice for random data
	bytes := make([]byte, 32) // 32 bytes = 256 bits of entropy

	// Generate random bytes using crypto/rand
	if _, err := rand.Read(bytes); err != nil {
		return "", fmt.Errorf("failed to generate random bytes: %w", err)
	}

	// Convert bytes to hexadecimal string
	return hex.EncodeToString(bytes), nil
}

// CreateVerificationToken creates a verification token for a user and stores it in the database.
//
// This function generates a new verification token, associates it with the given user ID,
// and inserts it into the database with an expiration time of 24 hours from creation.
//
// Parameters:
//   - db: Database interface for executing queries
//   - userID: The ID of the user for whom the token is being created
//
// Returns:
//   - *VerificationToken: The created verification token
//   - error: Any error encountered during token creation or database insertion
//
// Token Lifecycle:
//   - Generated using GenerateVerificationToken()
//   - Valid for 24 hours from creation
//   - Stored in the database for validation
//
// Example Usage:
//
//	token, err := CreateVerificationToken(db, user.ID)
//	if err != nil {
//	    return fmt.Errorf("failed to create verification token: %w", err)
//	}
//	// Use token.Token for email verification link
//
// Note: This function is currently not in use but kept for potential future implementation.
func CreateVerificationToken(db database.DB, userID int) (*VerificationToken, error) {
	// Generate a new verification token
	token, err := GenerateVerificationToken()
	if err != nil {
		return nil, fmt.Errorf("failed to generate token: %w", err)
	}

	// Create VerificationToken struct
	vt := &VerificationToken{
		UserID:    userID,
		Token:     token,
		ExpiresAt: time.Now().Add(24 * time.Hour),
		CreatedAt: time.Now(),
	}

	// SQL query to insert the token
	query := `
        INSERT INTO verification_tokens (user_id, token, expires_at, created_at)
        VALUES ($1, $2, $3, $4)
        RETURNING id`

	// Execute query and get the inserted ID
	err = db.QueryRow(query, vt.UserID, vt.Token, vt.ExpiresAt, vt.CreatedAt).Scan(&vt.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to insert verification token: %w", err)
	}

	return vt, nil
}

// GetVerificationTokenForUser retrieves the most recent valid verification token for a user.
//
// This function fetches the latest unused and non-expired verification token
// for the specified user. It's typically used when resending verification
// emails or checking token validity.
//
// Parameters:
//   - db: Database interface for executing queries
//   - userID: The ID of the user whose token to retrieve
//
// Returns:
//   - string: The verification token if found
//   - error: sql.ErrNoRows if no valid token exists, or other database errors
//
// Query Conditions:
//   - Token must not be used (used_at IS NULL)
//   - Token must not be expired (expires_at > NOW())
//   - Returns the most recently created token
//
// Example Usage:
//
//	token, err := GetVerificationTokenForUser(db, userID)
//	if err == sql.ErrNoRows {
//	    // No valid token exists, create new one
//	    return CreateVerificationToken(db, userID)
//	}
//	if err != nil {
//	    return "", fmt.Errorf("failed to get verification token: %w", err)
//	}
//
// Note: This function only returns the token string, not the full VerificationToken
// struct, as typically only the token string is needed for verification purposes.
func GetVerificationTokenForUser(db database.DB, userID int) (string, error) {
	var token string

	// SQL query to get the most recent valid token
	query := `
        SELECT token 
        FROM verification_tokens 
        WHERE user_id = $1 
        AND used_at IS NULL 
        AND expires_at > NOW() 
        ORDER BY created_at DESC 
        LIMIT 1`

	// Execute query and scan result
	err := db.QueryRow(query, userID).Scan(&token)
	if err != nil {
		if err == sql.ErrNoRows {
			return "", err // Return as is for specific handling
		}
		return "", fmt.Errorf("database error getting verification token: %w", err)
	}

	return token, nil
}

// VerifyEmail processes an email verification token and updates the user's verification status.
//
// This function uses a transaction to ensure atomicity when marking both the token as used
// and updating the user's verification status. It validates that the token exists,
// hasn't expired, and hasn't been used before.
//
// Parameters:
//   - db: Database interface for executing queries
//   - token: The verification token string to process
//
// Returns:
//   - error: "invalid or expired verification token" if token is invalid,
//     or other database errors if they occur
//
// Transaction Steps:
// 1. Mark the verification token as used
// 2. Update the user's verification status
//
// Example Usage:
//
//	err := VerifyEmail(db, "verification_token_string")
//	if err != nil {
//	    if strings.Contains(err.Error(), "invalid or expired") {
//	        return fmt.Errorf("verification failed: %w", err)
//	    }
//	    return fmt.Errorf("error during verification: %w", err)
//	}
//
// Note: This operation cannot be undone, and the token cannot be reused
// once marked as used.
func VerifyEmail(db database.DB, token string) error {
	// Start transaction
	tx, err := db.Begin()
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	// Mark token as used and get associated user ID
	var userID int
	query := `
        UPDATE verification_tokens 
        SET used_at = NOW()
        WHERE token = $1 
        AND expires_at > NOW() 
        AND used_at IS NULL
        RETURNING user_id`

	err = tx.QueryRow(query, token).Scan(&userID)
	if err != nil {
		if err == sql.ErrNoRows {
			return fmt.Errorf("invalid or expired verification token")
		}
		return fmt.Errorf("failed to process verification token: %w", err)
	}

	// Update user's verification status
	_, err = tx.Exec(`UPDATE users SET is_verified = true WHERE id = $1`, userID)
	if err != nil {
		return fmt.Errorf("failed to update user verification status: %w", err)
	}

	// Commit transaction
	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit verification transaction: %w", err)
	}

	return nil
}

// HashPassword securely hashes the user's password using bcrypt.
//
// This method replaces the plaintext password in the User struct with its
// hashed version. It uses bcrypt with the default cost factor for a balance
// of security and performance.
//
// Returns:
//   - error: If password hashing fails
//
// Security Features:
//   - Uses bcrypt hashing algorithm
//   - Includes automatic salt generation
//   - Uses default cost factor (10)
//   - Overwrites plaintext password
//
// Example Usage:
//
//	user := &User{Password: "plaintext_password"}
//	if err := user.HashPassword(); err != nil {
//	    return fmt.Errorf("failed to hash password: %w", err)
//	}
//	// Password is now hashed and original is destroyed
//
// Note: This method should be called before storing the user in the database
// to ensure passwords are never stored in plaintext.
func (u *User) HashPassword() error {
	// Generate bcrypt hash from password
	hashedPassword, err := bcrypt.GenerateFromPassword(
		[]byte(u.Password),
		bcrypt.DefaultCost,
	)
	if err != nil {
		return fmt.Errorf("failed to hash password: %w", err)
	}

	// Replace plaintext password with hash
	u.Password = string(hashedPassword)
	return nil
}

// CheckPassword verifies if the provided password matches the user's hashed password.
//
// This method uses bcrypt's secure comparison to verify the password without
// exposing the original hash. It ensures constant-time comparison to prevent
// timing attacks.
//
// Parameters:
//   - password: The plaintext password to verify
//
// Returns:
//   - error: nil if password matches, bcrypt.ErrMismatchedHashAndPassword if it doesn't,
//     or other errors if validation fails
//
// Security Features:
//   - Constant-time comparison
//   - No password logging
//   - Early validation of empty passwords
//
// Example Usage:
//
//	if err := user.CheckPassword(inputPassword); err != nil {
//	    if err == bcrypt.ErrMismatchedHashAndPassword {
//	        return fmt.Errorf("invalid credentials")
//	    }
//	    return fmt.Errorf("error validating password: %w", err)
//	}
//
// Note: This method should never log passwords or hashes in production
// to prevent security breaches.
func (u *User) CheckPassword(password string) error {
	// Validate password is not empty
	if password == "" {
		return fmt.Errorf("password cannot be empty")
	}

	// Compare password with stored hash
	return bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password))
}

// CreateUser creates a new user record in the database along with a verification token.
//
// This method performs several operations within a transaction:
// 1. Validates user input
// 2. Sets default username if not provided
// 3. Creates user record
// 4. Generates and stores verification token
//
// Parameters:
//   - db: Database interface for executing queries
//
// Returns:
//   - error: Validation error, database error, or "email already exists"
//
// Validation Rules:
//   - Email and password are required
//   - Password must be pre-hashed (starts with "$2a$")
//   - Email must be unique in the system
//
// Side Effects:
//   - Sets CreatedAt and UpdatedAt to current time
//   - Sets IsVerified to false
//   - Generates username from email if not provided
//   - Sets user ID from database
//
// Example Usage:
//
//	user := &User{
//	    Email: "user@example.com",
//	    Password: hashedPassword,
//	}
//	if err := user.CreateUser(db); err != nil {
//	    if err.Error() == "email already exists" {
//	        return handleDuplicateEmail(err)
//	    }
//	    return fmt.Errorf("failed to create user: %w", err)
//	}
func (u *User) CreateUser(db database.DB) error {
	// Validate required fields
	if u.Email == "" || u.Password == "" {
		return fmt.Errorf("email and password are required")
	}

	// Set default username if not provided
	if u.Username == "" {
		u.Username = strings.Split(u.Email, "@")[0]
	}

	// Ensure password is hashed
	if !strings.HasPrefix(u.Password, "$2a$") {
		return fmt.Errorf("password must be hashed before saving")
	}

	// Start transaction
	tx, err := db.Begin()
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	// Set timestamps and initial verification status
	u.CreatedAt = time.Now()
	u.UpdatedAt = u.CreatedAt
	u.IsVerified = false

	// Create user record
	query := `
        INSERT INTO users (email, username, password, is_verified, created_at, updated_at)
        VALUES ($1, $2, $3, $4, $5, $6)
        RETURNING id`

	err = tx.QueryRow(query, u.Email, u.Username, u.Password, u.IsVerified,
		u.CreatedAt, u.UpdatedAt).Scan(&u.ID)
	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok && pqErr.Code == "23505" {
			if strings.Contains(pqErr.Message, "email") {
				return fmt.Errorf("email already exists")
			}
			return fmt.Errorf("username already exists")
		}
		return fmt.Errorf("failed to create user: %w", err)
	}

	// Generate verification token
	token, err := GenerateVerificationToken()
	if err != nil {
		return fmt.Errorf("failed to generate verification token: %w", err)
	}

	// Store verification token
	verificationQuery := `
        INSERT INTO verification_tokens (user_id, token, expires_at, created_at)
        VALUES ($1, $2, $3, $4)`

	_, err = tx.Exec(verificationQuery,
		u.ID,
		token,
		time.Now().Add(24*time.Hour),
		time.Now(),
	)
	if err != nil {
		return fmt.Errorf("failed to create verification token: %w", err)
	}

	// Commit transaction
	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

// ResendVerificationToken creates a new verification token after invalidating existing ones.
//
// This function performs three steps:
// 1. Checks if the user is already verified
// 2. Invalidates all existing unused tokens for the user
// 3. Creates and returns a new verification token
//
// Parameters:
//   - db: Database interface for executing queries
//   - userID: The ID of the user requesting a new verification token
//
// Returns:
//   - *VerificationToken: Newly created verification token
//   - error: "user is already verified" or database errors
//
// Example Usage:
//
//	token, err := ResendVerificationToken(db, userID)
//	if err != nil {
//	    if err.Error() == "user is already verified" {
//	        return nil, fmt.Errorf("verification not needed: %w", err)
//	    }
//	    return nil, fmt.Errorf("failed to resend verification token: %w", err)
//	}
//
// Note: This function invalidates all existing tokens for security purposes,
// ensuring only one active token exists at a time.
func ResendVerificationToken(db database.DB, userID int) (*VerificationToken, error) {
	// Check if user is already verified
	var isVerified bool
	err := db.QueryRow(`SELECT is_verified FROM users WHERE id = $1`, userID).Scan(&isVerified)
	if err != nil {
		return nil, fmt.Errorf("failed to check user verification status: %w", err)
	}
	if isVerified {
		return nil, fmt.Errorf("user is already verified")
	}

	// Invalidate existing unused tokens
	_, err = db.Exec(`
        UPDATE verification_tokens 
        SET used_at = NOW() 
        WHERE user_id = $1 AND used_at IS NULL`,
		userID)
	if err != nil {
		return nil, fmt.Errorf("failed to invalidate existing tokens: %w", err)
	}

	// Create new verification token
	return CreateVerificationToken(db, userID)
}

// GetUserByEmail retrieves a user record from the database by email address.
//
// This function performs a direct lookup using the email address, which should
// be unique in the system. It returns the complete user record including
// the hashed password for authentication purposes.
//
// Parameters:
//   - db: Database interface for executing queries
//   - email: The email address to search for
//
// Returns:
//   - User: The found user's data
//   - error: sql.ErrNoRows if user not found, or other database errors
//
// Example Usage:
//
//	user, err := GetUserByEmail(db, "user@example.com")
//	if err == sql.ErrNoRows {
//	    return User{}, fmt.Errorf("user not found")
//	}
//	if err != nil {
//	    return User{}, fmt.Errorf("failed to fetch user: %w", err)
//	}
//
// Security Note:
//   - Returns hashed password for authentication
//   - Should only be used in authentication contexts
//   - Consider creating a separate method that excludes password for general use
func GetUserByEmail(db database.DB, email string) (User, error) {
	var user User

	// SQL query to fetch user by email
	query := `SELECT id, email, username, password, is_verified, created_at, updated_at 
              FROM users 
              WHERE email = $1`

	// Execute query and scan results into User struct
	err := db.QueryRow(query, email).Scan(
		&user.ID,
		&user.Email,
		&user.Username,
		&user.Password,
		&user.IsVerified,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	if err != nil {
		return user, err
	}
	return user, nil
}

// GetUserStatistics retrieves aggregated task statistics for a specific user.
//
// This function fetches pre-calculated statistics from the user_statistics view,
// providing various metrics about the user's tasks including counts by status
// and recent activity.
//
// Parameters:
//   - db: Database interface for executing queries
//   - userID: The ID of the user whose statistics to retrieve
//
// Returns:
//   - *UserStatistics: Pointer to structure containing all statistics
//   - error: sql.ErrNoRows if user not found, or other database errors
//
// Statistics Included:
//   - Total tasks count
//   - Tasks by status (completed, pending, in progress, deleted)
//   - Tasks created in the last 24 hours
//
// Example Usage:
//
//	stats, err := GetUserStatistics(db, userID)
//	if err == sql.ErrNoRows {
//	    return nil, fmt.Errorf("no statistics found for user")
//	}
//	if err != nil {
//	    return nil, fmt.Errorf("failed to fetch user statistics: %w", err)
//	}
//
// Note: This function relies on a database view 'user_statistics'
// which should be kept up to date with task changes.
func GetUserStatistics(db database.DB, userID int) (*UserStatistics, error) {
	stats := &UserStatistics{}

	// SQL query to fetch statistics from view
	query := `
        SELECT 
            user_id, username, total_tasks, completed_tasks,
            pending_tasks, in_progress_tasks, deleted_tasks,
            tasks_created_today
        FROM user_statistics
        WHERE user_id = $1`

	// Execute query and scan results into UserStatistics struct
	err := db.QueryRow(query, userID).Scan(
		&stats.UserID, &stats.Username, &stats.TotalTasks,
		&stats.CompletedTasks, &stats.PendingTasks,
		&stats.InProgressTasks, &stats.DeletedTasks,
		&stats.TasksCreatedToday,
	)
	if err != nil {
		return nil, err
	}
	return stats, nil
}
