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
	// This test remains the same as password hashing logic hasn't changed
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
				assert.True(t, len(u.Password) > 0)
				assert.True(t, u.Password != tt.password)
				assert.True(t, len(u.Password) > 50)
			}
		})
	}
}

func TestCheckPassword(t *testing.T) {
	// This test remains the same as password checking logic hasn't changed
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
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	user := &User{
		Email:    "test@example.com",
		Password: "password123",
	}
	err = user.HashPassword()
	assert.NoError(t, err)

	tests := []struct {
		name    string
		user    *User
		mockSQL func()
		wantErr bool
		errMsg  string
	}{
		{
			name: "Successful creation",
			user: user,
			mockSQL: func() {
				mock.ExpectQuery("INSERT INTO users").
					WithArgs(user.Email, sqlmock.AnyArg(), user.Password, sqlmock.AnyArg(), sqlmock.AnyArg()).
					WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))
			},
			wantErr: false,
		},
		{
			name: "Duplicate email",
			user: user,
			mockSQL: func() {
				mock.ExpectQuery("INSERT INTO users").
					WithArgs(user.Email, sqlmock.AnyArg(), user.Password, sqlmock.AnyArg(), sqlmock.AnyArg()).
					WillReturnError(&pq.Error{Code: "23505", Message: "email"})
			},
			wantErr: true,
			errMsg:  "email already exists",
		},
		{
			name:    "Empty email",
			user:    &User{Password: "password123"},
			mockSQL: func() {},
			wantErr: true,
			errMsg:  "email and password are required",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockSQL()
			err := tt.user.CreateUser(db)
			if tt.wantErr {
				assert.Error(t, err)
				assert.Equal(t, tt.errMsg, err.Error())
			} else {
				assert.NoError(t, err)
			}
			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestGetUserByEmail(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	now := time.Now()
	tests := []struct {
		name    string
		email   string
		mockSQL func()
		wantErr bool
	}{
		{
			name:  "User exists",
			email: "test@example.com",
			mockSQL: func() {
				rows := sqlmock.NewRows([]string{"id", "email", "username", "password", "created_at", "updated_at"}).
					AddRow(1, "test@example.com", "testuser", "$2a$10$N9qo8uLOickgx2ZMRZoMyeIjZAgcfl7p92ldGxad68LJZdL17lhWy", now, now)
				mock.ExpectQuery("SELECT (.+) FROM users WHERE email = \\$1").
					WithArgs("test@example.com").
					WillReturnRows(rows)
			},
			wantErr: false,
		},
		{
			name:  "User not found",
			email: "nonexistent@example.com",
			mockSQL: func() {
				mock.ExpectQuery("SELECT (.+) FROM users WHERE email = \\$1").
					WithArgs("nonexistent@example.com").
					WillReturnError(sql.ErrNoRows)
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockSQL()
			user, err := GetUserByEmail(db, tt.email)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.email, user.Email)
				assert.NotEmpty(t, user.Password)
			}
			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestGetUserByUsername(t *testing.T) {
	// This test can be kept for backward compatibility
	// but should be marked as deprecated if username-based lookup
	// will be removed in the future
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
