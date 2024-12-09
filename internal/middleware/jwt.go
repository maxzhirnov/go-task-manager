// Package middleware provides JWT authentication and authorization functionality
// for the task management application. It handles token generation, validation,
// and user context management.
package middleware

import (
	"context"
	"errors"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
)

// JWT configuration variables
var (
	// jwtSecret is used for signing and validating access tokens.
	// It should be set via JWT_SECRET environment variable in production.
	jwtSecret = []byte(getEnvWithDefault("JWT_SECRET", "default_secret_key_please_change_in_production"))
	// jwtRefreshSecret is used for signing and validating refresh tokens.
	// It should be set via JWT_REFRESH_SECRET environment variable in
	jwtRefreshSecret = []byte(getEnvWithDefault("JWT_REFRESH_SECRET", "default_refresh_secret_key_please_change_in_production"))
)

// getEnvWithDefault retrieves an environment variable value or returns a default if not set.
//
// Parameters:
//   - key: The environment variable name to look up
//   - defaultValue: The value to return if the environment variable is not set
//
// Returns:
//   - string: The environment variable value or the default value
//
// Example Usage:
//
//	secretKey := getEnvWithDefault("JWT_SECRET", "default_secret")
func getEnvWithDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

// Claims represents the custom JWT claims structure used for both access
// and refresh tokens. It extends jwt.StandardClaims to include user-specific
// information.
//
// Fields:
//   - UserID: The unique identifier of the authenticated user
//   - Username: The username of the authenticated user
//   - StandardClaims: Standard JWT claims (exp, iat, etc.)
//
// Note: This structure is used for both token generation and validation.
type Claims struct {
	UserID   int    `json:"user_id"`
	Username string `json:"username"`
	Email    string `json:"email"`
	jwt.StandardClaims
}

// GenerateJWT creates a new JWT access token for a user.
//
// It generates a signed JWT token containing user identification information
// and standard claims. The token is signed using HMAC-SHA256 (HS256) algorithm
// and expires after 1 hour from creation.
//
// Parameters:
//   - userID: The unique identifier of the user
//   - username: The username of the user
//
// Returns:
//   - string: The signed JWT token string
//   - error: An error if token generation fails
//
// Token Structure:
//
//	Header: {
//	  "alg": "HS256",
//	  "typ": "JWT"
//	}
//	Payload: {
//	  "user_id": 123,
//	  "username": "john_doe",
//	  "exp": 1516239022
//	}
//
// Example Usage:
//
//	token, err := GenerateJWT(user.ID, user.Username)
//	if err != nil {
//	    log.Printf("Failed to generate token: %v", err)
//	    return err
//	}
//
// Security Note:
//   - The token is signed with jwtSecret
//   - Expires in 1 hour from creation
//   - Contains user identification but no sensitive data
func GenerateJWT(userID int, username string, email string) (string, error) {
	// Set token expiration time to 1 hour from now
	expirationTime := time.Now().Add(1 * time.Hour)

	// Create claims with user information and expiration
	claims := &Claims{
		UserID:   userID,
		Username: username,
		Email:    email,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}

	// Create new token with claims using HS256 signing method
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Sign and return the token
	return token.SignedString(jwtSecret)
}

// GenerateRefreshToken creates a new JWT refresh token for a user.
//
// It generates a long-lived token (7 days) that can be used to obtain new
// access tokens without requiring re-authentication. The refresh token uses
// a separate secret key from the access token for additional security.
//
// Parameters:
//   - userID: The unique identifier of the user
//   - username: The username derived from user's email
//
// Returns:
//   - string: The signed refresh token
//   - error: An error if token generation fails
//
// Example Usage:
//
//	refreshToken, err := GenerateRefreshToken(user.ID, user.Username)
//	if err != nil {
//	    return "", fmt.Errorf("failed to generate refresh token: %w", err)
//	}
func GenerateRefreshToken(userID int, username string, email string) (string, error) {
	// Set expiration time to 7 days from now
	expirationTime := time.Now().Add(7 * 24 * time.Hour)

	// Create claims with user information and expiration
	claims := &Claims{
		UserID:   userID,
		Username: username,
		Email:    email,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}

	// Create and sign token using refresh secret
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtRefreshSecret)
}

