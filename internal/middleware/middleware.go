package middleware

import (
	"context"
	"errors"
	"net/http"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/techies/streamify/internal/database"
	"github.com/techies/streamify/internal/utils"
)

type contextKey string

const (
	UserIDKey    contextKey = "user_id"
	UserEmailKey contextKey = "user_email"
	UserRoleKey  contextKey = "user_role" // Added typed key for roles
	SessionIDKey contextKey = "session_id"
)

// GetUserID retrieves the user ID from context
func GetUserID(ctx context.Context) string {
	id, _ := ctx.Value(UserIDKey).(string)
	return id
}

// GetSessionID retrieves the session ID from context
func GetSessionID(ctx context.Context) string {
	id, _ := ctx.Value(SessionIDKey).(string)
	return id
}

// GetUserRole retrieves the user role from context
func GetUserRole(ctx context.Context) string {
	role, _ := ctx.Value(UserRoleKey).(string)
	return role
}

// AuthMiddleware validates JWT and injects user info into context
func AuthMiddleware(db *database.Queries, jwtSecret string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authHeader := r.Header.Get("Authorization")

			// 1. Validate Header Format
			if authHeader == "" {
				utils.RespondWithError(w, http.StatusUnauthorized, "Authorization header is required", nil)
				return
			}

			parts := strings.Fields(authHeader)
			if len(parts) != 2 || parts[0] != "Bearer" {
				utils.RespondWithError(w, http.StatusUnauthorized, "Authorization header must be Bearer {token}", nil)
				return
			}

			tokenString := parts[1]

			// 2. Parse and Validate Token
			token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
				if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
					return nil, errors.New("unexpected signing method")
				}
				return []byte(jwtSecret), nil
			})

			if err != nil {
				if errors.Is(err, jwt.ErrTokenExpired) {
					utils.RespondWithError(w, http.StatusUnauthorized, "Token has expired", nil)
					return
				}
				utils.RespondWithError(w, http.StatusUnauthorized, "Invalid token", err)
				return
			}

			// 3. Extract Claims
			claims, ok := token.Claims.(jwt.MapClaims)
			if !ok || !token.Valid {
				utils.RespondWithError(w, http.StatusUnauthorized, "Invalid token claims", nil)
				return
			}

			// 4. Set Context with a single chain
			ctx := r.Context()

			// Subject (UserID) is mandatory
			sub, ok := claims["sub"].(string)
			if !ok || sub == "" {
				utils.RespondWithError(w, http.StatusUnauthorized, "Token missing valid subject", nil)
				return
			}
			ctx = context.WithValue(ctx, UserIDKey, sub)

			// SessionID (sid) is used for server-side invalidation (logout)
			if sid, ok := claims["sid"].(string); ok && sid != "" {
				sessionID, err := uuid.Parse(sid)
				if err == nil {
					// Verify session exists and is not expired
					session, err := db.GetSessionByID(ctx, sessionID)
					if err != nil || time.Now().After(session.ExpiresAt) {
						utils.RespondWithError(w, http.StatusUnauthorized, "Session is invalid or expired", nil)
						return
					}
					ctx = context.WithValue(ctx, SessionIDKey, sid)
				}
			}

			// Optional fields
			if email, ok := claims["email"].(string); ok {
				ctx = context.WithValue(ctx, UserEmailKey, email)
			}
			if role, ok := claims["role"].(string); ok {
				ctx = context.WithValue(ctx, UserRoleKey, role)
			}

			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// AdminOnly middleware restricts access to admin users
func AdminOnly(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		role := GetUserRole(r.Context())
		if role != "admin" {
			utils.RespondWithError(w, http.StatusForbidden, "Admin access required", nil)
			return
		}
		next.ServeHTTP(w, r)
	})
}
