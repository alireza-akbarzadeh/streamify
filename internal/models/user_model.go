package models

import (
	"time"

	"github.com/techies/streamify/internal/database"
)

// UserResponse omits password and password hash
type UserResponse struct {
	ID         string    `json:"id"`
	Username   string    `json:"username"`
	FirstName  string    `json:"first_name"`
	LastName   string    `json:"last_name"`
	IsLocked   bool      `json:"is_locked"`
	Email      string    `json:"email"`
	IsVerified bool      `json:"is_verified"`
	Status     string    `json:"status"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

func NewUserResponse(u *database.User) *UserResponse {
	return &UserResponse{
		ID:         u.ID.String(),
		Username:   u.Username,
		FirstName:  u.FirstName.String,
		LastName:   u.LastName.String,
		IsLocked:   u.IsLocked,
		Email:      u.Email,
		IsVerified: u.IsVerified,
		Status:     u.Status,
		CreatedAt:  u.CreatedAt,
		UpdatedAt:  u.UpdatedAt,
	}
}
