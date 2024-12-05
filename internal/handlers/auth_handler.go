// @title Task Manager API
// @version 1.0
// @description Task management system with JWT authentication
// @host localhost:8080
// @BasePath /api
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
package handlers

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"strings"

	"github.com/maxzhirnov/go-task-manager/internal/middleware"
	"github.com/maxzhirnov/go-task-manager/internal/models"
	"github.com/maxzhirnov/go-task-manager/pkg/database"
	"github.com/maxzhirnov/go-task-manager/pkg/email"
)

// AuthHandler manages authentication-related HTTP requests.
// It handles user registration, login, token refresh, and email verification.
type AuthHandler struct {
	// DB provides database access for user operations
	DB database.DB

	// EmailService handles sending verification and notification emails
	EmailService email.EmailSender

	// GenerateJWT creates new JWT access tokens
	GenerateJWT func(userID int, username string) (string, error)

	// GenerateRefreshToken creates new refresh tokens
	GenerateRefreshToken func(userID int, username string) (string, error)

	// ValidateRefreshToken verifies and parses refresh tokens
	ValidateRefreshToken func(token string) (*middleware.Claims, error)
}

// NewAuthHandler creates a new instance of AuthHandler with default token handlers.
//
// Parameters:
//   - db: Database interface for user operations
//   - emailService: Service for sending emails
//
// Returns:
//   - *AuthHandler: Configured authentication handler
func NewAuthHandler(db database.DB, emailService email.EmailSender) *AuthHandler {
	return &AuthHandler{
		DB:                   db,
		EmailService:         emailService,
		GenerateJWT:          middleware.GenerateJWT,
		GenerateRefreshToken: middleware.GenerateRefreshToken,
		ValidateRefreshToken: middleware.ValidateRefreshToken,
	}
}

// RegisterRequest represents the expected JSON structure for registration requests.
type RegisterRequest struct {
	// Email address of the new user
	Email string `json:"email"`

	// Password for the new account (will be hashed before storage)
	Password string `json:"password"`
}

// RegisterHandler processes new user registration requests.
//
// It validates the registration input, creates a new user account,
// and initiates the email verification process. The handler expects
// a JSON request body containing email and password fields.
//
// The registration process includes:
// 1. Input validation
// 2. Password hashing
// 3. User creation in database
// 4. Verification token generation
// 5. Verification email sending
//
// HTTP Responses:
//   - 201 Created: Successful registration
//   - 400 Bad Request: Invalid input or missing required fields
//   - 409 Conflict: Email already exists
//   - 500 Internal Server Error: Server-side errors
//
// Example request:
//
//	POST /api/register
//	{
//	    "email": "user@example.com",
//	    "password": "userpassword"
//	}
//
// Example success response:
//
//	{
//	    "message": "User registered successfully"
//	}
func (h *AuthHandler) RegisterHandler(w http.ResponseWriter, r *http.Request) {
	// Parse and validate request body
	var req RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		JSONError(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Ensure required fields are present
	if req.Email == "" || req.Password == "" {
		JSONError(w, "Email and password are required", http.StatusBadRequest)
		return
	}

	// Initialize user model with request data
	user := &models.User{
		Email:    req.Email,
		Password: req.Password,
	}

	// Hash the user's password before storage
	if err := user.HashPassword(); err != nil {
		log.Printf("Failed to hash password: %v", err)
		JSONError(w, "Error hashing password", http.StatusInternalServerError)
		return
	}

	// Attempt to create the user in the database
	if err := user.CreateUser(h.DB); err != nil {
		// Handle duplicate email conflict
		if strings.Contains(err.Error(), "email already exists") {
			JSONError(w, "Email already exists", http.StatusConflict)
			return
		}
		// Handle other database errors
		log.Printf("Failed to create user: %v", err)
		JSONError(w, "Error creating user", http.StatusInternalServerError)
		return
	}

	// Generate and send verification email
	// Note: Registration continues even if verification steps fail
	token, err := models.GetVerificationTokenForUser(h.DB, user.ID)
	if err != nil {
		log.Printf("Failed to get verification token: %v", err)
	} else {
		if err := h.EmailService.SendVerificationEmail(user.Email, user.Username, token); err != nil {
			log.Printf("Failed to send verification email: %v", err)
		}
	}

	// Return success response
	log.Printf("User registered successfully: %s", user.Email)
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{
		"message": "User registered successfully",
	})
}

// LoginRequest represents the expected JSON structure for login requests.
type LoginRequest struct {
	// Email address of the user attempting to log in
	Email string `json:"email"`

	// Password to verify against stored hash
	Password string `json:"password"`
}

