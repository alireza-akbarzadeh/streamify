package users

import (
	"net/http"

	"github.com/google/uuid"
	"github.com/techies/streamify/internal/database"
	"github.com/techies/streamify/internal/logger"
	"github.com/techies/streamify/internal/middleware"
	"github.com/techies/streamify/internal/models"
	"github.com/techies/streamify/internal/service"
	"github.com/techies/streamify/internal/utils"
)

type UpdateProfileRequest struct {
	FirstName   *string `json:"first_name"`
	LastName    *string `json:"last_name"`
	Bio         *string `json:"bio"`
	AvatarUrl   *string `json:"avatar_url"`
	PhoneNumber *string `json:"phone_number"`
}

// UpdateProfile updates the authenticated user's profile information.
// @Summary      Update user profile
// @Description  Update the profile details of the currently authenticated user.
// @Tags         Users
// @Accept       json
// @Produce      json
// @Param        profile  body      UpdateProfileRequest  true  "Profile update details"
// @Success      200      {object}  models.UserResponse
// @Failure      400      {object}  utils.ErrorResponse
// @Failure      401      {object}  utils.ErrorResponse
// @Failure      500      {object}  utils.ErrorResponse
// @Security     BearerAuth
// @Router       /api/v1/users/{id} [put]
func (h *UserHandler) UpdateProfile(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// 1. Extract UserID safely
	userIDStr := middleware.GetUserID(ctx)
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		logger.Warn(ctx, "UpdateProfile: invalid user ID in context", "user_id", userIDStr)
		utils.RespondWithError(w, http.StatusUnauthorized, "Authentication required", nil)
		return
	}

	// 2. Parse Request
	var req UpdateProfileRequest
	if err := utils.ParseJSON(w, r, &req); err != nil {
		logger.Warn(ctx, "UpdateProfile: malformed request", "error", err)
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid JSON format", err)
		return
	}

	// 3. Service Call
	user, appErr := h.Service.UpdateProfile(ctx, service.UpdateProfileParams{
		UserID:      userID,
		FirstName:   req.FirstName,
		LastName:    req.LastName,
		Bio:         req.Bio,
		AvatarUrl:   req.AvatarUrl,
		PhoneNumber: req.PhoneNumber,
	})

	if appErr != nil {
		utils.RespondWithError(w, appErr.Code, appErr.Message, appErr.Err)
		return
	}

	utils.RespondWithJSON(w, http.StatusOK, models.UserResponse(MapUserToResponse(user)))
}

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
// @Success      200   {object}  map[string]string
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
