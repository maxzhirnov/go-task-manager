package handlers

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/maxzhirnov/go-task-manager/internal/middleware"
	"github.com/stretchr/testify/assert"
)

func TestRegisterHandler_Success(t *testing.T) {
	// Initialize sqlmock
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	// Mock the expected SQL query for user creation
	mock.ExpectQuery("INSERT INTO users").
		WithArgs("testuser", sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg()).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))

	// Initialize the handler
	authHandler := NewAuthHandler(db)

	// Create a test user payload
	userPayload := `{
		"username": "testuser",
		"password": "password123"
	}`

	// Create a test request
	req, err := http.NewRequest("POST", "/api/register", bytes.NewBufferString(userPayload))
	assert.NoError(t, err)

	// Create a response recorder
	rr := httptest.NewRecorder()

	// Call the handler
	authHandler.RegisterHandler(rr, req)

	// Assert the status code
	assert.Equal(t, http.StatusCreated, rr.Code)

	// Assert the response body
	var response map[string]string
	err = json.Unmarshal(rr.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "User registered successfully", response["message"])

	// Ensure all SQL expectations were met
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestRegisterHandler_DuplicateUser(t *testing.T) {
	// Initialize sqlmock
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	// Mock the expected SQL query to return a unique constraint violation
	mock.ExpectQuery("INSERT INTO users").
		WithArgs("testuser", sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg()).
		WillReturnError(errors.New("username already exists"))

	// Initialize the handler
	authHandler := NewAuthHandler(db)

	// Create a test user payload
	userPayload := `{
		"username": "testuser",
		"password": "password123"
	}`

	// Create a test request
	req, err := http.NewRequest("POST", "/api/register", bytes.NewBufferString(userPayload))
	assert.NoError(t, err)

	// Create a response recorder
	rr := httptest.NewRecorder()

	// Call the handler
	authHandler.RegisterHandler(rr, req)

	// Assert the status code
	assert.Equal(t, http.StatusConflict, rr.Code)

	// Assert the response body
	assert.Contains(t, rr.Body.String(), "Username already exists")

	// Ensure all SQL expectations were met
	assert.NoError(t, mock.ExpectationsWereMet())
}

// func TestLoginHandler_Success(t *testing.T) {
// 	// Initialize sqlmock
// 	db, mock, err := sqlmock.New()
// 	assert.NoError(t, err)
// 	defer db.Close()

// 	// Mock the expected SQL query to retrieve a user
// 	// Note: Using $1 instead of ? for PostgreSQL
// 	mock.ExpectQuery(`SELECT id, username, password, created_at, updated_at FROM users WHERE username = \$1`).
// 		WithArgs("testuser").
// 		WillReturnRows(sqlmock.NewRows([]string{"id", "username", "password", "created_at", "updated_at"}).
// 			AddRow(1, "testuser", "$2a$10$N9qo8uLOickgx2ZMRZoMyeIjZAgcfl7p92ldGxad68LJZdL17lhWy", time.Now(), time.Now()))

// 	// Initialize the handler with mock functions
// 	authHandler := &AuthHandler{
// 		DB: db,
// 		GenerateJWT: func(userID int, username string) (string, error) {
// 			return "mock-access-token", nil
// 		},
// 		GenerateRefreshToken: func(username string) (string, error) {
// 			return "mock-refresh-token", nil
// 		},
// 		ValidateRefreshToken: func(token string) (*middleware.Claims, error) {
// 			return &middleware.Claims{
// 				UserID:   1,
// 				Username: "testuser",
// 			}, nil
// 		},
// 	}

// 	// Create a test login payload
// 	loginPayload := `{
//         "username": "testuser",
//         "password": "password123"
//     }`

// 	// Create a test request
// 	req, err := http.NewRequest("POST", "/api/login", bytes.NewBufferString(loginPayload))
// 	assert.NoError(t, err)
// 	req.Header.Set("Content-Type", "application/json")

// 	// Create a response recorder
// 	rr := httptest.NewRecorder()

// 	// Call the handler
// 	authHandler.LoginHandler(rr, req)

// 	// Assert the status code
// 	assert.Equal(t, http.StatusOK, rr.Code)

// 	// Assert the response body
// 	var response map[string]string
// 	err = json.Unmarshal(rr.Body.Bytes(), &response)
// 	assert.NoError(t, err)
// 	assert.Equal(t, "mock-access-token", response["access_token"])
// 	assert.Equal(t, "mock-refresh-token", response["refresh_token"])

// 	// Ensure all SQL expectations were met
// 	assert.NoError(t, mock.ExpectationsWereMet())
// }

func TestRefreshTokenHandler_Success(t *testing.T) {
	// Initialize sqlmock
	db, _, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	// Mock refresh token validation
	mockValidateRefreshToken := func(token string) (*middleware.Claims, error) {
		return &middleware.Claims{
			UserID:   1,
			Username: "testuser",
		}, nil
	}

	// Mock JWT generation
	mockGenerateJWT := func(userID int, username string) (string, error) {
		return "mock-new-access-token", nil
	}

	// Initialize the handler with mock functions
	authHandler := &AuthHandler{
		DB:                   db,
		GenerateJWT:          mockGenerateJWT,
		ValidateRefreshToken: mockValidateRefreshToken,
	}

	// Create a test refresh token payload
	refreshPayload := `{
		"refresh_token": "mock-refresh-token"
	}`

	// Create a test request
	req, err := http.NewRequest("POST", "/api/refresh", bytes.NewBufferString(refreshPayload))
	assert.NoError(t, err)

	// Create a response recorder
	rr := httptest.NewRecorder()

	// Call the handler
	authHandler.RefreshTokenHandler(rr, req)

	// Assert the status code
	assert.Equal(t, http.StatusOK, rr.Code)

	// Assert the response body
	var token map[string]string
	err = json.Unmarshal(rr.Body.Bytes(), &token)
	assert.NoError(t, err)
	assert.Equal(t, "mock-new-access-token", token["access_token"])
}
