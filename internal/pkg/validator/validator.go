package validator

import (
	"fmt"
	"regexp"
	"strings"
)

// Validator provides validation functions
type Validator struct{}

// NewValidator creates a new validator
func NewValidator() *Validator {
	return &Validator{}
}

// ValidateEmail validates email format
func (v *Validator) ValidateEmail(email string) error {
	if email == "" {
		return fmt.Errorf("email is required")
	}

	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	if !emailRegex.MatchString(email) {
		return fmt.Errorf("invalid email format")
	}

	return nil
}

// ValidateUsername validates username format and length
func (v *Validator) ValidateUsername(username string) error {
	if username == "" {
		return fmt.Errorf("username is required")
	}

	if len(username) < 3 {
		return fmt.Errorf("username must be at least 3 characters")
	}

	if len(username) > 30 {
		return fmt.Errorf("username must be at most 30 characters")
	}

	usernameRegex := regexp.MustCompile(`^[a-zA-Z0-9_-]+$`)
	if !usernameRegex.MatchString(username) {
		return fmt.Errorf("username can only contain letters, numbers, underscores, and hyphens")
	}

	return nil
}

// ValidatePassword validates password strength
func (v *Validator) ValidatePassword(password string) error {
	if password == "" {
		return fmt.Errorf("password is required")
	}

	if len(password) < 8 {
		return fmt.Errorf("password must be at least 8 characters")
	}

	if len(password) > 128 {
		return fmt.Errorf("password must be at most 128 characters")
	}

	var (
		hasLower   = regexp.MustCompile(`[a-z]`).MatchString(password)
		hasUpper   = regexp.MustCompile(`[A-Z]`).MatchString(password)
		hasNumber  = regexp.MustCompile(`\d`).MatchString(password)
		hasSpecial = regexp.MustCompile(`[!@#$%^&*()_+\-=\[\]{};':\"\\|,.<>/?~` + "`" + `]`).MatchString(password)
	)

	if !hasLower {
		return fmt.Errorf("password must contain at least one lowercase letter")
	}

	if !hasUpper {
		return fmt.Errorf("password must contain at least one uppercase letter")
	}

	if !hasNumber {
		return fmt.Errorf("password must contain at least one number")
	}

	if !hasSpecial {
		return fmt.Errorf("password must contain at least one special character")
	}

	return nil
}

// ValidateName validates first name or last name
func (v *Validator) ValidateName(name, fieldName string) error {
	if name == "" {
		return nil // names are optional
	}

	if len(name) > 50 {
		return fmt.Errorf("%s must be at most 50 characters", fieldName)
	}

	name = strings.TrimSpace(name)
	nameRegex := regexp.MustCompile(`^[a-zA-ZÀ-ÿ\s'-]+$`)
	if !nameRegex.MatchString(name) {
		return fmt.Errorf("%s can only contain letters, spaces, apostrophes, and hyphens", fieldName)
	}

	return nil
}