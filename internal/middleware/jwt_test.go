package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/stretchr/testify/assert"
)

func TestGenerateJWT(t *testing.T) {
	tests := []struct {
		name     string
		userID   int
		username string
		wantErr  bool
	}{
		{
			name:     "Valid token generation",
			userID:   1,
			username: "testuser",
			wantErr:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			token, err := GenerateJWT(tt.userID, tt.username)
			if tt.wantErr {
				assert.Error(t, err)
				assert.Empty(t, token)
			} else {
				assert.NoError(t, err)
				assert.NotEmpty(t, token)

				// Validate the generated token
				claims, err := ValidateJWT(token)
				assert.NoError(t, err)
				assert.Equal(t, tt.userID, claims.UserID)
				assert.Equal(t, tt.username, claims.Username)
			}
		})
	}
}

func TestGenerateAndValidateRefreshToken(t *testing.T) {
	tests := []struct {
		name     string
		username string
		wantErr  bool
	}{
		{
			name:     "Valid refresh token generation",
			username: "testuser",
			wantErr:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			token, err := GenerateRefreshToken(tt.username)
			if tt.wantErr {
				assert.Error(t, err)
				assert.Empty(t, token)
			} else {
				assert.NoError(t, err)
				assert.NotEmpty(t, token)

				// Validate the generated refresh token
				claims, err := ValidateRefreshToken(token)
				assert.NoError(t, err)
				assert.Equal(t, tt.username, claims.Username)
			}
		})
	}
}

func TestJWTAuthMiddleware(t *testing.T) {
	tests := []struct {
		name           string
		setupAuth      func(r *http.Request)
		expectedStatus int
	}{
		{
			name: "Valid token",
			setupAuth: func(r *http.Request) {
				token, _ := GenerateJWT(1, "testuser")
				r.Header.Set("Authorization", "Bearer "+token)
			},
			expectedStatus: http.StatusOK,
		},
		{
			name: "Missing auth header",
			setupAuth: func(r *http.Request) {
				// Don't set Authorization header
			},
			expectedStatus: http.StatusUnauthorized,
		},
		{
			name: "Invalid token",
			setupAuth: func(r *http.Request) {
				r.Header.Set("Authorization", "Bearer invalid-token")
			},
			expectedStatus: http.StatusUnauthorized,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create a test handler that will be wrapped by the middleware
			nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				// Check if claims are in context for valid tokens
				if tt.expectedStatus == http.StatusOK {
					claims, ok := r.Context().Value("claims").(*Claims)
					assert.True(t, ok)
					assert.NotNil(t, claims)
				}
				w.WriteHeader(http.StatusOK)
			})

			// Create the middleware handler
			handler := JWTAuthMiddleware(nextHandler)

			// Create test request
			req := httptest.NewRequest("GET", "/test", nil)
			tt.setupAuth(req)

			// Create response recorder
			rr := httptest.NewRecorder()

			// Serve the request
			handler.ServeHTTP(rr, req)

			// Check status code
			assert.Equal(t, tt.expectedStatus, rr.Code)
		})
	}
}

func TestTokenExpiration(t *testing.T) {
	// Generate a token that expires in 1 second
	claims := &Claims{
		UserID:   1,
		Username: "testuser",
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(1 * time.Second).Unix(),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(jwtSecret)
	assert.NoError(t, err)

	// Wait for token to expire
	time.Sleep(2 * time.Second)

	// Try to validate expired token
	_, err = ValidateJWT(tokenString)
	assert.Error(t, err)
}
