package users

import (
	"net/http"

	"github.com/google/uuid"
	"github.com/techies/streamify/internal/database"
	"github.com/techies/streamify/internal/logger"
	"github.com/techies/streamify/internal/middleware"
	"github.com/techies/streamify/internal/models"
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

	// 3. Database Execution
	// Using the *string version of ToNullString
	err = h.App.DB.UpdateUserProfile(ctx, database.UpdateUserProfileParams{
		ID:          userID,
		FirstName:   utils.ToNullString(req.FirstName),
		LastName:    utils.ToNullString(req.LastName),
		Bio:         utils.ToNullString(req.Bio),
		AvatarUrl:   utils.ToNullString(req.AvatarUrl),
		PhoneNumber: utils.ToNullString(req.PhoneNumber),
	})

	if err != nil {
		logger.Error(ctx, "UpdateProfile: failed to update profile", err, "user_id", userID)
		utils.RespondWithError(w, http.StatusInternalServerError, "Failed to update profile", err)
		return
	}

	// 4. Return Updated State
	// Re-fetching ensures the client has the absolute 'truth' from the DB
	updatedUser, err := h.App.DB.GetUserById(ctx, userID)
	if err != nil {
		logger.Error(ctx, "UpdateProfile: failed to fetch updated user", err, "user_id", userID)
		utils.RespondWithJSON(w, http.StatusOK, map[string]string{"message": "Profile updated"})
		return
	}

	utils.RespondWithJSON(w, http.StatusOK, models.UserResponse(MapUserToResponse(updatedUser)))
}
