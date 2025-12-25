package users

import (
	"github.com/techies/streamify/internal/app"
	"github.com/techies/streamify/internal/database"
	"github.com/techies/streamify/internal/models"
	"github.com/techies/streamify/internal/service"
)

type UserHandler struct {
	App     *app.AppConfig
	Service *service.UserService
}

func NewUserHandler(app *app.AppConfig) *UserHandler {
	return &UserHandler{App: app}
}

// MapUserListToResponse handles slices of users efficiently
func MapUserListToResponse(dbUsers []database.User) []models.UserResponse {
	users := make([]models.UserResponse, len(dbUsers))
	for i, u := range dbUsers {
		users[i] = MapUserToResponse(u)
	}
	return users
}

func MapUserToResponse(u database.User) models.UserResponse {
	return models.UserResponse{
		ID:          u.ID,
		Username:    u.Username,
		Email:       u.Email,
		IsVerified:  u.IsVerified,
		CreatedAt:   u.CreatedAt,
		UpdatedAt:   u.UpdatedAt,
		Bio:         u.Bio.String,
		PhoneNumber: u.PhoneNumber.String,
		AvatarUrl:   u.AvatarUrl.String,
	}
}

type UserListResponse struct {
	Users   []models.UserResponse `json:"users"`
	Total   int64                 `json:"total"`
	Limit   int32                 `json:"limit"`
	Offset  int32                 `json:"offset"`
	HasMore bool                  `json:"has_more"`
}
