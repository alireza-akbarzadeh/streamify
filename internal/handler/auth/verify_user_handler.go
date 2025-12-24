package auth

import (
	"database/sql"
	"net/http"
	"time"

	"github.com/techies/streamify/internal/logger"
	"github.com/techies/streamify/internal/models"
	"github.com/techies/streamify/internal/utils"
)

// @Summary      Verify user email token
// @Description  Verifies a user's email using a verification token and marks the user as verified
// @Tags         Authentication
// @Accept       json
// @Produce      json
// @Param        token  query    string  true  "Verification token"
// @Success      303    {string} string "Redirects to frontend with verification status"
// @Failure      400    {object} utils.ErrorResponse
// @Failure      401    {object} utils.ErrorResponse
// @Failure      500    {object} utils.ErrorResponse
// @Router       /auth/verify [get]
func (h *AuthHandler) VerifyToken(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	token := r.URL.Query().Get("token")
	if token == "" {
		logger.Warn(ctx, "VerifyToken: token is required")
		utils.RespondWithError(w, http.StatusBadRequest, "Token is required")
		return
	}

	user, err := h.app.DB.GetUserByVerificationToken(ctx, sql.NullString{String: token, Valid: true})
	if err != nil {
		logger.Warn(ctx, "VerifyToken: invalid token", "token", token)
		utils.RespondWithError(w, http.StatusUnauthorized, "Invalid token")
		return
	}

	if !user.VerificationExpiresAt.Valid || user.VerificationExpiresAt.Time.Before(time.Now()) {
		logger.Warn(ctx, "VerifyToken: token expired", "user_id", user.ID)
		utils.RespondWithError(w, http.StatusUnauthorized, "Token expired")
		return
	}

	err = h.app.DB.VerifyUserByTokenByID(ctx, user.ID)
	if err != nil {
		logger.Error(ctx, "VerifyToken: failed to mark user as verified", err, "user_id", user.ID)
		utils.RespondWithError(w, http.StatusInternalServerError, "Verification failed")
		return
	}

	updatedUser, err := h.app.DB.GetUserById(ctx, user.ID)
	if err != nil {
		logger.Error(ctx, "VerifyToken: failed to fetch updated user", err, "user_id", user.ID)
		utils.RespondWithError(w, http.StatusInternalServerError, "Failed to fetch updated user", err)
		return
	}

	logger.Info(ctx, "User email verified successfully", "user_id", user.ID)
	utils.RespondWithJSON(w, http.StatusOK, map[string]interface{}{
		"message": "Email verified successfully",
		"user":    models.NewUserResponse(&updatedUser),
	})
}
