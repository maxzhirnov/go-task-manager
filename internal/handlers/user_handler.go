package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"
	"time"

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

// UpdateProfile handles user profile updates including username and password changes.
// It requires authentication and validates the current password before making any changes.
//
// Security considerations:
// - Requires valid JWT token
// - Validates current password
// - Logs sensitive operations
// - Sanitizes logs of sensitive data
//
// Flow:
// 1. Validate JWT claims
// 2. Parse and validate request
// 3. Verify user exists
// 4. Update profile with provided changes
// 5. Return updated profile data
func (h *UserHandler) UpdateProfile(w http.ResponseWriter, r *http.Request) {

	// Extract and validate JWT claims
	claims := r.Context().Value("claims").(*middleware.Claims)
	log.Printf("[Profile Update] Request received for user ID: %d", claims.UserID)

	/// Parse and validate request body
	var req UpdateProfileRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Printf("[Profile Update] Failed to decode request for user ID %d: %v",
			claims.UserID, err)
		JSONError(w, "Invalid request format", http.StatusBadRequest)
		return
	}

	// Log sanitized request data
	log.Printf("[Profile Update] Request details for user ID %d: "+
		"username_provided=%v, password_change_requested=%v",
		claims.UserID, req.Username != "", req.NewPassword != "")

	// Validate request data
	if err := validateProfileUpdate(req); err != nil {
		log.Printf("[Profile Update] Validation failed for user ID %d: %v",
			claims.UserID, err)
		JSONError(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Retrieve user from database
	user, err := models.GetUserByID(h.db, claims.UserID)
	if err != nil {
		log.Printf("[Profile Update] Failed to retrieve user ID %d: %v",
			claims.UserID, err)
		JSONError(w, "User not found", http.StatusNotFound)
		return
	}
	log.Printf("[Profile Update] Retrieved user ID %d successfully", claims.UserID)

	// Attempt profile update
	if err := user.UpdateProfile(h.db, claims.UserID, req.Username,
		req.NewPassword, req.CurrentPassword); err != nil {

		log.Printf("[Profile Update] Update failed for user ID %d: %v",
			claims.UserID, err)

		// Handle specific error cases
		switch {
		case strings.Contains(err.Error(), "invalid current password"):
			JSONError(w, "Invalid current password", http.StatusUnauthorized)
		case strings.Contains(err.Error(), "username already exists"):
			JSONError(w, "Username already taken", http.StatusConflict)
		default:
			JSONError(w, "Failed to update profile", http.StatusInternalServerError)
		}
		return
	}

	// Log successful update details
	log.Printf("[Profile Update] Successfully updated profile for user ID %d. "+
		"Changes: username=%v, password=%v",
		claims.UserID, req.Username != "", req.NewPassword != "")

	// Prepare and send success response
	response := map[string]interface{}{
		"message":    "Profile updated successfully",
		"username":   req.Username,
		"updated_at": time.Now().Format(time.RFC3339),
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(response); err != nil {
		log.Printf("[Profile Update] Failed to encode response for user ID %d: %v",
			claims.UserID, err)
		JSONError(w, "Internal server error", http.StatusInternalServerError)
		return
	}

}
