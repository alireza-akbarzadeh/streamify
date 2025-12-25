package auth

import (
	"database/sql"
	"errors"
	"net/http"
	"time"

	"github.com/techies/streamify/internal/database"
	"github.com/techies/streamify/internal/handler/token"
	"github.com/techies/streamify/internal/logger"
	"github.com/techies/streamify/internal/models"
	"github.com/techies/streamify/internal/utils"
	"golang.org/x/crypto/bcrypt"
)

var errInvalidCredentials = errors.New("invalid email or password")

// ========================
// Login Handler
// ========================

// LoginRequest represents login payload
type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=8"`
}

// LoginResponse represents login response
type LoginResponse struct {
	AccessToken string               `json:"access_token"`
	User        *models.UserResponse `json:"user"`
}

// @Summary      User login
// @Description  Authenticate user and return JWT access token + refresh token cookie
// @Tags         Authentication
// @Accept       json
// @Produce      json
// @Param        credentials  body      LoginRequest   true  "Login credentials"
// @Success      200          {object}  LoginResponse
// @Failure      400          {object}  utils.ErrorResponse
// @Failure      401          {object}  utils.ErrorResponse
// @Failure      500          {object}  utils.ErrorResponse
// @Router       /api/v1/auth/login [post]
func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	var req LoginRequest
	if err := utils.ParseJSON(w, r, &req); err != nil {
		logger.Warn(ctx, "Malformed login request", "error", err)
		utils.RespondWithError(w, http.StatusBadRequest, "Malformed request", err)
		return
	}

	user, err := h.App.DB.GetUserByEmail(ctx, req.Email)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			logger.Warn(ctx, "Login failed: user not found", "email", req.Email)
			utils.RespondWithError(w, http.StatusUnauthorized, errInvalidCredentials.Error())
			return
		}
		logger.Error(ctx, "Login failed: DB error", err, "email", req.Email)
		utils.RespondWithError(w, http.StatusInternalServerError, "Failed to fetch user", err)
		return
	}

	if user.IsLocked {
		logger.Warn(ctx, "Login attempt for locked user", "user_id", user.ID, "email", user.Email)
		utils.RespondWithError(w, http.StatusUnauthorized, "User account is locked")
		return
	}

	if !user.IsVerified {
		logger.Warn(ctx, "Login attempt for unverified user", "user_id", user.ID, "email", user.Email)
		utils.RespondWithError(w, http.StatusUnauthorized, "Please verify your email before logging in")
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password)); err != nil {
		logger.Warn(ctx, "Login failed: invalid password", "user_id", user.ID, "email", user.Email)
		utils.RespondWithError(w, http.StatusUnauthorized, errInvalidCredentials.Error())
		return
	}

	refreshToken, err := token.GenerateSecureToken(token.RefreshTokenLen)
	if err != nil {
		logger.Error(ctx, "Failed to generate refresh token", err, "user_id", user.ID)
		utils.RespondWithError(w, http.StatusInternalServerError, "Failed to generate refresh token", err)
		return
	}

	ip := utils.GetClientIP(r)
	IpAddress := utils.ToNullString(&ip)
	userAgent := r.UserAgent()
	UserAgent := utils.ToNullString(&userAgent)

	session, err := h.App.DB.CreateSession(ctx, database.CreateSessionParams{
		UserID:       user.ID,
		RefreshToken: refreshToken,
		IpAddress:    IpAddress,
		UserAgent:    UserAgent,
		ExpiresAt:    time.Now().Add(token.RefreshTokenTTL),
	})
	if err != nil {
		logger.Error(ctx, "Failed to save session", err, "user_id", user.ID)
		utils.RespondWithError(w, http.StatusInternalServerError, "Failed to save session", err)
		return
	}

	accessToken, err := utils.GenerateToken(user.ID, session.ID, token.AccessTokenTTL, h.App.JWTSecret)
	if err != nil {
		logger.Error(ctx, "Failed to generate access token", err, "user_id", user.ID)
		utils.RespondWithError(w, http.StatusInternalServerError, "Failed to generate access token", err)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "refresh_token",
		Value:    refreshToken,
		Path:     "/",
		Expires:  time.Now().Add(token.RefreshTokenTTL),
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteStrictMode,
	})

	logger.Info(ctx, "User logged in successfully", "user_id", user.ID, "email", user.Email)
	utils.RespondWithJSON(w, http.StatusOK, LoginResponse{
		AccessToken: accessToken,
		User:        models.NewUserResponse(&user),
	})
}