// LoginHandler processes user login requests and issues JWT tokens.
//
// It validates user credentials, checks email verification status,
// and generates both access and refresh tokens upon successful authentication.
//
// The login process includes:
// 1. Credential validation
// 2. Email verification check
// 3. Password verification
// 4. Token generation (access and refresh)
//
// HTTP Responses:
//   - 200 OK: Successful login with tokens
//   - 400 Bad Request: Invalid input or missing fields
//   - 401 Unauthorized: Invalid credentials
//   - 403 Forbidden: Email not verified
//   - 500 Internal Server Error: Server-side errors
//
// Example request:
//
//	POST /api/login
//	{
//	    "email": "user@example.com",
//	    "password": "userpassword"
//	}
//
// Example success response:
//
//	{
//	    "access_token": "eyJhbGc...",
//	    "refresh_token": "eyJhbGc..."
//	}
func (h *AuthHandler) LoginHandler(w http.ResponseWriter, r *http.Request) {
	// Parse and validate request body
	var req LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		JSONError(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Validate required fields
	if req.Email == "" || req.Password == "" {
		JSONError(w, "Email and password are required", http.StatusBadRequest)
		return
	}

	// Retrieve user from database
	user, err := models.GetUserByEmail(h.DB, req.Email)
	if err != nil {
		if err == sql.ErrNoRows {
			// User not found, return generic error for security
			JSONError(w, "Invalid credentials", http.StatusUnauthorized)
			return
		}
		log.Printf("Database error during login: %v", err)
		JSONError(w, "Server error", http.StatusInternalServerError)
		return
	}

	// Ensure email is verified
	if !user.IsVerified {
		JSONError(w, "Please verify your email before logging in", http.StatusForbidden)
		return
	}

	// Verify password
	if err := user.CheckPassword(req.Password); err != nil {
		log.Printf("Failed password check for user %s: %v", user.Email, err)
		JSONError(w, "Invalid credentials", http.StatusUnauthorized)
		return
	}

	// Generate access token
	accessToken, err := h.GenerateJWT(user.ID, user.Username)
	if err != nil {
		log.Printf("Failed to generate access token: %v", err)
		JSONError(w, "Failed to generate access token", http.StatusInternalServerError)
		return
	}

	// Generate refresh token
	refreshToken, err := h.GenerateRefreshToken(user.ID, user.Username)
	if err != nil {
		log.Printf("Failed to generate refresh token: %v", err)
		JSONError(w, "Failed to generate refresh token", http.StatusInternalServerError)
		return
	}

	// Send successful response with tokens
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"access_token":  accessToken,
		"refresh_token": refreshToken,
	})
}

// RefreshTokenRequest represents the expected JSON structure for token refresh requests.
type RefreshTokenRequest struct {
	// RefreshToken is the token used to obtain a new access token
	RefreshToken string `json:"refresh_token"`
}

// RefreshTokenHandler processes requests to refresh expired access tokens.
//
// It validates the provided refresh token and generates a new access token
// if the refresh token is valid. This endpoint is used to maintain user
// sessions without requiring re-authentication.
//
// The refresh process includes:
// 1. Refresh token validation
// 2. Claims extraction
// 3. New access token generation
//
// HTTP Responses:
//   - 200 OK: Successfully generated new access token
//   - 400 Bad Request: Invalid request format
//   - 401 Unauthorized: Invalid or expired refresh token
//   - 500 Internal Server Error: Token generation failure
//
// Example request:
//
//	POST /api/refresh
//	{
//	    "refresh_token": "eyJhbGc..."
//	}
//
// Example success response:
//
//	{
//	    "access_token": "eyJhbGc..."
//	}
func (h *AuthHandler) RefreshTokenHandler(w http.ResponseWriter, r *http.Request) {
	// Parse and validate request body
	var req RefreshTokenRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Printf("Failed to decode refresh token request: %v", err)
		JSONError(w, "Invalid input", http.StatusBadRequest)
		return
	}

	// Validate refresh token and extract claims
	claims, err := h.ValidateRefreshToken(req.RefreshToken)
	if err != nil {
		log.Printf("Invalid refresh token: %v", err)
		JSONError(w, "Invalid refresh token", http.StatusUnauthorized)
		return
	}

	// Fetch latest user data from database
	user, err := models.GetUserByID(h.DB, claims.UserID)
	if err != nil {
		log.Printf("Failed to fetch user data: %v", err)
		JSONError(w, "User not found", http.StatusNotFound)
		return
	}

	// Generate new access token using latest user data
	accessToken, err := h.GenerateJWT(user.ID, user.Username)
	if err != nil {
		log.Printf("Failed to generate new access token: %v", err)
		JSONError(w, "Failed to generate access token", http.StatusInternalServerError)
		return
	}

	// Send successful response with new access token
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"access_token": accessToken,
	})

	log.Printf("Successfully refreshed token for user ID: %d", claims.UserID)
}

