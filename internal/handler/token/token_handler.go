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
	app *app.AppConfig
}

func NewTokenHandler(app *app.AppConfig) *TokenHandler {
	return &TokenHandler{app: app}
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
// @Router       /auth/refresh [post]  // Remove /api/v1 from here

func (h *TokenHandler) RefreshToken(w http.ResponseWriter, r *http.Request) {
	// 1. Get refresh token from cookie
	cookie, err := r.Cookie("refresh_token")
	if err != nil {
		utils.RespondWithError(w, http.StatusUnauthorized, "Refresh token missing", err)
		return
	}

	// 2. Lookup session in DB
	session, err := h.app.DB.GetSessionByToken(r.Context(), cookie.Value)
	if err != nil {
		utils.RespondWithError(w, http.StatusUnauthorized, "Invalid session", err)
		return
	}

	// 3. Check if session expired
	if time.Now().After(session.ExpiresAt) {
		_ = h.app.DB.DeleteSessionByToken(r.Context(), session.RefreshToken)
		utils.RespondWithError(w, http.StatusUnauthorized, "Session expired", err)
		return
	}

	// 4. Generate new access token
	newAccessToken, err := utils.GenerateToken(session.UserID, AccessTokenTTL, h.app.JWTSecret)
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "Failed to generate access token", err)
		return
	}

	// 5. Generate new crypto-secure refresh token
	newRefreshToken, err := GenerateSecureToken(RefreshTokenLen)
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "Failed to generate refresh token", err)
		return
	}

	// 6. Rotate session in DB
	if err := h.app.DB.DeleteSessionByToken(r.Context(), session.RefreshToken); err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "Failed to rotate session", err)
		return
	}
	ip := utils.GetClientIP(r)
	IpAddress := utils.ToNullString(&ip)
	userAgent := r.UserAgent()
	UserAgent := utils.ToNullString(&userAgent)

	_, err = h.app.DB.CreateSession(r.Context(), database.CreateSessionParams{
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
// @Router       /auth/logout [post]
func (h *TokenHandler) Logout(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("refresh_token")
	if err == nil {
		// Delete the session from database if cookie exists
		_ = h.app.DB.DeleteSessionByToken(r.Context(), cookie.Value)
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

	if err := h.app.DB.DeleteAllUserSessions(r.Context(), userID); err != nil {
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
