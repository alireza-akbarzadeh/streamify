package users

import (
	"database/sql"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/techies/streamify/internal/logger"
	"github.com/techies/streamify/internal/models"
	"github.com/techies/streamify/internal/utils"
)

// GetUser returns a user by ID.
// @Summary      Get user by ID
// @Description  Get a single user by their unique ID.
// @Tags         Users
// @Accept       json
// @Produce      json
// @Param        id    path      string  true   "User ID (UUID)"
// @Success      200   {object}  models.UserResponse
// @Failure      400   {object}  utils.ErrorResponse
// @Failure      404   {object}  utils.ErrorResponse
// @Failure      500   {object}  utils.ErrorResponse
// @Security     BearerAuth
// @Router       /api/v1/users/{id} [get]
func (h *UserHandler) GetUser(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	idParam := chi.URLParam(r, "id")
	userID, err := uuid.Parse(idParam)
	if err != nil {
		logger.Warn(ctx, "GetUser: invalid user ID format", "user_id", idParam)
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid user ID format", nil)
		return
	}

	dbUser, err := h.App.DB.GetUserById(ctx, userID)
	if err != nil {
		if err == sql.ErrNoRows {
			logger.Warn(ctx, "GetUser: user not found", "user_id", userID)
			utils.RespondWithError(w, http.StatusNotFound, "User not found", nil)
			return
		}
		logger.Error(ctx, "GetUser: database error", err, "user_id", userID)
		utils.RespondWithError(w, http.StatusInternalServerError, "Database error", err)
		return
	}

	logger.Info(ctx, "GetUser: user fetched successfully", "user_id", userID)
	response := MapUserToResponse(dbUser)
	utils.RespondWithJSON(w, http.StatusOK, models.UserResponse(response))
}
