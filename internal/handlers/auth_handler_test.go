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
	"github.com/maxzhirnov/go-task-manager/pkg/email"
	"github.com/stretchr/testify/assert"
	"golang.org/x/crypto/bcrypt"
)

func TestRegisterHandler_Success(t *testing.T) {

	mockEmail := email.NewMockEmailService()

	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	// Mock the SQL query for user creation with email
	mock.ExpectQuery("INSERT INTO users").
		WithArgs("test@example.com", sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg()).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))

	authHandler := NewAuthHandler(db, mockEmail)

	userPayload := `{
        "email": "test@example.com",
        "password": "password123"
    }`

	req, err := http.NewRequest("POST", "/api/register", bytes.NewBufferString(userPayload))
	assert.NoError(t, err)

	rr := httptest.NewRecorder()
	authHandler.RegisterHandler(rr, req)

	assert.Equal(t, http.StatusCreated, rr.Code)

	var response map[string]string
	err = json.Unmarshal(rr.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "User registered successfully", response["message"])

	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestRegisterHandler_DuplicateEmail(t *testing.T) {
	mockEmail := email.NewMockEmailService()

	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	mock.ExpectQuery("INSERT INTO users").
		WithArgs("test@example.com", sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg()).
		WillReturnError(errors.New("email already exists"))

	authHandler := NewAuthHandler(db, mockEmail)

	userPayload := `{
        "email": "test@example.com",
        "password": "password123"
    }`

	req, err := http.NewRequest("POST", "/api/register", bytes.NewBufferString(userPayload))
	assert.NoError(t, err)

	rr := httptest.NewRecorder()
	authHandler.RegisterHandler(rr, req)

	assert.Equal(t, http.StatusConflict, rr.Code)
	assert.Contains(t, rr.Body.String(), "Email already exists")

	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestLoginHandler(t *testing.T) {
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
			payload: `{"email": "test@example.com", "password": "password123"}`,
			setupMock: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"id", "email", "username", "password", "created_at", "updated_at"}).
					AddRow(1, "test@example.com", "testuser", string(hashedPassword), time.Now(), time.Now())
				mock.ExpectQuery("SELECT (.+) FROM users WHERE email = \\$1").
					WithArgs("test@example.com").
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
			payload:      `{"email": "", "password": ""}`,
			setupMock:    func(mock sqlmock.Sqlmock) {},
			expectedCode: http.StatusBadRequest,
			expectedBody: map[string]string{
				"error": "Email and password are required",
			},
		},
		{
			name:    "Invalid credentials",
			payload: `{"email": "nonexistent@example.com", "password": "wrongpass"}`,
			setupMock: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery("SELECT (.+) FROM users WHERE email = \\$1").
					WithArgs("nonexistent@example.com").
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
			db, mock, err := sqlmock.New()
			assert.NoError(t, err)
			defer db.Close()

			tt.setupMock(mock)

			authHandler := &AuthHandler{
				DB: db,
				GenerateJWT: func(userID int, username string) (string, error) {
					return "mock-access-token", nil
				},
				GenerateRefreshToken: func(userID int, username string) (string, error) {
					return "mock-refresh-token", nil
				},
			}

			req, err := http.NewRequest("POST", "/api/login", strings.NewReader(tt.payload))
			assert.NoError(t, err)
			req.Header.Set("Content-Type", "application/json")

			rr := httptest.NewRecorder()
			authHandler.LoginHandler(rr, req)

			assert.Equal(t, tt.expectedCode, rr.Code)

			var response map[string]string
			err = json.NewDecoder(rr.Body).Decode(&response)
			assert.NoError(t, err)

			assert.Equal(t, tt.expectedBody, response)
			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}
