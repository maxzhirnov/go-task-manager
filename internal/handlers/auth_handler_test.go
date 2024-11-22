package handlers

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/maxzhirnov/go-task-manager/internal/middleware"
	"github.com/stretchr/testify/assert"
	"golang.org/x/crypto/bcrypt"
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

// TestLoginHandler_Success
func TestLoginHandler(t *testing.T) {
	// Generate a real bcrypt hash for "password123"
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte("password123"), bcrypt.DefaultCost)
	assert.NoError(t, err)

	tests := []struct {
		name         string
		payload      string
		setupMock    func(mock sqlmock.Sqlmock)
		expectedCode int
		expectedBody map[string]string
	}{
		{
			name:    "Successful login",
			payload: `{"username": "testuser", "password": "password123"}`,
			setupMock: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"id", "username", "password", "created_at", "updated_at"}).
					AddRow(1, "testuser", string(hashedPassword), time.Now(), time.Now())
				mock.ExpectQuery("SELECT (.+) FROM users WHERE username = \\$1").
					WithArgs("testuser").
					WillReturnRows(rows)
			},
			expectedCode: http.StatusOK,
			expectedBody: map[string]string{
				"access_token":  "mock-access-token",
				"refresh_token": "mock-refresh-token",
			},
		},
		{
			name:         "Empty credentials",
			payload:      `{"username": "", "password": ""}`,
			setupMock:    func(mock sqlmock.Sqlmock) {},
			expectedCode: http.StatusBadRequest,
			expectedBody: map[string]string{
				"error": "Username and password are required",
			},
		},
		{
			name:    "Invalid credentials",
			payload: `{"username": "nonexistent", "password": "wrongpass"}`,
			setupMock: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery("SELECT (.+) FROM users WHERE username = \\$1").
					WithArgs("nonexistent").
					WillReturnError(sql.ErrNoRows)
			},
			expectedCode: http.StatusUnauthorized,
			expectedBody: map[string]string{
				"error": "Invalid credentials",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Initialize sqlmock
			db, mock, err := sqlmock.New()
			assert.NoError(t, err)
			defer db.Close()

			// Setup mock expectations
			tt.setupMock(mock)

			// Initialize handler with mock functions
			authHandler := &AuthHandler{
				DB: db,
				GenerateJWT: func(userID int, username string) (string, error) {
					return "mock-access-token", nil
				},
				GenerateRefreshToken: func(username string) (string, error) {
					return "mock-refresh-token", nil
				},
			}

			// Create request
			req, err := http.NewRequest("POST", "/api/login", strings.NewReader(tt.payload))
			assert.NoError(t, err)
			req.Header.Set("Content-Type", "application/json")

			// Create response recorder
			rr := httptest.NewRecorder()

			// Call the handler
			authHandler.LoginHandler(rr, req)

			// Assert status code
			assert.Equal(t, tt.expectedCode, rr.Code)

			// Parse response body
			var response map[string]string
			err = json.NewDecoder(rr.Body).Decode(&response)
			assert.NoError(t, err)

			// Assert response body
			assert.Equal(t, tt.expectedBody, response)

			// Verify that all expected mock queries were executed
			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

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
