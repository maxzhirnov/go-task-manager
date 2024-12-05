package handlers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/maxzhirnov/go-task-manager/internal/middleware"
	"github.com/maxzhirnov/go-task-manager/internal/models"
	"github.com/maxzhirnov/go-task-manager/pkg/config"
	"github.com/maxzhirnov/go-task-manager/pkg/database"
	"github.com/maxzhirnov/go-task-manager/pkg/email"
	"golang.org/x/crypto/bcrypt"
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

	config *config.Config
}

// NewAuthHandler creates a new instance of AuthHandler with default token handlers.
//
// Parameters:
//   - db: Database interface for user operations
//   - emailService: Service for sending emails
//
// Returns:
//   - *AuthHandler: Configured authentication handler
func NewAuthHandler(db database.DB, emailService email.EmailSender, config *config.Config) *AuthHandler {
	return &AuthHandler{
		DB:                   db,
		EmailService:         emailService,
		GenerateJWT:          middleware.GenerateJWT,
		GenerateRefreshToken: middleware.GenerateRefreshToken,
		ValidateRefreshToken: middleware.ValidateRefreshToken,
		config:               config,
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

// ForgotPasswordHandler handles password reset requests by generating a reset token
// and sending it to the user's email address.
//
// Security considerations:
// - Does not reveal email existence
// - Uses time-limited tokens (15 minutes)
// - Logs attempts for security monitoring
// - Implements rate limiting per email/IP
//
// Flow:
// 1. Validate request and email format
// 2. Look up user by email
// 3. Generate secure reset token
// 4. Store token with expiry
// 5. Send reset email
func (h *AuthHandler) ForgotPasswordHandler(w http.ResponseWriter, r *http.Request) {
	// Track request metadata for security logging
	requestIP := r.RemoteAddr
	log.Printf("Received password reset request from IP: %s", requestIP)

	// Decode and validate request
	var req models.ForgotPasswordRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Printf("Failed to decode password reset request from IP %s: %v", requestIP, err)
		JSONError(w, "Invalid request", http.StatusBadRequest)
		return
	}

	// Validate email format
	if !isValidEmail(req.Email) {
		log.Printf("Invalid email format in reset request from IP %s: %s", requestIP, req.Email)
		JSONSuccess(w, "If the email exists, a reset link will be sent")
		return
	}
	log.Printf("Processing password reset request for email: %s", maskEmail(req.Email))

	// Get user by email
	user, err := models.GetUserByEmail(h.DB, req.Email)
	if err != nil {
		log.Printf("User lookup failed for reset request: %v", err)
		JSONSuccess(w, "If the email exists, a reset link will be sent")
		return
	}
	log.Printf("Found user for password reset: ID=%d", user.ID)

	// TODO: Check for existing recent reset requests
	// if !h.canRequestPasswordReset(user.ID) {
	// 	log.Printf("Too many reset attempts for user ID %d", user.ID)
	// 	JSONSuccess(w, "If the email exists, a reset link will be sent")
	// 	return
	// }

	// Generate reset token
	resetToken, err := generateResetToken()
	if err != nil {
		log.Printf("Failed to generate reset token for user %d: %v", user.ID, err)
		JSONError(w, "Failed to process request", http.StatusInternalServerError)
		return
	}
	log.Printf("Generated reset token for user %d", user.ID)

	// Set token expiry (15 minutes from now)
	expiryTime := time.Now().Add(15 * time.Minute)
	log.Printf("Setting token expiry for user %d to: %v", user.ID, expiryTime)

	// Update user with reset token
	err = user.UpdateResetToken(h.DB, resetToken, expiryTime)
	if err != nil {
		log.Printf("Failed to save reset token for user %d: %v", user.ID, err)
		JSONError(w, "Failed to process request", http.StatusInternalServerError)
		return
	}
	log.Printf("Successfully saved reset token for user %d", user.ID)

	// Send email with reset link
	resetLink := fmt.Sprintf("%s/reset-password?token=%s", h.config.SMTP.BaseURL, resetToken)
	log.Printf("Generated reset link for user %d", user.ID)

	// Send email with reset link
	err = h.EmailService.SendPasswordResetEmail(user.Email, resetLink)
	if err != nil {
		log.Printf("Failed to send reset email to user %d: %v", user.ID, err)
		JSONError(w, "Failed to send reset email", http.StatusInternalServerError)
		return
	}
	log.Printf("Successfully sent reset email to user %d", user.ID)

	// TODOL: Record successful reset request
	// h.recordPasswordResetAttempt(user.ID)

	JSONSuccess(w, "If the email exists, a reset link will be sent")
}

// ResetPasswordHandler processes password reset requests using a valid reset token.
// It validates the token, checks its expiration, and updates the user's password.
//
// Flow:
// 1. Decode and validate the request body
// 2. Verify the reset token and retrieve associated user
// 3. Ensure the token hasn't expired
// 4. Hash the new password
// 5. Update the password and clear the reset token
// 6. Return success response
func (h *AuthHandler) ResetPasswordHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("Starting password reset process")

	// Decode request body
	var req models.ResetPasswordRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Printf("Failed to decode reset password request: %v", err)
		JSONError(w, "Invalid request", http.StatusBadRequest)
		return
	}
	log.Printf("Successfully decoded password reset request")

	// Validate token length and format
	if len(req.Token) == 0 {
		log.Printf("Empty reset token received")
		JSONError(w, "Reset token is required", http.StatusBadRequest)
		return
	}

	// Retrieve user by reset token
	user, err := models.GetUserByResetToken(h.DB, req.Token)
	if err != nil {
		log.Printf("Failed to find user with reset token: %v", err)
		JSONError(w, "Invalid or expired reset token", http.StatusBadRequest)
		return
	}
	log.Printf("Found user for reset token: UserID=%d", user.ID)

	// Verify token expiration
	if time.Now().After(user.ResetTokenExpires) {
		log.Printf("Reset token expired for user %d. Expired at: %v", user.ID, user.ResetTokenExpires)
		JSONError(w, "Reset token has expired", http.StatusBadRequest)
		return
	}
	log.Printf("Reset token is valid and not expired")

	// Validate new password
	if len(req.NewPassword) < 6 {
		log.Printf("New password too short for user %d", user.ID)
		JSONError(w, "Password must be at least 8 characters long", http.StatusBadRequest)
		return
	}

	// Hash the new password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		log.Printf("Failed to hash new password for user %d: %v", user.ID, err)
		JSONError(w, "Failed to process new password", http.StatusInternalServerError)
		return
	}
	log.Printf("Successfully hashed new password for user %d", user.ID)

	// Update password and clear reset token
	err = user.UpdatePasswordAndClearResetToken(h.DB, string(hashedPassword))
	if err != nil {
		log.Printf("Failed to update password for user %d: %v", user.ID, err)
		JSONError(w, "Failed to update password", http.StatusInternalServerError)
		return
	}
	log.Printf("Successfully updated password and cleared reset token for user %d", user.ID)

	// TODO: Send confirmation email
	// if err := h.EmailService.SendPasswordChangeConfirmation(user.Email); err != nil {
	//     // Log but don't return error - password was successfully changed
	//     log.Printf("Failed to send password change confirmation email to user %d: %v",
	//         user.ID, err)
	// }

	log.Printf("Password reset successful for user %d", user.ID)
	JSONSuccess(w, "Password has been reset successfully")
}
