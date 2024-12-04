package handlers

import (
	"encoding/json"
	"log"
	"net/http"
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
