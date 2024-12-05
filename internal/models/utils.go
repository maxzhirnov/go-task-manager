package models

import "fmt"

func validateUsername(username string) error {
	if len(username) < 3 {
		return fmt.Errorf("username must be at least 3 characters long")
	}
	if len(username) > 30 {
		return fmt.Errorf("username must not exceed 30 characters")
	}
	// TODO: Add more validation rules as needed
	return nil
}

func validatePassword(password string) error {
	if len(password) < 6 {
		return fmt.Errorf("password must be at least 8 characters long")
	}
	// TODO: Add more validation rules as needed
	return nil
}
