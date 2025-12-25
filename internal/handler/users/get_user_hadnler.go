package users

import (
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

	user, appErr := h.Service.GetUser(ctx, userID)
	if appErr != nil {
		utils.RespondWithError(w, appErr.Code, appErr.Message, appErr.Err)
		return
	}

	logger.Info(ctx, "GetUser: user fetched successfully", "user_id", userID)
	response := MapUserToResponse(user)
	utils.RespondWithJSON(w, http.StatusOK, models.UserResponse(response))
}
