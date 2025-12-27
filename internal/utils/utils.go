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
	"os"
	"strconv"
	"strings"
	"time"

	"crypto/rand"
	"encoding/base64"

	"github.com/go-chi/chi/v5"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/techies/streamify/internal/database"
)

type ErrorResponse struct {
	Error string `json:"error"`
}

type AppError struct {
	Code    int    // HTTP Code
	Message string // Client-facing message
	Err     error  // The actual internal error (for logging)
}

func (e *AppError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("%s: %v", e.Message, e.Err)
	}
	return e.Message
}

// ToNullInt32 converts a *int to sql.NullInt32
func ToNullInt32(i *int) sql.NullInt32 {
	if i != nil {
		return sql.NullInt32{Int32: int32(*i), Valid: true}
	}
	return sql.NullInt32{Valid: false}
}

func ToNullUUIDFromString(s string) uuid.NullUUID {
	if s == "" {
		return uuid.NullUUID{}
	}
	id, err := uuid.Parse(s)
	if err != nil {
		return uuid.NullUUID{}
	}
	return uuid.NullUUID{UUID: id, Valid: true}
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
func GenerateToken(
	userID uuid.UUID,
	sessionID uuid.UUID,
	duration time.Duration,
	JWTSecret string,
	role database.UserRole,
	firstName string,
	lastName string,
	phoneNumber string,
	email string,
) (string, error) {
	claims := jwt.MapClaims{
		"sub":          userID.String(),
		"sid":          sessionID.String(),
		"exp":          time.Now().Add(duration).Unix(),
		"iat":          time.Now().Unix(),
		"role":         role,
		"first_name":   firstName,
		"last_name":    lastName,
		"phone_number": phoneNumber,
		"email":        email,
	}
	return jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString([]byte(JWTSecret))
}

// NormalizeEmail trims spaces and converts the email to lowercase
func NormalizeEmail(email string) string {
	return strings.ToLower(strings.TrimSpace(email))
}

// toNullString converts a string to sql.NullString
func ToNullString(s *string) sql.NullString {
	if s == nil {
		// This tells the DB "don't change this value"
		// when used with COALESCE in your SQL.
		return sql.NullString{Valid: false}
	}
	return sql.NullString{
		String: *s,
		Valid:  true,
	}
}

func ToNullUUID(u *uuid.UUID) uuid.NullUUID {
	if u != nil {
		return uuid.NullUUID{UUID: *u, Valid: true}
	}
	return uuid.NullUUID{}
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

func GetParam(r *http.Request, key string) string {
	if r == nil || key == "" {
		return ""
	}
	return chi.URLParam(r, key)
}

func UuidPtr(n uuid.NullUUID) *uuid.UUID {
	if n.Valid {
		return &n.UUID
	}
	return nil
}

func StrPtr(n sql.NullString) *string {
	if n.Valid {
		return &n.String
	}
	return nil
}

func IntPtr(n sql.NullInt32) *int {
	if n.Valid {
		i := int(n.Int32)
		return &i
	}
	return nil
}

func TimePtr(n sql.NullTime) *time.Time {
	if n.Valid {
		return &n.Time
	}
	return nil
}

// ParseInt Helper to keep the main handler logic clean
func ParseInt(value string, fallback int) int {
	if value == "" {
		return fallback
	}
	res, err := strconv.Atoi(value)
	if err != nil || res < 0 {
		return fallback
	}
	return res
}

func ParseUUIDPtr(value *string) (*uuid.UUID, error) {
	if value == nil {
		return nil, nil
	}

	id, err := uuid.Parse(*value)
	if err != nil {
		return nil, err
	}

	return &id, nil
}
