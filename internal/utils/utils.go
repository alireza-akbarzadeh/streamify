package utils

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"net/mail"
	"os"
	"regexp"
	"strings"
	"time"
	"unicode"

	"crypto/rand"
	"encoding/base64"

	"github.com/go-chi/chi/v5"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type ErrorResponse struct {
	Error string `json:"error"`
}

var (
	// ErrInvalidEmail is returned when email format is incorrect
	ErrInvalidEmail = errors.New("invalid email address format")
	// ErrPasswordTooShort is returned when password is less than 8 characters
	ErrPasswordTooShort = errors.New("password must be at least 8 characters long")
	// ErrPasswordTooWeak is returned when password lacks complexity
	ErrPasswordTooWeak = errors.New("password must contain at least one uppercase letter, one lowercase letter, and one number")
)

// ValidateEmail checks if the string is a valid email address
func ValidateEmail(email string) error {
	_, err := mail.ParseAddress(email)
	if err != nil {
		return ErrInvalidEmail
	}

	// Additional regex for stricter checking of the domain part
	emailRegex := regexp.MustCompile(`^[a-z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2,4}$`)
	if !emailRegex.MatchString(NormalizeEmail(email)) {
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

// GenerateRandomString returns a securely generated random string of the given length.
// Returns an error if crypto/rand fails.
func GenerateRandomString(n int) (string, error) {
	b := make([]byte, n)
	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}
	// base64 encode to make it URL-safe
	return base64.URLEncoding.EncodeToString(b)[:n], nil
}

func RespondWithError(w http.ResponseWriter, code int, msg string, err ...error) {
	fullMsg := msg
	if err != nil {
		fullMsg = fmt.Sprintf("%s: %v", msg, err)
	}

	if code >= 500 {
		log.Printf("Internal Server Error: %s", fullMsg)
	}

	RespondWithJSON(w, code, ErrorResponse{Error: fullMsg})
}

func RespondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	data, err := json.Marshal(payload)
	if err != nil {
		log.Printf("Failed to marshal JSON: %v", err)
		http.Error(w, "Failed to marshal JSON", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	if _, err := w.Write(data); err != nil {
		log.Printf("Failed to write response: %v", err)
		http.Error(w, "Failed to write response", http.StatusInternalServerError)
		return
	}
}

// readUUIDParam extracts a UUID from a path or query parameter.
// Returns the UUID and an error message (empty if no error).
func ReadUUIDParam(r *http.Request, param string) (uuid.UUID, string) {
	idStr := chi.URLParam(r, param)
	if idStr == "" {
		idStr = r.URL.Query().Get(param)
	}
	if idStr == "" {
		return uuid.Nil, "you need to provide " + param
	}
	id, err := uuid.Parse(idStr)
	if err != nil {
		return uuid.Nil, "invalid " + param + " format"
	}
	return id, ""
}

func GetEnvString(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}

// ParseJSON decodes the request body into the provided data structure.
// It limits the body size to 1MB to prevent memory exhaustion attacks.
// ... existing code ...
func ParseJSON(w http.ResponseWriter, r *http.Request, data interface{}) error {
	// 1. Limit the size of the request body (1MB)
	maxBytes := int64(1024 * 1024)
	r.Body = http.MaxBytesReader(w, r.Body, maxBytes)

	// 2. Decode the JSON
	decoder := json.NewDecoder(r.Body)
	// This helps prevent some types of JSON attacks
	decoder.DisallowUnknownFields()

	err := decoder.Decode(data)
	if err != nil {
		if errors.Is(err, io.EOF) {
			return errors.New("request body cannot be empty")
		}
		return err
	}

	return nil
}

// WriteJSON is a helper to send JSON responses back to the client
func WriteJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(data); err != nil {
		log.Printf("Failed to encode JSON response: %v", err)
		http.Error(w, "Failed to encode JSON response", http.StatusInternalServerError)
	}
}

// Internal helper for token generation
func GenerateToken(userID uuid.UUID, duration time.Duration, JWTSecret string) (string, error) {
	claims := jwt.MapClaims{
		"sub": userID.String(),
		"exp": time.Now().Add(duration).Unix(),
		"iat": time.Now().Unix(),
	}
	return jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString([]byte(JWTSecret))
}

// NormalizeEmail trims spaces and converts the email to lowercase
func NormalizeEmail(email string) string {
	return strings.ToLower(strings.TrimSpace(email))
}

// toNullString converts a string to sql.NullString
func ToNullString(s string) sql.NullString {
	if s == "" {
		return sql.NullString{Valid: false}
	}
	return sql.NullString{String: s, Valid: true}
}

// normalizeIP strips port if present
func NormalizeIP(addr string) string {
	host, _, err := net.SplitHostPort(addr)
	if err != nil {
		return addr
	}
	return host
}

func GetClientIP(r *http.Request) string {
	// Check X-Forwarded-For header (may contain multiple IPs)
	if xff := r.Header.Get("X-Forwarded-For"); xff != "" {
		parts := strings.Split(xff, ",")
		return strings.TrimSpace(parts[0])
	}
	// Fallback to RemoteAddr
	ip, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		return r.RemoteAddr
	}
	return ip
}
