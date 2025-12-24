package users

import (
	"net/http"

	"github.com/google/uuid"
	"github.com/techies/streamify/internal/middleware"
	"github.com/techies/streamify/internal/utils"
)

// DeleteUser handles soft-deleting a user
// @Summary      Soft-delete the current user
// @Description  Marks the authenticated user as deleted (soft delete). The user will no longer be able to log in, but their data is retained.
// @Tags         users
// @Security     BearerAuth
// @Produce      json
// @Success      204  {object}  map[string]string  "User deleted successfully"
// @Failure      400  {object}  map[string]string  "Invalid user ID"
// @Failure      500  {object}  map[string]string  "Failed to delete user"
// @Router       /users/me [delete]
func (h *UserHandler) DeleteUser(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// Get user ID from context (middleware)
	userIDStr := middleware.GetUserID(ctx)
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid user ID", err)
		return
	}

	// Call the repository to soft-delete the user
	if err := h.app.DB.SoftDeleteUser(ctx, userID); err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "Failed to delete user", err)
		return
	}

	message := "Your account has been soft-deleted. It will be permanently removed after 40 days."

	utils.RespondWithJSON(w, http.StatusOK, map[string]string{"message": message})
}

// PermanentlyDeleteOldSoftDeletedUsers permanently deletes users who were soft-deleted and meet the criteria
//
// @Summary      Permanently delete old soft-deleted users
// @Description  Removes from the database all users who were previously soft-deleted and are eligible for permanent deletion (e.g., deleted for a certain period).
// @Tags         users
// @Security     BearerAuth
// @Produce      json
// @Success      200  {object}  map[string]string  "Old soft-deleted users permanently deleted successfully"
// @Failure      500  {object}  map[string]string  "Failed to permanently delete old soft-deleted users"
// @Router       /admin/users/purge-soft-deleted [delete]
func (h *UserHandler) PermanentlyDeleteOldSoftDeletedUsers(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// Call the repository to permanently delete old soft-deleted users
	if err := h.app.DB.PermanentlyDeleteOldSoftDeletedUsers(ctx); err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "Failed to permanently delete old soft-deleted users", err)
		return
	}

	// Return 200 OK
	utils.RespondWithJSON(w, http.StatusOK, map[string]string{"message": "Old soft-deleted users permanently deleted successfully"})
}
