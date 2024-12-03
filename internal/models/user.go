package models

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/lib/pq"
	"github.com/maxzhirnov/go-task-manager/pkg/database"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID        int       `json:"id"`
	Username  string    `json:"username"`
	Email     string    `json:"email"`
	Password  string    `json:"password,omitempty"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
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

	// Set username from email if not provided
	if u.Username == "" {
		u.Username = strings.Split(u.Email, "@")[0]
	}

	// Check if password is hashed
	if !strings.HasPrefix(u.Password, "$2a$") {
		return fmt.Errorf("password must be hashed before saving")
	}

	query := `
        INSERT INTO users (email, username, password, created_at, updated_at)
        VALUES ($1, $2, $3, $4, $5)
        RETURNING id`

	u.CreatedAt = time.Now()
	u.UpdatedAt = time.Now()

	err := db.QueryRow(query, u.Email, u.Username, u.Password, u.CreatedAt, u.UpdatedAt).Scan(&u.ID)
	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok && pqErr.Code == "23505" {
			if strings.Contains(pqErr.Message, "email") {
				return fmt.Errorf("email already exists")
			}
			return fmt.Errorf("username already exists")
		}
		return err
	}

	return nil
}

// Add new method for finding user by email
func GetUserByEmail(db database.DB, email string) (User, error) {
	var user User
	query := `SELECT id, email, username, password, created_at, updated_at FROM users WHERE email = $1`
	err := db.QueryRow(query, email).Scan(&user.ID, &user.Email, &user.Username, &user.Password, &user.CreatedAt, &user.UpdatedAt)
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
