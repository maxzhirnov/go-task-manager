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

type AuthHandler struct {
	DB                   database.DB
	EmailService         email.EmailSender
	GenerateJWT          func(userID int, username string) (string, error)
	GenerateRefreshToken func(userID int, username string) (string, error)
	ValidateRefreshToken func(token string) (*middleware.Claims, error)
}

func NewAuthHandler(db database.DB, emailService email.EmailSender) *AuthHandler {
	return &AuthHandler{
		DB:                   db,
		EmailService:         emailService,
		GenerateJWT:          middleware.GenerateJWT,
		GenerateRefreshToken: middleware.GenerateRefreshToken,
		ValidateRefreshToken: middleware.ValidateRefreshToken,
	}
}

type RegisterRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// @Summary Register new user
// @Description Register a new user in the system
// @Tags auth
// @Accept json
// @Produce json
// @Param user body models.User true "User registration details"
// @Success 201 {object} map[string]string "User registered successfully"
// @Failure 400 {object} models.ErrorResponse "Invalid input"
// @Failure 409 {object} models.ErrorResponse "Username already exists"
// @Failure 500 {object} models.ErrorResponse "Internal server error"
// @Router /register [post]
func (h *AuthHandler) RegisterHandler(w http.ResponseWriter, r *http.Request) {
	var req RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		JSONError(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Validate email and password
	if req.Email == "" || req.Password == "" {
		JSONError(w, "Email and password are required", http.StatusBadRequest)
		return
	}

	user := &models.User{
		Email:    req.Email,
		Password: req.Password,
	}

	if err := user.HashPassword(); err != nil {
		JSONError(w, "Error hashing password", http.StatusInternalServerError)
		return
	}

	if err := user.CreateUser(h.DB); err != nil {
		if strings.Contains(err.Error(), "email already exists") {
			JSONError(w, "Email already exists", http.StatusConflict)
			return
		}
		JSONError(w, "Error creating user", http.StatusInternalServerError)
		return
	}

	token, err := models.GetVerificationTokenForUser(h.DB, user.ID)
	if err != nil {
		log.Printf("Failed to get verification token: %v", err)
	} else {
		if err := h.EmailService.SendVerificationEmail(user.Email, user.Username, token); err != nil {
			log.Printf("Failed to send verification email: %v", err)
		}
	}

	log.Printf("User registered successfully")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{
		"message": "User registered successfully",
	})
}

// LoginRequest represents the login request payload
// @Description Login request structure
type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// @Summary Login user
// @Description Authenticate user and return access and refresh tokens
// @Tags auth
// @Accept json
// @Produce json
// @Param credentials body LoginRequest true "User credentials"
// @Success 200 {object} map[string]string "Returns access_token and refresh_token"
// @Failure 400 {object} models.ErrorResponse "Invalid input"
// @Failure 401 {object} models.ErrorResponse "Invalid credentials"
// @Failure 500 {object} models.ErrorResponse "Internal server error"
// @Router /login [post]
func (h *AuthHandler) LoginHandler(w http.ResponseWriter, r *http.Request) {
	var req LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		JSONError(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Check if email or password is empty
	if req.Email == "" || req.Password == "" {
		JSONError(w, "Email and password are required", http.StatusBadRequest)
		return
	}

	// Get the user from the database
	user, err := models.GetUserByEmail(h.DB, req.Email)
	if err != nil {
		if err == sql.ErrNoRows {
			JSONError(w, "Invalid credentials", http.StatusUnauthorized)
			return
		}
		JSONError(w, "Server error", http.StatusInternalServerError)
		return
	}

	// Check if email is verified
	if !user.IsVerified {
		JSONError(w, "Please verify your email before logging in", http.StatusForbidden)
		return
	}

	// Check the password
	if err := user.CheckPassword(req.Password); err != nil {
		JSONError(w, "Invalid credentials", http.StatusUnauthorized)
		return
	}

	// Generate tokens
	accessToken, err := h.GenerateJWT(user.ID, user.Username)
	if err != nil {
		JSONError(w, "Failed to generate access token", http.StatusInternalServerError)
		return
	}

	refreshToken, err := h.GenerateRefreshToken(user.ID, user.Username)
	if err != nil {
		JSONError(w, "Failed to generate refresh token", http.StatusInternalServerError)
		return
	}

	// Return tokens
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"access_token":  accessToken,
		"refresh_token": refreshToken,
	})
}

// @Summary Refresh access token
// @Description Get new access token using refresh token
// @Tags auth
// @Accept json
// @Produce json
// @Param refresh_token body map[string]string true "Refresh token"
// @Success 200 {object} map[string]string "Returns new access_token"
// @Failure 400 {object} models.ErrorResponse "Invalid input"
// @Failure 401 {object} models.ErrorResponse "Invalid refresh token"
// @Failure 500 {object} models.ErrorResponse "Internal server error"
// @Router /refresh [post]
func (h *AuthHandler) RefreshTokenHandler(w http.ResponseWriter, r *http.Request) {
	var req struct {
		RefreshToken string `json:"refresh_token"`
	}

	// Parse the request body
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		JSONError(w, "Invalid input", http.StatusBadRequest)
		return
	}

	// Validate the refresh token
	claims, err := h.ValidateRefreshToken(req.RefreshToken)
	if err != nil {
		JSONError(w, "Invalid refresh token", http.StatusUnauthorized)
		return
	}

	// Generate a new access token
	accessToken, err := h.GenerateJWT(claims.UserID, claims.Username)
	if err != nil {
		JSONError(w, "Failed to generate access token", http.StatusInternalServerError)
		return
	}

	log.Printf("New access token: %s", accessToken)
	// Return the new access token
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"access_token": accessToken})
}

func (h *AuthHandler) VerifyEmailHandler(w http.ResponseWriter, r *http.Request) {
	token := r.URL.Query().Get("token")
	if token == "" {
		JSONError(w, "Verification token is required", http.StatusBadRequest)
		return
	}

	err := models.VerifyEmail(h.DB, token)
	if err != nil {
		if strings.Contains(err.Error(), "invalid or expired") {
			JSONError(w, "Invalid or expired verification token", http.StatusBadRequest)
			return
		}
		JSONError(w, "Error verifying email", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{
		"message": "Email verified successfully",
	})
}

// Добавим метод для повторной отправки верификационного письма
func (h *AuthHandler) ResendVerificationHandler(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Email string `json:"email"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		JSONError(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	user, err := models.GetUserByEmail(h.DB, req.Email)
	if err != nil {
		JSONError(w, "User not found", http.StatusNotFound)
		return
	}

	if user.IsVerified {
		JSONError(w, "Email is already verified", http.StatusBadRequest)
		return
	}

	token, err := models.CreateVerificationToken(h.DB, user.ID)
	if err != nil {
		JSONError(w, "Error generating verification token", http.StatusInternalServerError)
		return
	}

	if err := h.EmailService.SendVerificationEmail(user.Email, user.Username, token.Token); err != nil {
		log.Printf("Failed to send verification email: %v", err)
		JSONError(w, "Error sending verification email", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{
		"message": "Verification email sent successfully",
	})
}
