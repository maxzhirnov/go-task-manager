// Package models provides data structures and database operations
// for the task management application.
package models

import (
	"crypto/rand"
	"database/sql"
	"encoding/hex"
	"fmt"
	"log"
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

	ResetPasswordToken string
	ResetTokenExpires  time.Time
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

	TasksLastWeek        int     `json:"tasks_last_week"`
	TasksThisWeek        int     `json:"tasks_this_week"`
	WeeklyTrendUp        bool    `json:"weekly_trend_up"`
	WeeklyTrendValue     int     `json:"weekly_trend_value"`
	PendingTasksLastWeek int     `json:"pending_tasks_last_week"`
	PendingTrendUp       bool    `json:"pending_trend_up"`
	PendingTrendValue    int     `json:"pending_trend_value"`
	AverageDailyTasks    float64 `json:"average_daily_tasks"`
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

	query := `
        WITH weekly_stats AS (
            SELECT 
                COUNT(*) FILTER (WHERE created_at >= NOW() - INTERVAL '7 days') as this_week,
                COUNT(*) FILTER (WHERE created_at >= NOW() - INTERVAL '14 days' 
                    AND created_at < NOW() - INTERVAL '7 days') as last_week,
                COUNT(*) FILTER (
                    WHERE status = 'pending' 
                    AND created_at >= NOW() - INTERVAL '7 days'
                ) as pending_this_week,
                COUNT(*) FILTER (
                    WHERE status = 'pending' 
                    AND created_at >= NOW() - INTERVAL '14 days'
                    AND created_at < NOW() - INTERVAL '7 days'
                ) as pending_last_week
            FROM tasks 
            WHERE user_id = $1
        ),
        daily_average AS (
            SELECT 
                ROUND(COUNT(*)::decimal / 
                    GREATEST(
                        EXTRACT(DAY FROM (NOW() - MIN(created_at)))::decimal,
                        1
                    ), 2) as avg_daily_tasks
            FROM tasks 
            WHERE user_id = $1
            AND created_at >= NOW() - INTERVAL '30 days'
        )
        SELECT 
            us.user_id, 
            us.username, 
            us.total_tasks, 
            us.completed_tasks,
            us.pending_tasks, 
            us.in_progress_tasks, 
            us.deleted_tasks,
            us.tasks_created_today,
            ws.last_week as tasks_last_week,
            ws.this_week as tasks_this_week,
            CASE 
                WHEN ws.last_week = 0 THEN true
                ELSE ws.this_week > ws.last_week 
            END as weekly_trend_up,
            CASE 
                WHEN ws.last_week = 0 THEN 
                    CASE 
                        WHEN ws.this_week = 0 THEN 0
                        ELSE 100
                    END
                ELSE ((ws.this_week - ws.last_week)::float / ws.last_week * 100)::int
            END as weekly_trend_value,
            ws.pending_last_week,
            CASE 
                WHEN ws.pending_last_week = 0 THEN true
                ELSE ws.pending_this_week > ws.pending_last_week 
            END as pending_trend_up,
            CASE 
                WHEN ws.pending_last_week = 0 THEN 
                    CASE 
                        WHEN ws.pending_this_week = 0 THEN 0
                        ELSE 100
                    END
                ELSE ((ws.pending_this_week - ws.pending_last_week)::float / ws.pending_last_week * 100)::int
            END as pending_trend_value,
            COALESCE(da.avg_daily_tasks, 0) as average_daily_tasks
        FROM user_statistics us
        CROSS JOIN weekly_stats ws
        CROSS JOIN daily_average da
        WHERE us.user_id = $1`

	err := db.QueryRow(query, userID).Scan(
		&stats.UserID,
		&stats.Username,
		&stats.TotalTasks,
		&stats.CompletedTasks,
		&stats.PendingTasks,
		&stats.InProgressTasks,
		&stats.DeletedTasks,
		&stats.TasksCreatedToday,
		&stats.TasksLastWeek,
		&stats.TasksThisWeek,
		&stats.WeeklyTrendUp,
		&stats.WeeklyTrendValue,
		&stats.PendingTasksLastWeek,
		&stats.PendingTrendUp,
		&stats.PendingTrendValue,
		&stats.AverageDailyTasks,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to get user statistics: %w", err)
	}

	return stats, nil
}

