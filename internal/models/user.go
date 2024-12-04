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

type User struct {
	ID         int       `json:"id"`
	Username   string    `json:"username"`
	Email      string    `json:"email"`
	Password   string    `json:"password,omitempty"`
	IsVerified bool      `json:"is_verified"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

type VerificationToken struct {
	ID        int        `json:"id"`
	UserID    int        `json:"user_id"`
	Token     string     `json:"token"`
	ExpiresAt time.Time  `json:"expires_at"`
	CreatedAt time.Time  `json:"created_at"`
	UsedAt    *time.Time `json:"used_at,omitempty"`
}

type UserStatistics struct {
	UserID            int    `json:"user_id"`
	Username          string `json:"username"`
	TotalTasks        int    `json:"total_tasks"`
	CompletedTasks    int    `json:"completed_tasks"`
	PendingTasks      int    `json:"pending_tasks"`
	InProgressTasks   int    `json:"in_progress_tasks"`
	DeletedTasks      int    `json:"deleted_tasks"`
	TasksCreatedToday int    `json:"tasks_created_today"`
}

// GenerateVerificationToken generates a random verification token
func GenerateVerificationToken() (string, error) {
	bytes := make([]byte, 32)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}

// CreateVerificationToken creates a verification token for a user
// Not using for now
func CreateVerificationToken(db database.DB, userID int) (*VerificationToken, error) {
	token, err := GenerateVerificationToken()
	if err != nil {
		return nil, err
	}

	vt := &VerificationToken{
		UserID:    userID,
		Token:     token,
		ExpiresAt: time.Now().Add(24 * time.Hour),
		CreatedAt: time.Now(),
	}

	query := `
        INSERT INTO verification_tokens (user_id, token, expires_at, created_at)
        VALUES ($1, $2, $3, $4)
        RETURNING id`

	err = db.QueryRow(query, vt.UserID, vt.Token, vt.ExpiresAt, vt.CreatedAt).Scan(&vt.ID)
	if err != nil {
		return nil, err
	}

	return vt, nil
}

func GetVerificationTokenForUser(db database.DB, userID int) (string, error) {
	var token string
	query := `
        SELECT token 
        FROM verification_tokens 
        WHERE user_id = $1 
        AND used_at IS NULL 
        AND expires_at > NOW() 
        ORDER BY created_at DESC 
        LIMIT 1`

	err := db.QueryRow(query, userID).Scan(&token)
	if err != nil {
		return "", err
	}
	return token, nil
}

func VerifyEmail(db database.DB, token string) error {
	tx, err := db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

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
		return err
	}

	// Обновляем статус верификации пользователя
	_, err = tx.Exec(`UPDATE users SET is_verified = true WHERE id = $1`, userID)
	if err != nil {
		return err
	}

	return tx.Commit()
}

// HashPassword hashes the user's password
func (u *User) HashPassword() error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	u.Password = string(hashedPassword) // Replace plaintext password with hash
	return nil
}

// CheckPassword compares a plaintext password with the hashed password
func (u *User) CheckPassword(password string) error {
	log.Printf("Attempting to compare:")
	log.Printf("Stored hash: %s", u.Password)
	log.Printf("Input password: %s", password)

	if password == "" {
		return fmt.Errorf("password cannot be empty")
	}

	err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password))
	if err != nil {
		log.Printf("Password comparison failed: %v", err)
	}
	return err
}

// CreateUser inserts a new user into the database
func (u *User) CreateUser(db database.DB) error {
	if u.Email == "" || u.Password == "" {
		return fmt.Errorf("email and password are required")
	}

	if u.Username == "" {
		u.Username = strings.Split(u.Email, "@")[0]
	}

	if !strings.HasPrefix(u.Password, "$2a$") {
		return fmt.Errorf("password must be hashed before saving")
	}

	tx, err := db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// Создаем пользователя
	query := `
        INSERT INTO users (email, username, password, is_verified, created_at, updated_at)
        VALUES ($1, $2, $3, $4, $5, $6)
        RETURNING id`

	u.CreatedAt = time.Now()
	u.UpdatedAt = time.Now()
	u.IsVerified = false

	err = tx.QueryRow(query, u.Email, u.Username, u.Password, u.IsVerified,
		u.CreatedAt, u.UpdatedAt).Scan(&u.ID)
	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok && pqErr.Code == "23505" {
			if strings.Contains(pqErr.Message, "email") {
				return fmt.Errorf("email already exists")
			}
			return fmt.Errorf("username already exists")
		}
		return err
	}

	// Создаем токен верификации
	token, err := GenerateVerificationToken()
	if err != nil {
		return err
	}

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
		return err
	}

	return tx.Commit()
}

func ResendVerificationToken(db database.DB, userID int) (*VerificationToken, error) {
	// Проверяем, не верифицирован ли уже пользователь
	var isVerified bool
	err := db.QueryRow(`SELECT is_verified FROM users WHERE id = $1`, userID).Scan(&isVerified)
	if err != nil {
		return nil, err
	}
	if isVerified {
		return nil, fmt.Errorf("user is already verified")
	}

	// Деактивируем старые токены
	_, err = db.Exec(`
        UPDATE verification_tokens 
        SET used_at = NOW() 
        WHERE user_id = $1 AND used_at IS NULL`,
		userID)
	if err != nil {
		return nil, err
	}

	// Создаем новый токен
	return CreateVerificationToken(db, userID)
}

// Add new method for finding user by email
func GetUserByEmail(db database.DB, email string) (User, error) {
	var user User
	query := `SELECT id, email, username, password, is_verified, created_at, updated_at 
              FROM users 
              WHERE email = $1`
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

// GetUserByUsername retrieves a user by their username
func GetUserByUsername(db database.DB, username string) (User, error) {
	var user User
	query := `SELECT id, username, password, created_at, updated_at FROM users WHERE username = $1`
	err := db.QueryRow(query, username).Scan(&user.ID, &user.Username, &user.Password, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		return user, err
	}
	if user.Password == "" {
		return user, fmt.Errorf("invalid password hash")
	}
	return user, nil
}

func GetUserStatistics(db database.DB, userID int) (*UserStatistics, error) {
	stats := &UserStatistics{}
	query := `
        SELECT 
            user_id, username, total_tasks, completed_tasks,
            pending_tasks, in_progress_tasks, deleted_tasks,
            tasks_created_today
        FROM user_statistics
        WHERE user_id = $1`

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
