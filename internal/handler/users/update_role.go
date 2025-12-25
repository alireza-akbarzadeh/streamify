package users

import (
	"net/http"

	"github.com/google/uuid"
	"github.com/techies/streamify/internal/database"
	"github.com/techies/streamify/internal/logger"
	"github.com/techies/streamify/internal/utils"
)

// UpdateUserRoleRequest defines the payload for updating user role
type UpdateUserRoleRequest struct {
	Role string `json:"role" validate:"required,oneof=customer admin owner"`
}

// UpdateUserRole updates a user's role. Admin only.
// @Summary      Update user role
// @Description  Allows an admin to update the role of any user.
// @Tags         Users
// @Accept       json
// @Produce      json
// @Param        id    path      string                 true  "User ID"
// @Param        role  body      UpdateUserRoleRequest  true  "Role update details"
// @Success      200   {object}  models.UserResponse
// @Failure      400   {object}  utils.ErrorResponse
// @Failure      401   {object}  utils.ErrorResponse
// @Failure      403   {object}  utils.ErrorResponse
// @Failure      500   {object}  utils.ErrorResponse
// @Security     BearerAuth
// @Router       /api/v1/users/{id}/role [put]
func (h *UserHandler) UpdateUserRole(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// 1. Parse user ID from path
	userIDStr := utils.GetParam(r, "id")
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		logger.Warn(ctx, "UpdateUserRole: invalid user ID", "user_id", userIDStr)
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid user ID", err)
		return
	}

	// 2. Parse request body
	var req UpdateUserRoleRequest
	if err := utils.ParseJSON(w, r, &req); err != nil {
		logger.Warn(ctx, "UpdateUserRole: invalid request body", "error", err)
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid JSON format", err)
		return
	}

	// 3. Call service
	appErr := h.Service.UpdateUserRole(ctx, userID, database.UserRole(req.Role))
	if appErr != nil {
		utils.RespondWithError(w, appErr.Code, appErr.Message, appErr.Err)
		return
	}

	utils.RespondWithJSON(w, http.StatusOK, map[string]string{
		"message": "User role updated successfully",
	})
}