// GetUserByID retrieves a user from the database by their unique identifier.
//
// Parameters:
//   - db: database.DB interface for database operations
//   - id: int, the unique identifier of the user
//
// Returns:
//   - User: the user object if found
//   - error: nil if successful, otherwise an error describing the issue
//
// Errors:
//   - If the user is not found, it returns a sql.ErrNoRows wrapped in a custom error
//   - Any other database errors are wrapped and returned
func GetUserByID(db database.DB, id int) (User, error) {
	// Start timing the operation
	start := time.Now()

	// Log the start of the operation
	log.Printf("GetUserByID: Starting user retrieval for ID: %d", id)

	var user User

	// SQL query to fetch user details
	query := `
        SELECT id, email, username, password, is_verified, created_at, updated_at 
        FROM users 
        WHERE id = $1`

	// Execute the query and scan results into User struct
	err := db.QueryRow(query, id).Scan(
		&user.ID,
		&user.Email,
		&user.Username,
		&user.Password,
		&user.IsVerified,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	// Handle potential errors
	if err != nil {
		if err == sql.ErrNoRows {
			log.Printf("GetUserByID: User not found for ID: %d", id)
			return User{}, fmt.Errorf("user not found: %w", err)
		}
		log.Printf("GetUserByID: Database error while fetching user ID %d: %v", id, err)
		return User{}, fmt.Errorf("failed to get user by ID: %w", err)
	}

	// Log successful retrieval
	log.Printf("GetUserByID: Successfully retrieved user ID: %d, Username: %s", user.ID, user.Username)

	// Log operation duration
	log.Printf("GetUserByID: Operation completed in %v", time.Since(start))

	return user, nil
}

// UpdateProfile updates a user's profile information in the database.
// It supports updating username and/or password, requiring current password verification.
//
// Parameters:
//   - db: database interface for transactions
//   - userID: the ID of the user to update
//   - username: new username (empty string if no change)
//   - newPassword: new password (empty string if no change)
//   - currentPassword: current password for verification
//
// Security measures:
//   - Verifies current password before any changes
//   - Uses database transactions for atomic updates
//   - Hashes passwords using bcrypt
//   - Validates all inputs before processing
//
// Returns:
//   - error: nil if successful, otherwise contains the reason for failure
func (u *User) UpdateProfile(db database.DB, userID int, username, newPassword, currentPassword string) error {
	operation := "UpdateProfile"
	log.Printf("[%s] Starting profile update for user ID: %d", operation, userID)

	// Input validation
	if currentPassword == "" {
		log.Printf("[%s] Current password not provided for user ID: %d", operation, userID)
		return fmt.Errorf("current password is required")
	}

	// Verify current password
	if err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(currentPassword)); err != nil {
		log.Printf("[%s] Password verification failed for user ID: %d", operation, userID)
		return fmt.Errorf("invalid current password")
	}
	log.Printf("[%s] Password verification successful for user ID: %d", operation, userID)

	// Start transaction for atomic updates
	tx, err := db.Begin()
	if err != nil {
		log.Printf("[%s] Failed to start transaction for user ID %d: %v", operation, userID, err)
		return fmt.Errorf("database transaction error: %w", err)
	}
	defer func() {
		if tx != nil {
			if err := tx.Rollback(); err != nil && err != sql.ErrTxDone {
				log.Printf("[%s] Failed to rollback transaction for user ID %d: %v", operation, userID, err)
			}
		}
	}()

	// Build update query dynamically
	updates := []string{}
	args := []interface{}{}
	argCount := 1

	// Handle username update
	if username != "" {
		if username != u.Username {
			if err := validateUsername(username); err != nil {
				log.Printf("[%s] Invalid username format for user ID %d: %v", operation, userID, err)
				return fmt.Errorf("invalid username: %w", err)
			}
			log.Printf("[%s] Adding username update for user ID %d: %s", operation, userID, username)
			updates = append(updates, fmt.Sprintf("username = $%d", argCount))
			args = append(args, username)
			argCount++
		}
	}

	// Handle password update
	if newPassword != "" {
		if err := validatePassword(newPassword); err != nil {
			log.Printf("[%s] Invalid new password format for user ID %d: %v", operation, userID, err)
			return fmt.Errorf("invalid new password: %w", err)
		}

		log.Printf("[%s] Processing password update for user ID %d", operation, userID)
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
		if err != nil {
			log.Printf("[%s] Failed to hash new password for user ID %d: %v", operation, userID, err)
			return fmt.Errorf("password hashing error: %w", err)
		}
		updates = append(updates, fmt.Sprintf("password = $%d", argCount))
		args = append(args, string(hashedPassword))
		argCount++
	}

	// Add timestamp update
	updates = append(updates, fmt.Sprintf("updated_at = $%d", argCount))
	args = append(args, time.Now())
	argCount++

	// Add userID for WHERE clause
	args = append(args, userID)

	// Check if there are any updates to perform
	if len(updates) == 0 {
		log.Printf("[%s] No updates requested for user ID: %d", operation, userID)
		return nil
	}

	// Construct and execute update query
	query := fmt.Sprintf(`
        UPDATE users 
        SET %s 
        WHERE id = $%d`,
		strings.Join(updates, ", "),
		argCount,
	)

	log.Printf("[%s] Executing update query for user ID: %d", operation, userID)
	result, err := tx.Exec(query, args...)
	if err != nil {
		log.Printf("[%s] Failed to execute update query for user ID %d: %v", operation, userID, err)
		return fmt.Errorf("database update error: %w", err)
	}

	// Verify update success
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		log.Printf("[%s] Failed to get rows affected for user ID %d: %v", operation, userID, err)
		return fmt.Errorf("database result error: %w", err)
	}
	if rowsAffected == 0 {
		log.Printf("[%s] No rows updated for user ID: %d", operation, userID)
		return fmt.Errorf("user not found")
	}

	// Commit transaction
	if err := tx.Commit(); err != nil {
		log.Printf("[%s] Failed to commit transaction for user ID %d: %v", operation, userID, err)
		return fmt.Errorf("transaction commit error: %w", err)
	}
	tx = nil // Prevent rollback in deferred function

	log.Printf("[%s] Successfully updated profile for user ID: %d", operation, userID)
	return nil
}