// VerifyEmailHandler processes email verification requests.
//
// It validates the verification token provided in the URL query parameters
// and updates the user's email verification status. This endpoint is typically
// accessed via a link sent to the user's email during registration.
//
// The verification process includes:
// 1. Token validation
// 2. Token expiration check
// 3. User verification status update
//
// URL Parameters:
//   - token: The verification token sent to user's email
//
// HTTP Responses:
//   - 200 OK: Email successfully verified
//   - 400 Bad Request: Missing, invalid, or expired token
//   - 500 Internal Server Error: Server-side verification errors
//
// Example request:
//
//	GET /api/verify-email?token=abc123...
//
// Example success response:
//
//	{
//	    "message": "Email verified successfully"
//	}
func (h *AuthHandler) VerifyEmailHandler(w http.ResponseWriter, r *http.Request) {
	// Extract verification token from URL query parameters
	token := r.URL.Query().Get("token")
	if token == "" {
		log.Printf("Email verification attempted without token")
		JSONError(w, "Verification token is required", http.StatusBadRequest)
		return
	}

	// Attempt to verify the email using the provided token
	err := models.VerifyEmail(h.DB, token)
	if err != nil {
		if strings.Contains(err.Error(), "invalid or expired") {
			// Handle invalid or expired token
			log.Printf("Invalid or expired verification token: %v", err)
			JSONError(w, "Invalid or expired verification token", http.StatusBadRequest)
			return
		}
		// Handle other verification errors
		log.Printf("Email verification failed: %v", err)
		JSONError(w, "Error verifying email", http.StatusInternalServerError)
		return
	}

	// Send successful verification response
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{
		"message": "Email verified successfully",
	})

	log.Printf("Email verification successful for token: %s", token)
}

// ResendVerificationRequest represents the expected JSON structure for verification resend requests.
type ResendVerificationRequest struct {
	// Email address of the user requesting verification resend
	Email string `json:"email"`
}

// ResendVerificationHandler processes requests to resend email verification links.
//
// It checks if the user exists and isn't already verified, generates a new
// verification token, and sends a new verification email. This endpoint
// helps users who didn't receive or lost their original verification email.
//
// The resend process includes:
// 1. User existence verification
// 2. Verification status check
// 3. New token generation
// 4. Verification email sending
//
// HTTP Responses:
//   - 200 OK: Verification email successfully sent
//   - 400 Bad Request: Invalid request or already verified email
//   - 404 Not Found: User not found
//   - 500 Internal Server Error: Token generation or email sending errors
//
// Example request:
//
//	POST /api/resend-verification
//	{
//	    "email": "user@example.com"
//	}
//
// Example success response:
//
//	{
//	    "message": "Verification email sent successfully"
//	}
func (h *AuthHandler) ResendVerificationHandler(w http.ResponseWriter, r *http.Request) {
	// Parse and validate request body
	var req ResendVerificationRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Printf("Failed to decode resend verification request: %v", err)
		JSONError(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Retrieve user from database
	user, err := models.GetUserByEmail(h.DB, req.Email)
	if err != nil {
		log.Printf("User not found for verification resend: %s", req.Email)
		JSONError(w, "User not found", http.StatusNotFound)
		return
	}

	// Check if email is already verified
	if user.IsVerified {
		log.Printf("Attempted to resend verification for already verified email: %s", req.Email)
		JSONError(w, "Email is already verified", http.StatusBadRequest)
		return
	}

	// Generate new verification token
	token, err := models.CreateVerificationToken(h.DB, user.ID)
	if err != nil {
		log.Printf("Failed to generate verification token: %v", err)
		JSONError(w, "Error generating verification token", http.StatusInternalServerError)
		return
	}

	// Send verification email
	if err := h.EmailService.SendVerificationEmail(user.Email, user.Username, token.Token); err != nil {
		log.Printf("Failed to send verification email: %v", err)
		JSONError(w, "Error sending verification email", http.StatusInternalServerError)
		return
	}

	// Send successful response
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{
		"message": "Verification email sent successfully",
	})

	log.Printf("Verification email resent successfully to: %s", req.Email)
}
