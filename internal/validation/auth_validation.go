package validation

import (
	"errors"
	"net/mail"
	"regexp"
	"strings"
	"unicode"

	"github.com/techies/streamify/internal/utils"
)

var (
	// ErrInvalidEmail is returned when email format is incorrect
	ErrInvalidEmail = errors.New("invalid email address format")
	// ErrPasswordTooShort is returned when password is less than 8 characters
	ErrPasswordTooShort = errors.New("password must be at least 8 characters long")
	// ErrPasswordTooWeak is returned when password lacks complexity
	ErrPasswordTooWeak = errors.New("password must contain at least one uppercase letter, one lowercase letter, and one number")

	// emailRegex is a stricter regex for email validation, compiled once
	emailRegex = regexp.MustCompile(`^[a-z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2,}$`)
)

// ValidateEmail checks if the string is a valid email address
func ValidateEmail(email string) error {
	_, err := mail.ParseAddress(email)
	if err != nil {
		return ErrInvalidEmail
	}

	// Use the package-level compiled regex for stricter checking of the domain part
	if !emailRegex.MatchString(utils.NormalizeEmail(email)) {
		return ErrInvalidEmail
	}

	return nil
}

// ValidatePassword checks for length and complexity
func ValidatePassword(password string) error {
	if len(password) < 8 {
		return ErrPasswordTooShort
	}

	var (
		hasUpper bool
		hasLower bool
		hasDigit bool
	)

	for _, char := range password {
		switch {
		case unicode.IsUpper(char):
			hasUpper = true
		case unicode.IsLower(char):
			hasLower = true
		case unicode.IsDigit(char):
			hasDigit = true
		}
	}

	if !hasUpper || !hasLower || !hasDigit {
		return ErrPasswordTooWeak
	}

	return nil
}

func ValidateUsername(username string) error {
	username = strings.TrimSpace(username)
	if len(username) < 3 {
		return errors.New("username must be at least 3 characters")
	}
	if len(username) > 30 {
		return errors.New("username cannot exceed 30 characters")
	}
	return nil
}

// ValidateRegisterInput coordinates the validation of all registration fields
func ValidateRegisterInput(username, email, password string) error {
	if err := ValidateUsername(username); err != nil {
		return err
	}
	if err := ValidateEmail(email); err != nil {
		return err
	}
	return ValidatePassword(password)
}
