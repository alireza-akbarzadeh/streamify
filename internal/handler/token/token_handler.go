package token

import (
	"crypto/rand"
	"encoding/base64"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/techies/streamify/internal/app"
	"github.com/techies/streamify/internal/database"
	"github.com/techies/streamify/internal/middleware"
	"github.com/techies/streamify/internal/utils"
)

type TokenHandler struct {
	App *app.AppConfig
}

func NewTokenHandler(app *app.AppConfig) *TokenHandler {
	return &TokenHandler{App: app}
}

// Token expiration durations
const (
	AccessTokenTTL  = 15 * time.Minute
	RefreshTokenTTL = 7 * 24 * time.Hour
	RefreshTokenLen = 64
)

// generateSecureToken returns a crypto-secure random string
func GenerateSecureToken(n int) (string, error) {
	b := make([]byte, n)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(b), nil
}

// RefreshToken rotates the user's refresh token and issues a new access token
// @Summary      Refresh access token
// @Description  Validate and rotate refresh token, issue a new access token
// @Tags         Authentication
// @Accept       json
// @Produce      json
// @Success      200  {object}  map[string]string
// @Failure      401  {object}  utils.ErrorResponse
// @Router       /api/v1/auth/refresh [post]
func (h *TokenHandler) RefreshToken(w http.ResponseWriter, r *http.Request) {
	// 1. Get refresh token from cookie
	cookie, err := r.Cookie("refresh_token")
	if err != nil {
		utils.RespondWithError(w, http.StatusUnauthorized, "Refresh token missing", nil)
		return
	}

	// 2. Lookup session in DB
	session, err := h.App.DB.GetSessionByToken(r.Context(), cookie.Value)
	if err != nil {
		utils.RespondWithError(w, http.StatusUnauthorized, "Invalid session", err)
		return
	}

	// 3. Check if session expired
	if time.Now().After(session.ExpiresAt) {
		_ = h.App.DB.DeleteSessionByToken(r.Context(), session.RefreshToken)
		utils.RespondWithError(w, http.StatusUnauthorized, "Session expired", err)
		return
	}

	// 4. Generate new crypto-secure refresh token
	newRefreshToken, err := GenerateSecureToken(RefreshTokenLen)
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "Failed to generate refresh token", err)
		return
	}

	// 5. Rotate session in DB
	if err := h.App.DB.DeleteSessionByToken(r.Context(), session.RefreshToken); err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "Failed to rotate session", err)
		return
	}
	ip := utils.GetClientIP(r)
	IpAddress := utils.ToNullString(&ip)
	userAgent := r.UserAgent()
	UserAgent := utils.ToNullString(&userAgent)

	newSession, err := h.App.DB.CreateSession(r.Context(), database.CreateSessionParams{
		UserID:       session.UserID,
		RefreshToken: newRefreshToken,
		IpAddress:    IpAddress,
		UserAgent:    UserAgent,
		ExpiresAt:    time.Now().Add(RefreshTokenTTL),
	})
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "Failed to create new session", err)
		return
	}
	user, err := h.App.DB.GetUserById(r.Context(), session.UserID)
	if err != nil {
		utils.RespondWithError(w, http.StatusUnauthorized, "User not found", err)
		return
	}
	// 6. Generate new access token
	newAccessToken, err := utils.GenerateToken(
		user.ID,
		newSession.ID,
		AccessTokenTTL,
		h.App.JWTSecret,
		user.Role,
		user.FirstName.String,
		user.LastName.String,
		user.PhoneNumber.String,
		user.Email,
	)

	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "Failed to generate access token", err)
		return
	}

	// 7. Set secure refresh token cookie
	http.SetCookie(w, &http.Cookie{
		Name:     "refresh_token",
		Value:    newRefreshToken,
		Path:     "/",
		Expires:  time.Now().Add(RefreshTokenTTL),
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteStrictMode,
	})

	// 8. Respond with new access token
	utils.RespondWithJSON(w, http.StatusOK, map[string]string{
		"access_token": newAccessToken,
	})
}

// ======================
// Logout Current Session
// ======================

// Logout invalidates the current session
// @Summary      User logout
// @Description  Invalidate the current session and clear refresh token cookie
// @Tags         Authentication
// @Produce      json
// @Success      200  {object}  map[string]string
// @Security     BearerAuth
// @Router       /api/v1/auth/logout [post]
func (h *TokenHandler) Logout(w http.ResponseWriter, r *http.Request) {
	// 1. Try to invalidate session by sid from JWT
	sid := middleware.GetSessionID(r.Context())
	if sid != "" {
		sessionID, err := uuid.Parse(sid)
		if err == nil {
			// Get session to get the refresh token (if we want to be thorough, but we can just delete by ID if we add DeleteSessionByID)
			// Actually, let's just add DeleteSessionByID query.
			_ = h.App.DB.DeleteSessionByID(r.Context(), sessionID)
		}
	}

	// 2. Also try to clear by refresh_token cookie (for web clients)
	cookie, err := r.Cookie("refresh_token")
	if err == nil {
		// Delete the session from database if cookie exists
		_ = h.App.DB.DeleteSessionByToken(r.Context(), cookie.Value)
	}

	// Clear the cookie in the browser
	http.SetCookie(w, &http.Cookie{
		Name:     "refresh_token",
		Value:    "",
		Path:     "/",
		Expires:  time.Unix(0, 0),
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteStrictMode,
	})

	utils.RespondWithJSON(w, http.StatusOK, map[string]string{"message": "Logged out successfully"})
}

// ======================
// Logout All Sessions
// ======================

// LogoutAllDevices invalidates all sessions for the authenticated user
// @Summary      Logout from all devices
// @Description  Revoke all sessions and clear current refresh token
// @Tags         Authentication
// @Produce      json
// @Success      200  {object}  map[string]string
// @Failure      401  {object}  utils.ErrorResponse
// @Failure      500  {object}  utils.ErrorResponse
// @Router       /api/v1/auth/logout-all [post]
func (h *TokenHandler) LogoutAllDevices(w http.ResponseWriter, r *http.Request) {
	userIDStr := middleware.GetUserID(r.Context()) // string
	if userIDStr == "" {
		utils.RespondWithError(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	// Convert string to uuid.UUID
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		utils.RespondWithError(w, http.StatusUnauthorized, "Invalid user ID")
		return
	}

	if err := h.App.DB.DeleteAllUserSessions(r.Context(), userID); err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "Failed to revoke sessions")
		return
	}

	h.clearRefreshCookie(w)
	utils.RespondWithJSON(w, http.StatusOK, map[string]string{
		"message": "Logged out from all devices",
	})
}

func (h *TokenHandler) clearRefreshCookie(w http.ResponseWriter) {
	http.SetCookie(w, &http.Cookie{
		Name:     "refresh_token",
		Value:    "",
		Path:     "/", // matches login/refresh
		MaxAge:   -1,
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteStrictMode,
	})
}
