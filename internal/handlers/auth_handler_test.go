package handlers

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/maxzhirnov/go-task-manager/internal/middleware"
	"github.com/maxzhirnov/go-task-manager/pkg/email"
	"github.com/stretchr/testify/assert"
	"golang.org/x/crypto/bcrypt"
)

func newTestAuthHandler(t *testing.T) (*AuthHandler, sqlmock.Sqlmock, func()) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Failed to create mock DB: %v", err)
	}

	mockEmail := email.NewMockEmailService()
	handler := &AuthHandler{
		DB:           db,
		EmailService: mockEmail,
		GenerateJWT: func(userID int, username string, email string) (string, error) {
			return "mock-access-token", nil
		},
		GenerateRefreshToken: func(userID int, username string, email string) (string, error) {
			return "mock-refresh-token", nil
		},
	}

	cleanup := func() {
		db.Close()
	}

	return handler, mock, cleanup
}

// Helper function to create a test request
func createTestRequest(t *testing.T, method, endpoint string, body interface{}) *http.Request {
	var reqBody io.Reader
	if body != nil {
		jsonBody, err := json.Marshal(body)
		if err != nil {
			t.Fatalf("Failed to marshal request body: %v", err)
		}
		reqBody = bytes.NewBuffer(jsonBody)
	}

	req, err := http.NewRequest(method, endpoint, reqBody)
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")
	return req
}

