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

	if err := h.EmailService.SendWelcomeEmail(user.Email, user.Username); err != nil {
		// Log the error but don't fail the registration
		log.Printf("Failed to send welcome email: %v", err)
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

	log.Printf("Login attempt for user: %s with password: %s", req.Email, req.Password)

	// Check if password or username is empty FIRST
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

	log.Printf("Retrieved user: %s, stored hash: %s", user.Username, user.Password)

	// Check the password
	if err := user.CheckPassword(req.Password); err != nil {
		JSONError(w, "Invalid credentials", http.StatusUnauthorized)
		return
	}

	// Generate the access token
	accessToken, err := h.GenerateJWT(user.ID, user.Username)
	if err != nil {
		JSONError(w, "Failed to generate access token", http.StatusInternalServerError)
		return
	}

	// Generate the refresh token
	refreshToken, err := h.GenerateRefreshToken(user.ID, user.Username)
	if err != nil {
		JSONError(w, "Failed to generate refresh token", http.StatusInternalServerError)
		return
	}

	log.Printf("New refresh token: %s", refreshToken)
	// Return both tokens to the client
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
