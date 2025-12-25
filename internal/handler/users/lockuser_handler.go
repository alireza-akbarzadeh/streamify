package users

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/techies/streamify/internal/logger"
	"github.com/techies/streamify/internal/utils"
)

// @Summary      Lock user account
// @Description  Locks a user account by setting is_locked to true
// @Tags         Users
// @Accept       json
// @Produce      json
// @Param        id   path      string  true  "User ID (UUID)"
// @Success      200  {object}  map[string]string
// @Failure      400  {object}  utils.ErrorResponse
// @Failure      500  {object}  utils.ErrorResponse
// @Security     BearerAuth
// @Router       /api/v1/users/{id}/lock [post]
func (h *UserHandler) LockUser(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	id := chi.URLParam(r, "id")
	uid, err := uuid.Parse(id)
	if err != nil {
		logger.Warn(ctx, "LockUser: invalid user ID", "id", id)
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid user ID")
		return
	}
	err = h.App.DB.LockUser(ctx, uid)
	if err != nil {
		logger.Error(ctx, "LockUser: failed to lock user", err, "user_id", uid)
		utils.RespondWithError(w, http.StatusInternalServerError, "Failed to lock user", err)
		return
	}
	// Invalidate all sessions for this user
	if err := h.App.DB.DeleteAllUserSessions(ctx, uid); err != nil {
		logger.Error(ctx, "LockUser: failed to delete user sessions", err, "user_id", uid)
		// You may want to continue, or return an error depending on your policy
	}
	logger.Info(ctx, "User locked and sessions invalidated", "user_id", uid)
	utils.RespondWithJSON(w, http.StatusOK, map[string]string{"message": "User locked and logged out"})
}

// @Summary      Unlock user account
// @Description  Unlocks a user account by setting is_locked to false
// @Tags         Users
// @Accept       json
// @Produce      json
// @Param        id   path      string  true  "User ID (UUID)"
// @Success      200  {object}  map[string]string
// @Failure      400  {object}  utils.ErrorResponse
// @Failure      500  {object}  utils.ErrorResponse
// @Security     BearerAuth
// @Router       /api/v1/users/{id}/unlock [post]
func (h *UserHandler) UnLockUser(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	id := chi.URLParam(r, "id")
	uid, err := uuid.Parse(id)
	if err != nil {
		logger.Warn(ctx, "UnLockUser: invalid user ID", "id", id)
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid user ID")
		return
	}
	err = h.App.DB.UnlockUser(ctx, uid)
	if err != nil {
		logger.Error(ctx, "UnLockUser: failed to unlock user", err, "user_id", uid)
		utils.RespondWithError(w, http.StatusInternalServerError, "Failed to unlock user", err)
		return
	}
	logger.Info(ctx, "User unlocked successfully", "user_id", uid)
	utils.RespondWithJSON(w, http.StatusOK, map[string]string{"message": "User unlocked"})

}
