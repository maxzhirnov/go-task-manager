package handlers

import (
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/mail"
	"strings"
)

// JSONError writes a standardized JSON error response to the HTTP response writer.
//
// This utility function ensures consistent error response format across all handlers.
// It sets the appropriate content type, status code, and logs the error message.
//
// Parameters:
//   - w: HTTP ResponseWriter to write the error response to
//   - message: Error message to be included in the response
//   - status: HTTP status code for the response
//
// Response Format:
//
//	{
//	    "error": "error message here"
//	}
//
// Example Usage:
//
//	if err != nil {
//	    JSONError(w, "Invalid input", http.StatusBadRequest)
//	    return
//	}
//
// Note: This function automatically logs the error message for monitoring
// and debugging purposes.
func JSONError(w http.ResponseWriter, message string, status int) {
	// Log the error message
	log.Printf("Error: %s", message)

	// Set response headers
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	// Write JSON response
	if err := json.NewEncoder(w).Encode(map[string]string{"error": message}); err != nil {
		log.Printf("Failed to encode error response: %v", err)
	}
}

// generateResetToken creates a cryptographically secure random token for password reset.
// The token is URL-safe base64 encoded and has 32 bytes of entropy.
//
// Returns:
// - A URL-safe base64 encoded string of 43 characters
//
// Security considerations:
// - Uses crypto/rand for secure random number generation
// - 32 bytes provides 256 bits of entropy
// - URL-safe encoding ensures token can be safely used in URLs
func generateResetToken() (string, error) {
	// Create a buffer for 32 random bytes (256 bits of entropy)
	b := make([]byte, 32)

	// Generate cryptographically secure random bytes
	_, err := rand.Read(b)
	if err != nil {
		return "", fmt.Errorf("failed to generate random bytes: %w", err)
	}

	// Encode using URL-safe base64 (for safe usage in URLs)
	// The resulting string will be 43 characters long
	token := base64.URLEncoding.EncodeToString(b)

	return token, nil
}

// isValidEmail performs basic email format validation
func isValidEmail(email string) bool {
	// Basic email format check
	_, err := mail.ParseAddress(email)
	return err == nil
}

// maskEmail masks part of the email for logging purposes
// Example: j***@example.com
func maskEmail(email string) string {
	parts := strings.Split(email, "@")
	if len(parts) != 2 {
		return "invalid_email"
	}
	if len(parts[0]) <= 1 {
		return email
	}
	return parts[0][:1] + "***@" + parts[1]
}

// validateProfileUpdate performs validation on profile update request data
func validateProfileUpdate(req UpdateProfileRequest) error {
	// Username validation
	if req.Username != "" {
		if len(req.Username) < 3 {
			return fmt.Errorf("username must be at least 3 characters long")
		}
		if len(req.Username) > 30 {
			return fmt.Errorf("username must not exceed 30 characters")
		}
		// Add more username validation rules as needed
	}

	// Password validation
	if req.NewPassword != "" {
		if len(req.NewPassword) < 8 {
			return fmt.Errorf("new password must be at least 8 characters long")
		}
		// Add more password validation rules as needed
	}

	// Current password is required for any changes
	if req.CurrentPassword == "" {
		return fmt.Errorf("current password is required")
	}

	return nil
}