// UpdateResetToken updates a user's password reset token and its expiration time.
// This is used during the password reset process to store the temporary reset token.
//
// Parameters:
//   - db: database interface for executing the update
//   - token: the generated reset token string
//   - expiry: timestamp when the token should expire
//
// Returns:
//   - error: nil if successful, otherwise contains the reason for failure
//
// Security considerations:
//   - Token expiration is enforced via database timestamp
//   - Updates timestamp to track modification time
//   - Uses parameterized queries to prevent SQL injection
func (u *User) UpdateResetToken(db database.DB, token string, expiry time.Time) error {
	operation := "UpdateResetToken"
	log.Printf("[%s] Starting reset token update for user ID: %d", operation, u.ID)

	// Validate inputs
	if token == "" {
		log.Printf("[%s] Empty token provided for user ID: %d", operation, u.ID)
		return fmt.Errorf("reset token cannot be empty")
	}

	if expiry.Before(time.Now()) {
		log.Printf("[%s] Invalid expiry time provided for user ID: %d", operation, u.ID)
		return fmt.Errorf("expiry time must be in the future")
	}

	// SQL query to update reset token information
	query := `
        UPDATE users 
        SET reset_password_token = $1, 
            reset_token_expires = $2,
            updated_at = NOW()
        WHERE id = $3`

	log.Printf("[%s] Executing update query for user ID: %d with expiry: %v",
		operation, u.ID, expiry.Format(time.RFC3339))

	// Execute the update query
	result, err := db.Exec(query, token, expiry, u.ID)
	if err != nil {
		log.Printf("[%s] Failed to update reset token for user ID %d: %v",
			operation, u.ID, err)
		return fmt.Errorf("failed to update reset token: %w", err)
	}

	// Verify the update was successful
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		log.Printf("[%s] Failed to get rows affected for user ID %d: %v",
			operation, u.ID, err)
		return fmt.Errorf("failed to verify update: %w", err)
	}

	if rowsAffected == 0 {
		log.Printf("[%s] No user found with ID: %d", operation, u.ID)
		return fmt.Errorf("user not found")
	}

	log.Printf("[%s] Successfully updated reset token for user ID: %d",
		operation, u.ID)

	// Log token expiration time for monitoring
	log.Printf("[%s] Reset token for user ID %d will expire at: %v",
		operation, u.ID, expiry.Format(time.RFC3339))

	return nil
}

