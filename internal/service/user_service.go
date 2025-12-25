package service

import (
	"context"
	"database/sql"
	"errors"
	"net/http"

	"github.com/google/uuid"
	"github.com/techies/streamify/internal/app"
	"github.com/techies/streamify/internal/database"
	"github.com/techies/streamify/internal/utils"
)

type UserService struct {
	BaseService
	cfg *app.AppConfig
}

func NewUserService(db *database.Queries, cfg *app.AppConfig) *UserService {
	return &UserService{
		BaseService: NewBaseService(db),
		cfg:         cfg,
	}
}

type UpdateProfileParams struct {
	UserID      uuid.UUID `validate:"required"`
	FirstName   *string
	LastName    *string
	Bio         *string
	AvatarUrl   *string
	PhoneNumber *string
}

func (s *UserService) UpdateProfile(ctx context.Context, params UpdateProfileParams) (database.User, *utils.AppError) {
	if err := validate.Struct(params); err != nil {
		return database.User{}, &utils.AppError{
			Code:    http.StatusBadRequest,
			Message: "Validation failed",
			Err:     err,
		}
	}

	err := s.DB.UpdateUserProfile(ctx, database.UpdateUserProfileParams{
		ID:          params.UserID,
		FirstName:   utils.ToNullString(params.FirstName),
		LastName:    utils.ToNullString(params.LastName),
		Bio:         utils.ToNullString(params.Bio),
		AvatarUrl:   utils.ToNullString(params.AvatarUrl),
		PhoneNumber: utils.ToNullString(params.PhoneNumber),
	})

	if err != nil {
		return database.User{}, &utils.AppError{
			Code:    http.StatusInternalServerError,
			Message: "Failed to update profile",
			Err:     err,
		}
	}

	user, err := s.DB.GetUserById(ctx, params.UserID)
	if err != nil {
		return database.User{}, &utils.AppError{
			Code:    http.StatusInternalServerError,
			Message: "Failed to fetch updated user",
			Err:     err,
		}
	}

	return user, nil
}

func (s *UserService) GetUser(ctx context.Context, id uuid.UUID) (database.User, *utils.AppError) {
	user, err := s.DB.GetUserById(ctx, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return database.User{}, &utils.AppError{
				Code:    http.StatusNotFound,
				Message: "User not found",
			}
		}
		return database.User{}, &utils.AppError{
			Code:    http.StatusInternalServerError,
			Message: "Database error",
			Err:     err,
		}
	}
	return user, nil
}

type ListUsersParams struct {
	Limit       int32
	Offset      int32
	Username    string
	Email       string
	PhoneNumber string
}

type ListUsersResult struct {
	Users []database.User
	Total int64
}

func (s *UserService) ListUsers(ctx context.Context, params ListUsersParams) (ListUsersResult, *utils.AppError) {
	arg := database.GetUsersParams{
		Limit:             params.Limit,
		Offset:            params.Offset,
		SearchUsername:    sql.NullString{String: params.Username, Valid: params.Username != ""},
		SearchEmail:       sql.NullString{String: params.Email, Valid: params.Email != ""},
		SearchPhoneNumber: sql.NullString{String: params.PhoneNumber, Valid: params.PhoneNumber != ""},
	}

	dbUsers, err := s.DB.GetUsers(ctx, arg)
	if err != nil {
		return ListUsersResult{}, &utils.AppError{
			Code:    http.StatusInternalServerError,
			Message: "Failed to fetch users",
			Err:     err,
		}
	}

	total, err := s.DB.CountUsers(ctx, database.CountUsersParams{
		SearchUsername:    arg.SearchUsername,
		SearchEmail:       arg.SearchEmail,
		SearchPhoneNumber: arg.SearchPhoneNumber,
	})
	if err != nil {
		return ListUsersResult{}, &utils.AppError{
			Code:    http.StatusInternalServerError,
			Message: "Failed to count users",
			Err:     err,
		}
	}

	return ListUsersResult{
		Users: dbUsers,
		Total: total,
	}, nil
}

func (s *UserService) UpdateUserRole(
	ctx context.Context,
	userID uuid.UUID,
	role database.UserRole,
) *utils.AppError {

	// 1️⃣ Validate role (domain rule)
	switch role {
	case database.UserRoleUser,
		database.UserRoleCustomer,
		database.UserRoleAdmin:
		// owner intentionally excluded
	default:
		return &utils.AppError{
			Code:    http.StatusBadRequest,
			Message: "Invalid or forbidden user role",
		}
	}

	// 2️⃣ Execute DB update
	err := s.DB.UpdateUserRole(ctx, database.UpdateUserRoleParams{
		ID:   userID,
		Role: role,
	})

	if err != nil {
		return &utils.AppError{
			Code:    http.StatusInternalServerError,
			Message: "Failed to update user role",
			Err:     err,
		}
	}

	return nil
}
