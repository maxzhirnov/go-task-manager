package handlers

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/maxzhirnov/go-task-manager/internal/middleware"
	"github.com/maxzhirnov/go-task-manager/internal/models"
	"github.com/maxzhirnov/go-task-manager/pkg/database"
)

type AuthHandler struct {
	DB                   database.DB
	GenerateJWT          func(userID int, username string) (string, error)
	GenerateRefreshToken func(username string) (string, error)
	ValidateRefreshToken func(token string) (*middleware.Claims, error)
}

func NewAuthHandler(db database.DB) *AuthHandler {
	return &AuthHandler{
		DB:                   db,
		GenerateJWT:          middleware.GenerateJWT,
		GenerateRefreshToken: middleware.GenerateRefreshToken,
		ValidateRefreshToken: middleware.ValidateRefreshToken,
	}
}

// RegisterHandler registers a new user
func (h *AuthHandler) RegisterHandler(w http.ResponseWriter, r *http.Request) {
	var user models.User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		log.Printf("Error decoding input: %v", err)
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	log.Printf("Registering user with password: %s", user.Password)

	// Hash the password
	if err := user.HashPassword(); err != nil {
		log.Printf("Error hashing password: %v", err)
		http.Error(w, "Failed to hash password", http.StatusInternalServerError)
		return
	}

	log.Printf("Password hashed to: %s", user.Password)

	// Log the hashed password for debugging
	log.Printf("Hashed password for user %s: %s", user.Username, user.Password)

	// Create the user
	if err := user.CreateUser(h.DB); err != nil {
		log.Printf("Error creating user: %v", err)

		// Check for duplicate username error
		if err.Error() == "username already exists" {
			http.Error(w, "Username already exists", http.StatusConflict)
			return
		}

		http.Error(w, "Failed to create user", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{"message": "User registered successfully"})
}

// LoginRequest is a struct for the login request
type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// LoginHandler logs in a user
func (h *AuthHandler) LoginHandler(w http.ResponseWriter, r *http.Request) {
	var req LoginRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Printf("Failed to decode login request: %v", err)
		JSONError(w, "Invalid input", http.StatusBadRequest)
		return
	}

	log.Printf("Login attempt for user: %s with password: %s", req.Username, req.Password)

	// Check if password or username is empty FIRST
	if req.Username == "" || req.Password == "" {
		JSONError(w, "Username and password are required", http.StatusBadRequest)
		return
	}

	// Get the user from the database
	user, err := models.GetUserByUsername(h.DB, req.Username)
	if err != nil {
		log.Printf("Failed to get user by username: %v", err)
		JSONError(w, "Invalid credentials", http.StatusUnauthorized)
		return
	}

	log.Printf("Retrieved user: %s, stored hash: %s", user.Username, user.Password)

	// Check the password
	if err := user.CheckPassword(req.Password); err != nil {
		log.Printf("Password check failed: %v", err)
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
	refreshToken, err := h.GenerateRefreshToken(user.Username)
	if err != nil {
		JSONError(w, "Failed to generate refresh token", http.StatusInternalServerError)
		return
	}

	// Return both tokens to the client
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"access_token":  accessToken,
		"refresh_token": refreshToken,
	})
}

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

	// Return the new access token
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"access_token": accessToken})
}
