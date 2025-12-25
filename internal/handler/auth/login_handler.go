package auth

import (
	"net/http"
	"time"

	"github.com/techies/streamify/internal/handler/token"
	"github.com/techies/streamify/internal/logger"
	"github.com/techies/streamify/internal/models"
	"github.com/techies/streamify/internal/service"
	"github.com/techies/streamify/internal/utils"
)

// ========================
// Login Handler
// ========================

// LoginRequest represents login payload
type LoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=8"`
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
func (h *Handler) Login(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	var req LoginRequest
	if err := utils.ParseJSON(w, r, &req); err != nil {
		logger.Warn(ctx, "Malformed login request", "error", err)
		utils.RespondWithError(w, http.StatusBadRequest, "Malformed request", err)
		return
	}

	result, appErr := h.Service.Login(ctx, service.LoginParams{
		Email:     req.Email,
		Password:  req.Password,
		IP:        utils.GetClientIP(r),
		UserAgent: r.UserAgent(),
	})

	if appErr != nil {
		utils.RespondWithError(w, appErr.Code, appErr.Message, appErr.Err)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "refresh_token",
		Value:    result.RefreshToken,
		Path:     "/",
		Expires:  time.Now().Add(token.RefreshTokenTTL),
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteStrictMode,
	})

	logger.Info(ctx, "User logged in successfully", "user_id", result.User.ID, "email", result.User.Email)
	utils.RespondWithJSON(w, http.StatusOK, LoginResponse{
		AccessToken: result.AccessToken,
		User:        models.NewUserResponse(&result.User),
	})
}
