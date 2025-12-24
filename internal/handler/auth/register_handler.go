package auth

import (
	"database/sql"
	"errors"
	"net/http"
	"strings"
	"time"

	"github.com/techies/streamify/internal/database"
	"github.com/techies/streamify/internal/handler/token"
	"github.com/techies/streamify/internal/logger"
	"github.com/techies/streamify/internal/models"
	"github.com/techies/streamify/internal/utils"
	"golang.org/x/crypto/bcrypt"
)

// RegisterRequest represents registration payload
type RegisterRequest struct {
	Email     string `json:"email" binding:"required,email"`
	Password  string `json:"password" binding:"required,min=8"`
	Username  string `json:"username" binding:"required,min=3"`
	FirstName string `json:"first_name" binding:"required,min=1"`
	LastName  string `json:"last_name" binding:"required,min=1"`
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
		logger.Warn(ctx, "Register: invalid request payload", "error", err)
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid request payload", err)
		return
	}

	// 1. Validate Email
	if err := utils.ValidateEmail(req.Email); err != nil {
		logger.Warn(ctx, "Register: invalid email", "email", req.Email, "error", err)
		utils.RespondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	// 2. Validate Password
	if err := utils.ValidatePassword(req.Password); err != nil {
		logger.Warn(ctx, "Register: invalid password", "error", err)
		utils.RespondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	// 3. Validate Username
	if len(strings.TrimSpace(req.Username)) < 3 {
		logger.Warn(ctx, "Register: invalid username", "username", req.Username)
		utils.RespondWithError(w, http.StatusBadRequest, "Username must be at least 3 characters long")
		return
	}

	logger.Debug(ctx, "registration attempt", "email", req.Email)

	// Normalize email
	email := utils.NormalizeEmail(req.Email)

	// Check if user exists
	existingUser, err := h.app.DB.GetUserByEmail(ctx, email)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		utils.RespondWithError(w, http.StatusInternalServerError, "Failed to check existing user", err)
		return
	}
	if existingUser.Email != "" {
		utils.RespondWithError(w, http.StatusConflict, "User with this email already exists", err)
		return
	}
	if len(strings.TrimSpace(req.FirstName)) < 1 {
		logger.Warn(ctx, "Register: missing first name")
		utils.RespondWithError(w, http.StatusBadRequest, "First name is required")
		return
	}
	if len(strings.TrimSpace(req.LastName)) < 1 {
		logger.Warn(ctx, "Register: missing last name")
		utils.RespondWithError(w, http.StatusBadRequest, "Last name is required")
		return
	}
	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		logger.Error(ctx, "failed to hash password", err)
		utils.RespondWithError(w, http.StatusInternalServerError, "Failed to register user", err)
		return
	}
	verificationToken, err := token.GenerateSecureToken(32)
	if err != nil {
		logger.Error(ctx, "failed to generate verification token", err)
		utils.RespondWithError(w, http.StatusInternalServerError, "Failed to generate verification token", err)
		return
	}

	verificationExpiresAt := time.Now().Add(24 * time.Hour)
	// Create user
	newUserParams := &database.CreateUserParams{
		Username:              req.Username,
		Email:                 email,
		PasswordHash:          string(hashedPassword),
		VerificationToken:     sql.NullString{String: verificationToken, Valid: true},
		VerificationExpiresAt: sql.NullTime{Time: verificationExpiresAt, Valid: true},
		FirstName:             sql.NullString{String: req.FirstName, Valid: true},
		LastName:              sql.NullString{String: req.LastName, Valid: true},
	}
	createdUser, err := h.app.DB.CreateUser(ctx, *newUserParams)
	if err != nil {
		logger.Error(ctx, "failed to create user", err, "email", email)
		utils.RespondWithError(w, http.StatusInternalServerError, "Failed to register user", err)
		return
	}

	logger.Info(r.Context(), "user registered successfully", "user_id", createdUser.ID, "email", createdUser.Email)
	utils.RespondWithJSON(w, http.StatusCreated, models.NewUserResponse(&createdUser))
}