// GetUserByResetToken retrieves a user from the database using a password reset token.
// This function is used during the password reset process to validate the reset token
// and retrieve the associated user.
//
// Parameters:
//   - db: database interface for executing the query
//   - token: the reset token string to look up
//
// Returns:
//   - User: the user associated with the token if found
//   - error: nil if successful, otherwise contains the reason for failure
//
// Security considerations:
//   - Does not reveal token validity through timing
//   - Logs are sanitized of sensitive information
//   - Uses parameterized queries to prevent SQL injection
func GetUserByResetToken(db database.DB, token string) (User, error) {
	operation := "GetUserByResetToken"
	start := time.Now()

	log.Printf("[%s] Starting user lookup by reset token", operation)

	// Validate input
	if token == "" {
		log.Printf("[%s] Empty token provided", operation)
		return User{}, fmt.Errorf("reset token cannot be empty")
	}

	var user User

	// SQL query to fetch user by reset token
	query := `
        SELECT id, email, username, password, is_verified, 
               created_at, updated_at, reset_token_expires
        FROM users 
        WHERE reset_password_token = $1`

	// Execute query and scan results
	err := db.QueryRow(query, token).Scan(
		&user.ID,
		&user.Email,
		&user.Username,
		&user.Password,
		&user.IsVerified,
		&user.CreatedAt,
		&user.UpdatedAt,
		&user.ResetTokenExpires,
	)

	// Handle potential errors
	if err != nil {
		if err == sql.ErrNoRows {
			log.Printf("[%s] No user found with provided reset token", operation)
			return User{}, fmt.Errorf("invalid or expired reset token")
		}
		log.Printf("[%s] Database error while fetching user: %v", operation, err)
		return User{}, fmt.Errorf("failed to get user by reset token: %w", err)
	}

	// Check if token has expired
	if time.Now().After(user.ResetTokenExpires) {
		log.Printf("[%s] Reset token expired for user ID: %d. Expired at: %v",
			operation, user.ID, user.ResetTokenExpires.Format(time.RFC3339))
		return User{}, fmt.Errorf("reset token has expired")
	}

	// Log successful retrieval (without sensitive information)
	log.Printf("[%s] Successfully retrieved user ID: %d with reset token",
		operation, user.ID)
	log.Printf("[%s] Operation completed in %v", operation, time.Since(start))

	return user, nil
}

// UpdatePasswordAndClearResetToken updates the user's password and clears the reset token.
// This function is typically called after a successful password reset.
//
// Parameters:
//   - db: database interface for executing the update
//   - hashedPassword: the new password hash to set
//
// Returns:
//   - error: nil if successful, otherwise contains the reason for failure
//
// Security considerations:
//   - Clears reset token to prevent reuse
//   - Updates timestamp for audit trail
//   - Uses parameterized query to prevent SQL injection
func (u *User) UpdatePasswordAndClearResetToken(db database.DB, hashedPassword string) error {
	operation := "UpdatePasswordAndClearResetToken"
	log.Printf("[%s] Starting password update and token clear for user ID: %d", operation, u.ID)

	// Validate input
	if hashedPassword == "" {
		log.Printf("[%s] Empty hashed password provided for user ID: %d", operation, u.ID)
		return fmt.Errorf("hashed password cannot be empty")
	}

	// SQL query to update password and clear reset token
	query := `
        UPDATE users 
        SET password = $1,
            reset_password_token = NULL,
            reset_token_expires = NULL,
            updated_at = NOW()
        WHERE id = $2`

	log.Printf("[%s] Executing update query for user ID: %d", operation, u.ID)

	// Execute the update query
	result, err := db.Exec(query, hashedPassword, u.ID)
	if err != nil {
		log.Printf("[%s] Failed to update password and clear reset token for user ID %d: %v",
			operation, u.ID, err)
		return fmt.Errorf("failed to update password: %w", err)
	}

	// Verify the update was successful
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		log.Printf("[%s] Failed to get rows affected for user ID %d: %v",
			operation, u.ID, err)
		return fmt.Errorf("failed to verify update: %w", err)
	}

	if rowsAffected == 0 {
		log.Printf("[%s] No user found with ID: %d", operation, u.ID)
		return fmt.Errorf("user not found")
	}

	log.Printf("[%s] Successfully updated password and cleared reset token for user ID: %d",
		operation, u.ID)

	return nil
}
