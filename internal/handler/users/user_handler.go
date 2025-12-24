package users

import (
	"time"

	"github.com/google/uuid"
	"github.com/techies/streamify/internal/app"
	"github.com/techies/streamify/internal/database"
)

type UserHandler struct {
	app *app.AppConfig
}

func NewUserHandler(app *app.AppConfig) *UserHandler {
	return &UserHandler{app: app}
}

// MapUserListToResponse handles slices of users efficiently
func MapUserListToResponse(dbUsers []database.User) []UserResponse {
	users := make([]UserResponse, len(dbUsers))
	for i, u := range dbUsers {
		users[i] = MapUserToResponse(u)
	}
	return users
}

func MapUserToResponse(u database.User) UserResponse {
	return UserResponse{
		ID:         u.ID,
		Username:   u.Username,
		Email:      u.Email,
		IsVerified: u.IsVerified, // Updated to match your logic
		// Using RFC3339 is standard for modern APIs
		CreatedAt: u.CreatedAt.Format(time.RFC3339),
		UpdatedAt: u.UpdatedAt.Format(time.RFC3339),
	}
}

// UserResponse is the public contract
type UserResponse struct {
	ID         uuid.UUID `json:"id"`
	Username   string    `json:"username"`
	Email      string    `json:"email"`
	IsVerified bool      `json:"is_verified"`
	CreatedAt  string    `json:"created_at"`
	UpdatedAt  string    `json:"updated_at"`
}

type UserListResponse struct {
	Users   []UserResponse `json:"users"`
	Total   int64          `json:"total"`
	Limit   int32          `json:"limit"`
	Offset  int32          `json:"offset"`
	HasMore bool           `json:"has_more"`
}
