package models

// ErrorResponse represents the structure of error responses
// @Description Error response structure
type ErrorResponse struct {
	// The error message
	// @example "Invalid input data"
	Error string `json:"error"`
}
