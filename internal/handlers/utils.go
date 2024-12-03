package handlers

import (
	"encoding/json"
	"log"
	"net/http"
)

// JSONError sends a JSON-formatted error response
func JSONError(w http.ResponseWriter, message string, status int) {
	log.Printf("Error: %s", message)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(map[string]string{"error": message})
}
