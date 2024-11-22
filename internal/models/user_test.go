package models

import (
	"database/sql"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/lib/pq"
	"github.com/stretchr/testify/assert"
)

func TestHashPassword(t *testing.T) {
	tests := []struct {
		name     string
		password string
		wantErr  bool
	}{
		{
			name:     "Valid password",
			password: "password123",
			wantErr:  false,
		},
		{
			name:     "Empty password",
			password: "",
			wantErr:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u := &User{Password: tt.password}
			err := u.HashPassword()

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				// Verify that the password was actually hashed
				assert.True(t, len(u.Password) > 0)
				assert.True(t, u.Password != tt.password)
				assert.True(t, len(u.Password) > 50) // bcrypt hashes are typically longer than 50 chars
			}
		})
	}
}

func TestCheckPassword(t *testing.T) {
	// Create a user with a known password
	u := &User{Password: "password123"}
	err := u.HashPassword()
	assert.NoError(t, err)

	tests := []struct {
		name     string
		password string
		wantErr  bool
	}{
		{
			name:     "Correct password",
			password: "password123",
			wantErr:  false,
		},
		{
			name:     "Wrong password",
			password: "wrongpassword",
			wantErr:  true,
		},
		{
			name:     "Empty password",
			password: "",
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := u.CheckPassword(tt.password)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestCreateUser(t *testing.T) {
	// Initialize sqlmock
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	// Create a test user with a hashed password
	user := &User{
		Username: "testuser",
		Password: "password123",
	}
	err = user.HashPassword()
	assert.NoError(t, err)

	// Test cases
	tests := []struct {
		name    string
		user    *User
		mockSQL func()
		wantErr bool
	}{
		{
			name: "Successful creation",
			user: user,
			mockSQL: func() {
				mock.ExpectQuery("INSERT INTO users").
					WithArgs(user.Username, user.Password, sqlmock.AnyArg(), sqlmock.AnyArg()).
					WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))
			},
			wantErr: false,
		},
		{
			name: "Duplicate username",
			user: user,
			mockSQL: func() {
				mock.ExpectQuery("INSERT INTO users").
					WithArgs(user.Username, user.Password, sqlmock.AnyArg(), sqlmock.AnyArg()).
					WillReturnError(&pq.Error{Code: "23505"})
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockSQL()
			err := tt.user.CreateUser(db)
			if tt.wantErr {
				assert.Error(t, err)
				if tt.name == "Duplicate username" {
					assert.Equal(t, "username already exists", err.Error())
				}
			} else {
				assert.NoError(t, err)
			}
			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestGetUserByUsername(t *testing.T) {
	// Initialize sqlmock
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	now := time.Now()
	tests := []struct {
		name     string
		username string
		mockSQL  func()
		wantErr  bool
	}{
		{
			name:     "User exists",
			username: "testuser",
			mockSQL: func() {
				rows := sqlmock.NewRows([]string{"id", "username", "password", "created_at", "updated_at"}).
					AddRow(1, "testuser", "$2a$10$N9qo8uLOickgx2ZMRZoMyeIjZAgcfl7p92ldGxad68LJZdL17lhWy", now, now)
				mock.ExpectQuery("SELECT (.+) FROM users WHERE username = \\$1").
					WithArgs("testuser").
					WillReturnRows(rows)
			},
			wantErr: false,
		},
		{
			name:     "User not found",
			username: "nonexistent",
			mockSQL: func() {
				mock.ExpectQuery("SELECT (.+) FROM users WHERE username = \\$1").
					WithArgs("nonexistent").
					WillReturnError(sql.ErrNoRows)
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockSQL()
			user, err := GetUserByUsername(db, tt.username)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.username, user.Username)
				assert.NotEmpty(t, user.Password)
			}
			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}
