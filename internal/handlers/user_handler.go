package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"

	"github.com/maxzhirnov/go-task-manager/internal/middleware"
	"github.com/maxzhirnov/go-task-manager/internal/models"
	"github.com/maxzhirnov/go-task-manager/pkg/database"
)

type UserHandler struct {
	db database.DB
}

type UpdateProfileRequest struct {
	Email           string `json:"email"`
	Username        string `json:"username"`
	NewPassword     string `json:"new_password,omitempty"`
	CurrentPassword string `json:"current_password"`
}

func NewUserHandler(db database.DB) *UserHandler {
	return &UserHandler{db: db}
}

func (h *UserHandler) UpdateProfile(w http.ResponseWriter, r *http.Request) {
	claims := r.Context().Value("claims").(*middleware.Claims)
	log.Printf("Received profile update request for user %s with id %d", claims.Username, claims.UserID)

	// Parse request body
	var req UpdateProfileRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Printf("Failed to decode request body: %v", err)
		JSONError(w, "Invalid request format", http.StatusBadRequest)
		return
	}
	log.Printf("Request data: %+v", req)

	// Get user from database
	user, err := models.GetUserByID(h.db, claims.UserID)
	if err != nil {
		log.Printf("Failed to get user: %v", err)
		JSONError(w, "User not found", http.StatusNotFound)
		return
	}
	log.Printf("User: %+v", user)

	// Update profile
	if err := user.UpdateProfile(h.db, claims.UserID, req.Username, req.NewPassword, req.CurrentPassword); err != nil {
		log.Printf("Failed to update profile: %v", err)
		switch {
		case strings.Contains(err.Error(), "invalid current password"):
			JSONError(w, "Invalid current password", http.StatusUnauthorized)
		default:
			JSONError(w, "Failed to update profile", http.StatusInternalServerError)
		}
		return
	}

	// Return success response
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"message":  "Profile updated successfully",
		"username": req.Username,
	})

}
