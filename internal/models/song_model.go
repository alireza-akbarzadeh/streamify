package models

import (
	"time"

	"github.com/google/uuid"
	"github.com/techies/streamify/internal/database"
	"github.com/techies/streamify/internal/utils"
)

type SongResponse struct {
	ID          uuid.UUID  `json:"id"`
	Title       string     `json:"title"`
	ArtistID    uuid.UUID  `json:"artist_id"`
	AlbumID     *uuid.UUID `json:"album_id,omitempty"`
	Genre       *string    `json:"genre,omitempty"`
	URL         string     `json:"url"`
	Duration    int        `json:"duration"`
	ReleaseDate *time.Time `json:"release_date,omitempty"`
	Bitrate     *int       `json:"bitrate,omitempty"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
	DeletedAt   *time.Time `json:"deleted_at,omitempty"`
}

func MapSongResponse(s *database.Song) *SongResponse {
	return &SongResponse{
		ID:          s.ID,
		Title:       s.Title,
		ArtistID:    s.ArtistID.UUID,
		AlbumID:     utils.UuidPtr(s.AlbumID),
		Genre:       utils.StrPtr(s.Genre),
		URL:         s.Url.String,
		Duration:    int(s.Duration.Int32),
		ReleaseDate: utils.TimePtr(s.ReleaseDate),
		Bitrate:     utils.IntPtr(s.Bitrate),
		CreatedAt:   s.CreatedAt,
		UpdatedAt:   s.UpdatedAt.Time,
		DeletedAt:   utils.TimePtr(s.DeletedAt),
	}
}

// MapSongListToResponse maps a slice of database.Song to a slice of SongResponse
func MapSongListToResponse(songs []database.Song) []SongResponse {
	responses := make([]SongResponse, 0, len(songs))
	for _, song := range songs {
		responses = append(responses, *MapSongResponse(&song))
	}
	return responses
}
