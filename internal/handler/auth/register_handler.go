package auth

import (
	"database/sql"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/techies/streamify/internal/database"
	"github.com/techies/streamify/internal/handler/token"
	"github.com/techies/streamify/internal/handler/users"
	"github.com/techies/streamify/internal/utils"
	"github.com/techies/streamify/internal/validation"
	"golang.org/x/crypto/bcrypt"
)

// RegisterRequest represents registration payload
type RegisterRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=8"`
	Username string `json:"username" binding:"required,min=3"`
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
// @Router       /auth/register [post]
func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	var req RegisterRequest
	if err := utils.ParseJSON(w, r, &req); err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid request payload", err)
		return
	}

	// Input Validation (Email, Password, Username)
	if err := validation.ValidateRegisterInput(req.Username, req.Email, req.Password); err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, err.Error(), nil)
		return
	}

	email := utils.NormalizeEmail(req.Email)

	// Check availability
	exists, err := h.app.DB.GetUserByEmail(ctx, email)
	if err != nil {
		// Log the error and return 500 Internal Server Error
		utils.RespondWithError(w, http.StatusInternalServerError, "Database error during registration", err)
		return
	}
	if exists.ID != uuid.Nil {
		utils.RespondWithError(w, http.StatusConflict, "Email already registered", nil)
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "Failed to hash password", err)
		return
	}

	// Verification Logic
	vToken, err := token.GenerateSecureToken(32)
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "Failed to generate verification token", err)
		return
	}
	vExpires := time.Now().Add(24 * time.Hour)

	// Create User with NULL profile fields
	user, err := h.app.DB.CreateUser(ctx, database.CreateUserParams{
		Username:              req.Username,
		Email:                 email,
		PasswordHash:          string(hashedPassword),
		VerificationToken:     sql.NullString{String: vToken, Valid: true},
		VerificationExpiresAt: sql.NullTime{Time: vExpires, Valid: true},
		// All profile fields are explicitly NULL/Valid:false
	})

	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "Registration failed", err)
		return
	}

	utils.RespondWithJSON(w, http.StatusCreated, users.MapUserToResponse(user))
}
