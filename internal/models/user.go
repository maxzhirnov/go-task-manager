package models

import (
	"fmt"
	"log"
	"time"

	"github.com/lib/pq"
	"github.com/maxzhirnov/go-task-manager/pkg/database"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID        int       `json:"id"`
	Username  string    `json:"username"`
	Password  string    `json:"-"` // Don't expose the password
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
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
	return bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password))
}

// CreateUser inserts a new user into the database
func (u *User) CreateUser(db database.DB) error {
	query := `
        INSERT INTO users (username, password, created_at, updated_at)
        VALUES ($1, $2, $3, $4)
        RETURNING id`

	u.CreatedAt = time.Now()
	u.UpdatedAt = time.Now()

	log.Printf("Storing hashed password for user %s: %s", u.Username, u.Password)

	err := db.QueryRow(query, u.Username, u.Password, u.CreatedAt, u.UpdatedAt).Scan(&u.ID)
	if err != nil {
		// Check if the error is due to a unique constraint violation
		if pqErr, ok := err.(*pq.Error); ok && pqErr.Code == "23505" { // Unique constraint violation
			return fmt.Errorf("username already exists")
		}
		return err
	}

	return nil
}

// GetUserByUsername retrieves a user by their username
func GetUserByUsername(db database.DB, username string) (User, error) {
	var user User
	query := `SELECT id, username, password, created_at, updated_at FROM users WHERE username = $1`
	err := db.QueryRow(query, username).Scan(&user.ID, &user.Username, &user.Password, &user.CreatedAt, &user.UpdatedAt)
	return user, err
}