func TestRegisterHandlerSuccess(t *testing.T) {
	handler, mock, cleanup := newTestAuthHandler(t)
	defer cleanup()

	// Begin transaction
	mock.ExpectBegin()

	// Mock user insertion with exact columns and values
	mock.ExpectQuery(`INSERT INTO users \(email, username, password, is_verified, created_at, updated_at\) VALUES \(\$1, \$2, \$3, \$4, \$5, \$6\) RETURNING id`).
		WithArgs(
			"test@example.com", // email
			"test",             // username (derived from email)
			sqlmock.AnyArg(),   // hashed password
			false,              // is_verified
			sqlmock.AnyArg(),   // created_at
			sqlmock.AnyArg(),   // updated_at
		).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))

	// Mock verification token insertion
	mock.ExpectExec(`INSERT INTO verification_tokens \(user_id, token, expires_at, created_at\) VALUES \(\$1, \$2, \$3, \$4\)`).
		WithArgs(
			1,                // user_id
			sqlmock.AnyArg(), // token
			sqlmock.AnyArg(), // expires_at
			sqlmock.AnyArg(), // created_at
		).
		WillReturnResult(sqlmock.NewResult(1, 1))

	// Expect commit
	mock.ExpectCommit()

	req := createTestRequest(t, "POST", "/api/register", RegisterRequest{
		Email:    "test@example.com",
		Password: "password123",
	})

	rr := httptest.NewRecorder()
	handler.RegisterHandler(rr, req)

	assert.Equal(t, http.StatusCreated, rr.Code)

	var response map[string]string
	err := json.NewDecoder(rr.Body).Decode(&response)
	assert.NoError(t, err)
	assert.Equal(t, "User registered successfully", response["message"])

	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestRegisterHandler(t *testing.T) {
	tests := []struct {
		name         string
		payload      interface{}
		setupMock    func(mock sqlmock.Sqlmock)
		expectedCode int
		expectedBody map[string]string
	}{
		{
			name:         "Invalid JSON body",
			payload:      "invalid json",
			setupMock:    func(mock sqlmock.Sqlmock) {},
			expectedCode: http.StatusBadRequest,
			expectedBody: map[string]string{
				"error": "Invalid request body",
			},
		},
		{
			name: "Empty email",
			payload: RegisterRequest{
				Email:    "",
				Password: "password123",
			},
			setupMock:    func(mock sqlmock.Sqlmock) {},
			expectedCode: http.StatusBadRequest,
			expectedBody: map[string]string{
				"error": "Email and password are required",
			},
		},
		{
			name: "Empty password",
			payload: RegisterRequest{
				Email:    "test@example.com",
				Password: "",
			},
			setupMock:    func(mock sqlmock.Sqlmock) {},
			expectedCode: http.StatusBadRequest,
			expectedBody: map[string]string{
				"error": "Email and password are required",
			},
		},
		{
			name: "Duplicate email",
			payload: RegisterRequest{
				Email:    "existing@example.com",
				Password: "password123",
			},
			setupMock: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectQuery(`INSERT INTO users`).
					WillReturnError(fmt.Errorf("email already exists"))
				mock.ExpectRollback()
			},
			expectedCode: http.StatusConflict,
			expectedBody: map[string]string{
				"error": "Email already exists",
			},
		},
		{
			name: "Database error",
			payload: RegisterRequest{
				Email:    "test@example.com",
				Password: "password123",
			},
			setupMock: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectQuery(`INSERT INTO users`).
					WillReturnError(fmt.Errorf("database error"))
				mock.ExpectRollback()
			},
			expectedCode: http.StatusInternalServerError,
			expectedBody: map[string]string{
				"error": "Error creating user",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup
			handler, mock, cleanup := newTestAuthHandler(t)
			defer cleanup()

			// Configure mock expectations
			tt.setupMock(mock)

			// Create request
			var body []byte
			var err error
			if str, ok := tt.payload.(string); ok {
				body = []byte(str)
			} else {
				body, err = json.Marshal(tt.payload)
				assert.NoError(t, err)
			}

			req := httptest.NewRequest("POST", "/api/register", bytes.NewBuffer(body))
			req.Header.Set("Content-Type", "application/json")
			rr := httptest.NewRecorder()

			// Execute request
			handler.RegisterHandler(rr, req)

			// Assert response
			assert.Equal(t, tt.expectedCode, rr.Code)

			var response map[string]string
			err = json.NewDecoder(rr.Body).Decode(&response)
			assert.NoError(t, err)
			assert.Equal(t, tt.expectedBody, response)

			// Verify mock expectations
			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestLoginHandler(t *testing.T) {
	tests := []struct {
		name           string
		payload        interface{}
		setupMock      func(mock sqlmock.Sqlmock)
		expectedStatus int
		expectedBody   map[string]string
	}{
		{
			name: "Successful login",
			payload: LoginRequest{
				Email:    "test@example.com",
				Password: "password123",
			},
			setupMock: func(mock sqlmock.Sqlmock) {
				hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("password123"), bcrypt.DefaultCost)
				rows := sqlmock.NewRows([]string{
					"id", "email", "username", "password", "is_verified", "created_at", "updated_at",
				}).AddRow(1, "test@example.com", "testuser", string(hashedPassword), true, time.Now(), time.Now())

				mock.ExpectQuery("SELECT (.+) FROM users WHERE email = \\$1").
					WithArgs("test@example.com").
					WillReturnRows(rows)
			},
			expectedStatus: http.StatusOK,
			expectedBody: map[string]string{
				"access_token":  "mock-access-token",
				"refresh_token": "mock-refresh-token",
			},
		},
		{
			name:           "Invalid JSON body",
			payload:        "invalid json",
			setupMock:      func(mock sqlmock.Sqlmock) {},
			expectedStatus: http.StatusBadRequest,
			expectedBody: map[string]string{
				"error": "Invalid request body",
			},
		},
		{
			name: "Empty email",
			payload: LoginRequest{
				Email:    "",
				Password: "password123",
			},
			setupMock:      func(mock sqlmock.Sqlmock) {},
			expectedStatus: http.StatusBadRequest,
			expectedBody: map[string]string{
				"error": "Email and password are required",
			},
		},
		{
			name: "Empty password",
			payload: LoginRequest{
				Email:    "test@example.com",
				Password: "",
			},
			setupMock:      func(mock sqlmock.Sqlmock) {},
			expectedStatus: http.StatusBadRequest,
			expectedBody: map[string]string{
				"error": "Email and password are required",
			},
		},
		{
			name: "User not found",
			payload: LoginRequest{
				Email:    "nonexistent@example.com",
				Password: "password123",
			},
			setupMock: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery("SELECT (.+) FROM users WHERE email = \\$1").
					WithArgs("nonexistent@example.com").
					WillReturnError(sql.ErrNoRows)
			},
			expectedStatus: http.StatusUnauthorized,
			expectedBody: map[string]string{
				"error": "Invalid credentials",
			},
		},
		{
			name: "Unverified email",
			payload: LoginRequest{
				Email:    "unverified@example.com",
				Password: "password123",
			},
			setupMock: func(mock sqlmock.Sqlmock) {
				hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("password123"), bcrypt.DefaultCost)
				rows := sqlmock.NewRows([]string{
					"id", "email", "username", "password", "is_verified", "created_at", "updated_at",
				}).AddRow(1, "unverified@example.com", "testuser", string(hashedPassword), false, time.Now(), time.Now())

				mock.ExpectQuery("SELECT (.+) FROM users WHERE email = \\$1").
					WithArgs("unverified@example.com").
					WillReturnRows(rows)
			},
			expectedStatus: http.StatusForbidden,
			expectedBody: map[string]string{
				"error": "Please verify your email before logging in",
			},
		},
		{
			name: "Wrong password",
			payload: LoginRequest{
				Email:    "test@example.com",
				Password: "wrongpassword",
			},
			setupMock: func(mock sqlmock.Sqlmock) {
				hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("password123"), bcrypt.DefaultCost)
				rows := sqlmock.NewRows([]string{
					"id", "email", "username", "password", "is_verified", "created_at", "updated_at",
				}).AddRow(1, "test@example.com", "testuser", string(hashedPassword), true, time.Now(), time.Now())

				mock.ExpectQuery("SELECT (.+) FROM users WHERE email = \\$1").
					WithArgs("test@example.com").
					WillReturnRows(rows)
			},
			expectedStatus: http.StatusUnauthorized,
			expectedBody: map[string]string{
				"error": "Invalid credentials",
			},
		},
		{
			name: "Database error",
			payload: LoginRequest{
				Email:    "test@example.com",
				Password: "password123",
			},
			setupMock: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery("SELECT (.+) FROM users WHERE email = \\$1").
					WithArgs("test@example.com").
					WillReturnError(fmt.Errorf("database error"))
			},
			expectedStatus: http.StatusInternalServerError,
			expectedBody: map[string]string{
				"error": "Server error",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			handler, mock, cleanup := newTestAuthHandler(t)
			defer cleanup()

			tt.setupMock(mock)

			var req *http.Request
			if str, ok := tt.payload.(string); ok {
				req = createTestRequest(t, "POST", "/api/login", str)
			} else {
				req = createTestRequest(t, "POST", "/api/login", tt.payload)
			}

			rr := httptest.NewRecorder()
			handler.LoginHandler(rr, req)

			assert.Equal(t, tt.expectedStatus, rr.Code)

			var response map[string]string
			err := json.NewDecoder(rr.Body).Decode(&response)
			assert.NoError(t, err)
			assert.Equal(t, tt.expectedBody, response)

			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestRefreshTokenHandler(t *testing.T) {
	tests := []struct {
		name           string
		payload        interface{}
		setupHandler   func() *AuthHandler
		expectedStatus int
		expectedBody   map[string]string
	}{
		// {
		// 	name: "Successful token refresh",
		// 	payload: struct {
		// 		RefreshToken string `json:"refresh_token"`
		// 	}{
		// 		RefreshToken: "valid-refresh-token",
		// 	},
		// 	setupHandler: func() *AuthHandler {
		// 		db, _, err := sqlmock.New()
		// 		if err != nil {
		// 			t.Fatal(err)
		// 		}

		// 		return &AuthHandler{
		// 			DB: db,
		// 			ValidateRefreshToken: func(token string) (*middleware.Claims, error) {
		// 				return &middleware.Claims{
		// 					UserID:   1,
		// 					Username: "testuser",
		// 				}, nil
		// 			},
		// 			GenerateJWT: func(userID int, username string) (string, error) {
		// 				return "mock-access-token", nil
		// 			},
		// 		}
		// 	},
		// 	expectedStatus: http.StatusOK,
		// 	expectedBody: map[string]string{
		// 		"access_token": "mock-access-token",
		// 	},
		// },
		{
			name:    "Invalid JSON body",
			payload: "invalid json",
			setupHandler: func() *AuthHandler {
				db, _, err := sqlmock.New()
				if err != nil {
					t.Fatal(err)
				}
				return &AuthHandler{DB: db}
			},
			expectedStatus: http.StatusBadRequest,
			expectedBody: map[string]string{
				"error": "Invalid input",
			},
		},
		{
			name: "Invalid refresh token",
			payload: struct {
				RefreshToken string `json:"refresh_token"`
			}{
				RefreshToken: "invalid-token",
			},
			setupHandler: func() *AuthHandler {
				db, _, err := sqlmock.New()
				if err != nil {
					t.Fatal(err)
				}
				return &AuthHandler{
					DB: db,
					ValidateRefreshToken: func(token string) (*middleware.Claims, error) {
						return nil, fmt.Errorf("invalid token")
					},
				}
			},
			expectedStatus: http.StatusUnauthorized,
			expectedBody: map[string]string{
				"error": "Invalid refresh token",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			handler := tt.setupHandler()

			var req *http.Request
			if str, ok := tt.payload.(string); ok {
				req = createTestRequest(t, "POST", "/api/refresh", str)
			} else {
				req = createTestRequest(t, "POST", "/api/refresh", tt.payload)
			}

			rr := httptest.NewRecorder()
			handler.RefreshTokenHandler(rr, req)

			assert.Equal(t, tt.expectedStatus, rr.Code)

			var response map[string]string
			err := json.NewDecoder(rr.Body).Decode(&response)
			assert.NoError(t, err)
			assert.Equal(t, tt.expectedBody, response)
		})
	}
}

func TestVerifyEmailHandler(t *testing.T) {
	tests := []struct {
		name           string
		token          string
		setupMock      func(mock sqlmock.Sqlmock)
		expectedStatus int
		expectedBody   map[string]string
	}{
		{
			name:  "Successful verification",
			token: "valid-token",
			setupMock: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectQuery("UPDATE verification_tokens").
					WithArgs("valid-token").
					WillReturnRows(sqlmock.NewRows([]string{"user_id"}).AddRow(1))
				mock.ExpectExec("UPDATE users SET is_verified = true").
					WithArgs(1).
					WillReturnResult(sqlmock.NewResult(0, 1))
				mock.ExpectCommit()
			},
			expectedStatus: http.StatusOK,
			expectedBody: map[string]string{
				"message": "Email verified successfully",
			},
		},
		{
			name:           "Missing token",
			token:          "",
			setupMock:      func(mock sqlmock.Sqlmock) {},
			expectedStatus: http.StatusBadRequest,
			expectedBody: map[string]string{
				"error": "Verification token is required",
			},
		},
		{
			name:  "Invalid or expired token",
			token: "invalid-token",
			setupMock: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectQuery("UPDATE verification_tokens").
					WithArgs("invalid-token").
					WillReturnError(fmt.Errorf("invalid or expired verification token"))
				mock.ExpectRollback()
			},
			expectedStatus: http.StatusBadRequest,
			expectedBody: map[string]string{
				"error": "Invalid or expired verification token",
			},
		},
		{
			name:  "Database error",
			token: "valid-token",
			setupMock: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectQuery("UPDATE verification_tokens").
					WithArgs("valid-token").
					WillReturnError(fmt.Errorf("database error"))
				mock.ExpectRollback()
			},
			expectedStatus: http.StatusInternalServerError,
			expectedBody: map[string]string{
				"error": "Error verifying email",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			handler, mock, cleanup := newTestAuthHandler(t)
			defer cleanup()

			tt.setupMock(mock)

			req := httptest.NewRequest("GET", "/api/verify-email?token="+tt.token, nil)
			rr := httptest.NewRecorder()

			handler.VerifyEmailHandler(rr, req)

			assert.Equal(t, tt.expectedStatus, rr.Code)

			var response map[string]string
			err := json.NewDecoder(rr.Body).Decode(&response)
			assert.NoError(t, err)
			assert.Equal(t, tt.expectedBody, response)

			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}