// ValidateJWT validates an access token and extracts its claims.
//
// It verifies the token's signature using the JWT secret key and checks
// its validity. This function is used to authenticate requests by validating
// the access token provided in the Authorization header.
//
// Parameters:
//   - tokenString: The JWT access token to validate
//
// Returns:
//   - *Claims: The token's claims if validation is successful
//   - error: An error if the token is invalid, expired, or malformed
//
// Error cases:
//   - Token is expired
//   - Token has invalid signature
//   - Token is malformed
//   - Token claims cannot be parsed
//
// Example Usage:
//
//	claims, err := ValidateJWT(tokenString)
//	if err != nil {
//	    return nil, fmt.Errorf("invalid token: %w", err)
//	}
//	userID := claims.UserID
//
// Security Note:
//   - Uses jwtSecret for validation
//   - Automatically checks token expiration
//   - Verifies token signature
func ValidateJWT(tokenString string) (*Claims, error) {
	// Parse and validate the token with claims
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return jwtSecret, nil
	})

	// Handle parsing errors (expired, invalid signature, malformed)
	if err != nil {
		return nil, err
	}

	// Extract and validate claims
	claims, ok := token.Claims.(*Claims)
	if !ok || !token.Valid {
		return nil, errors.New("invalid token")
	}

	return claims, nil
}

// ValidateRefreshToken validates a refresh token and extracts its claims.
//
// It verifies the refresh token's signature using a separate refresh secret key
// and checks its validity. This function is used during the token refresh process
// to validate refresh tokens before issuing new access tokens.
//
// Parameters:
//   - tokenString: The refresh token to validate
//
// Returns:
//   - *Claims: The token's claims if validation is successful
//   - error: jwt.ErrSignatureInvalid or parsing error if validation fails
//
// Error cases:
//   - Token is expired
//   - Token has invalid signature
//   - Token is malformed
//   - Token claims cannot be parsed
//
// Example Usage:
//
//	claims, err := ValidateRefreshToken(refreshToken)
//	if err != nil {
//	    if err == jwt.ErrSignatureInvalid {
//	        return nil, fmt.Errorf("invalid refresh token signature")
//	    }
//	    return nil, fmt.Errorf("refresh token validation failed: %w", err)
//	}
//
// Security Note:
//   - Uses separate jwtRefreshSecret for validation
//   - Returns specific error for invalid signatures
//   - Automatically checks token expiration
//   - Used only for refresh token operations
func ValidateRefreshToken(tokenString string) (*Claims, error) {
	// Parse and validate the refresh token with claims
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return jwtRefreshSecret, nil
	})

	// Handle parsing errors (expired, malformed)
	if err != nil {
		return nil, err
	}

	// Extract and validate claims
	claims, ok := token.Claims.(*Claims)
	if !ok || !token.Valid {
		return nil, jwt.ErrSignatureInvalid
	}

	return claims, nil
}

// JWTAuthMiddleware provides JWT authentication for HTTP endpoints.
//
// This middleware:
// 1. Extracts the JWT from the Authorization header
// 2. Validates the token
// 3. Adds the claims to the request context
// 4. Passes the request to the next handler if authentication succeeds
//
// Authorization Header Format:
//
//	Authorization: Bearer <token>
//
// Context Value:
//
//	Key: "claims"
//	Value: *Claims containing UserID and Username
//
// HTTP Responses:
//   - 401 Unauthorized:
//   - Missing Authorization header
//   - Invalid token format
//   - Expired token
//   - Invalid signature
//
// Example Usage:
//
//	router.Handle("/api/protected",
//	    JWTAuthMiddleware(http.HandlerFunc(protectedHandler)))
//
// Protected Handler Access:
//
//	func protectedHandler(w http.ResponseWriter, r *http.Request) {
//	    claims := r.Context().Value("claims").(*Claims)
//	    userID := claims.UserID
//	    // ... handler logic
//	}
func JWTAuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Extract Authorization header
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, "Authorization header required", http.StatusUnauthorized)
			return
		}

		// Extract token from Bearer scheme
		tokenString := strings.TrimPrefix(authHeader, "Bearer ")

		// Validate token and extract claims
		claims, err := ValidateJWT(tokenString)
		if err != nil {
			http.Error(w, "Invalid token", http.StatusUnauthorized)
			return
		}

		// Add claims to request context
		ctx := context.WithValue(r.Context(), "claims", claims)
		r = r.WithContext(ctx)

		// Pass to next handler
		next.ServeHTTP(w, r)
	})
}
