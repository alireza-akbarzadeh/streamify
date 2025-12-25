package auth

import (
	"net/http"

	"github.com/techies/streamify/internal/handler/users"
	"github.com/techies/streamify/internal/models"
	"github.com/techies/streamify/internal/service"
	"github.com/techies/streamify/internal/utils"
)

// RegisterRequest represents registration payload
type RegisterRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=8"`
	Username string `json:"username" validate:"required,min=3"`
}

// ========================
// Register Handler
// ========================

// @Summary      User registration
// @Description  Register a new user account
// @Tags         Authentication
// @Accept       json
// @Produce      json
// @Param        user  body      RegisterRequest  true  "User registration details"
// @Success      201   {object}  models.UserResponse
// @Failure      400   {object}  utils.ErrorResponse
// @Failure      409   {object}  utils.ErrorResponse
// @Failure      500   {object}  utils.ErrorResponse
// @Router       /api/v1/auth/register [post]
func (h *Handler) Register(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	var req RegisterRequest
	if err := utils.ParseJSON(w, r, &req); err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid request payload", err)
		return
	}

	user, appErr := h.Service.Register(ctx, service.RegisterParams{
		Username: req.Username,
		Email:    req.Email,
		Password: req.Password,
	})

	if appErr != nil {
		utils.RespondWithError(w, appErr.Code, appErr.Message, appErr.Err)
		return
	}

	response := users.MapUserToResponse(user)
	utils.RespondWithJSON(w, http.StatusCreated, models.UserResponse(response))
}
