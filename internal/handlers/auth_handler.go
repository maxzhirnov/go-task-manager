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
	DB database.DB
}

func NewAuthHandler(db database.DB) *AuthHandler {
	return &AuthHandler{DB: db}
}

// RegisterHandler registers a new user
func (h *AuthHandler) RegisterHandler(w http.ResponseWriter, r *http.Request) {
	var user models.User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	// Hash the password
	if err := user.HashPassword(); err != nil {
		http.Error(w, "Failed to hash password", http.StatusInternalServerError)
		return
	}

	// Create the user
	if err := user.CreateUser(h.DB); err != nil {
		// Log the error to help with debugging
		log.Printf("Error creating user: %v", err)

		// Return a more specific error message if it's a duplicate username
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

// LoginHandler logs in a user
func (h *AuthHandler) LoginHandler(w http.ResponseWriter, r *http.Request) {
	var req models.User
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		JSONError(w, "Invalid input", http.StatusBadRequest)
		return
	}

	// Get the user from the database
	user, err := models.GetUserByUsername(h.DB, req.Username)
	if err != nil {
		JSONError(w, "Invalid credentials", http.StatusUnauthorized)
		return
	}

	// Check the password
	if err := user.CheckPassword(req.Password); err != nil {
		JSONError(w, "Invalid credentials", http.StatusUnauthorized)
		return
	}

	// Generate a JWT access token
	accessToken, err := middleware.GenerateJWT(user.Username)
	if err != nil {
		JSONError(w, "Failed to generate access token", http.StatusInternalServerError)
		return
	}

	// Generate a refresh token
	refreshToken, err := middleware.GenerateRefreshToken(user.Username)
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
	claims, err := middleware.ValidateRefreshToken(req.RefreshToken)
	if err != nil {
		JSONError(w, "Invalid refresh token", http.StatusUnauthorized)
		return
	}

	// Generate a new access token
	accessToken, err := middleware.GenerateJWT(claims.Username)
	if err != nil {
		JSONError(w, "Failed to generate access token", http.StatusInternalServerError)
		return
	}

	// Return the new access token
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"access_token": accessToken})
}
