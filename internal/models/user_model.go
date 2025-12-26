package models

import (
	"time"

	"github.com/google/uuid"
	"github.com/techies/streamify/internal/database"
)

// UserResponse omits password and password hash
type UserResponse struct {
	ID          uuid.UUID         `json:"id"`
	Username    string            `json:"username"`
	FirstName   string            `json:"first_name"`
	LastName    string            `json:"last_name"`
	IsLocked    bool              `json:"is_locked"`
	Email       string            `json:"email"`
	IsVerified  bool              `json:"is_verified"`
	Status      string            `json:"status"`
	CreatedAt   time.Time         `json:"created_at"`
	UpdatedAt   time.Time         `json:"updated_at"`
	AvatarUrl   string            `json:"avatar_url,omitempty"`
	Bio         string            `json:"bio,omitempty"`
	PhoneNumber string            `json:"phone_number,omitempty"`
	Role        database.UserRole `json:"role"`
}

type UpdateSongRequest struct {
	Title       *string `json:"title" validate:"omitempty,min=1"`
	ArtistID    *string `json:"artist_id" validate:"omitempty,uuid4"`
	AlbumID     *string `json:"album_id" validate:"omitempty,uuid4"`
	Genre       *string `json:"genre" validate:"omitempty,min=1"`
	URL         *string `json:"url" validate:"omitempty,url"`
	Duration    *int    `json:"duration" validate:"omitempty,gte=0"`
	ReleaseDate *string `json:"release_date" validate:"omitempty,datetime=2006-01-02"`
	Bitrate     *int    `json:"bitrate" validate:"omitempty,gte=0"`
}

func NewUserResponse(u *database.User) *UserResponse {
	return &UserResponse{
		ID:          u.ID,
		Username:    u.Username,
		FirstName:   u.FirstName.String,
		LastName:    u.LastName.String,
		IsLocked:    u.IsLocked,
		Email:       u.Email,
		IsVerified:  u.IsVerified,
		Status:      u.Status,
		Role:        u.Role,
		AvatarUrl:   u.AvatarUrl.String,
		Bio:         u.Bio.String,
		PhoneNumber: u.PhoneNumber.String,
		CreatedAt:   u.CreatedAt,
		UpdatedAt:   u.UpdatedAt,
	}
}
