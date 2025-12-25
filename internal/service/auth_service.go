package service

import (
	"context"
	"database/sql"
	"errors"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/techies/streamify/internal/app"
	"github.com/techies/streamify/internal/database"
	"github.com/techies/streamify/internal/handler/token"
	"github.com/techies/streamify/internal/utils"
	"golang.org/x/crypto/bcrypt"
)

type AuthService struct {
	BaseService
	cfg *app.AppConfig
}

func NewAuthService(db *database.Queries, cfg *app.AppConfig) *AuthService {
	return &AuthService{
		BaseService: NewBaseService(db),
		cfg:         cfg,
	}
}

type RegisterParams struct {
	Username string `validate:"required,min=3,max=30"`
	Email    string `validate:"required,email"`
	Password string `validate:"required,min=8"`
}

func (s *AuthService) Register(ctx context.Context, params RegisterParams) (database.User, *utils.AppError) {
	if err := validate.Struct(params); err != nil {
		return database.User{}, &utils.AppError{
			Code:    http.StatusBadRequest,
			Message: "Validation failed",
			Err:     err,
		}
	}

	email := utils.NormalizeEmail(params.Email)

	// Check availability
	exists, err := s.DB.GetUserByEmail(ctx, email)
	if err == nil && exists.ID != uuid.Nil {
		return database.User{}, &utils.AppError{
			Code:    http.StatusConflict,
			Message: "Email already registered",
		}
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(params.Password), bcrypt.DefaultCost)
	if err != nil {
		return database.User{}, &utils.AppError{
			Code:    http.StatusInternalServerError,
			Message: "Failed to hash password",
			Err:     err,
		}
	}

	vToken, _ := token.GenerateSecureToken(32)
	vExpires := time.Now().Add(24 * time.Hour)

	user, err := s.DB.CreateUser(ctx, database.CreateUserParams{
		Username:              params.Username,
		Email:                 email,
		PasswordHash:          string(hashedPassword),
		VerificationToken:     sql.NullString{String: vToken, Valid: true},
		VerificationExpiresAt: sql.NullTime{Time: vExpires, Valid: true},
	})

	if err != nil {
		return database.User{}, &utils.AppError{
			Code:    http.StatusInternalServerError,
			Message: "Registration failed",
			Err:     err,
		}
	}

	return user, nil
}

type LoginParams struct {
	Email     string `validate:"required,email"`
	Password  string `validate:"required"`
	IP        string
	UserAgent string
}

type LoginResult struct {
	User         database.User
	AccessToken  string
	RefreshToken string
}

func (s *AuthService) Login(ctx context.Context, params LoginParams) (LoginResult, *utils.AppError) {
	if err := validate.Struct(params); err != nil {
		return LoginResult{}, &utils.AppError{
			Code:    http.StatusBadRequest,
			Message: "Validation failed",
			Err:     err,
		}
	}

	user, err := s.DB.GetUserByEmail(ctx, params.Email)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return LoginResult{}, &utils.AppError{
				Code:    http.StatusUnauthorized,
				Message: "Invalid email or password",
			}
		}
		return LoginResult{}, &utils.AppError{
			Code:    http.StatusInternalServerError,
			Message: "Database error",
			Err:     err,
		}
	}

	if user.IsLocked {
		return LoginResult{}, &utils.AppError{
			Code:    http.StatusUnauthorized,
			Message: "User account is locked",
		}
	}

	if !user.IsVerified {
		return LoginResult{}, &utils.AppError{
			Code:    http.StatusUnauthorized,
			Message: "Please verify your email before logging in",
		}
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(params.Password)); err != nil {
		return LoginResult{}, &utils.AppError{
			Code:    http.StatusUnauthorized,
			Message: "Invalid email or password",
		}
	}

	refreshToken, _ := token.GenerateSecureToken(token.RefreshTokenLen)
	session, err := s.DB.CreateSession(ctx, database.CreateSessionParams{
		UserID:       user.ID,
		RefreshToken: refreshToken,
		IpAddress:    utils.ToNullString(&params.IP),
		UserAgent:    utils.ToNullString(&params.UserAgent),
		ExpiresAt:    time.Now().Add(token.RefreshTokenTTL),
	})
	if err != nil {
		return LoginResult{}, &utils.AppError{
			Code:    http.StatusInternalServerError,
			Message: "Failed to save session",
			Err:     err,
		}
	}

	accessToken, err := utils.GenerateToken(user.ID, session.ID, token.AccessTokenTTL, s.cfg.JWTSecret)
	if err != nil {
		return LoginResult{}, &utils.AppError{
			Code:    http.StatusInternalServerError,
			Message: "Failed to generate access token",
			Err:     err,
		}
	}

	return LoginResult{
		User:         user,
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}
