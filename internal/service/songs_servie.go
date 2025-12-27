package service

import (
	"context"
	"database/sql"

	"github.com/google/uuid"
	"github.com/techies/streamify/internal/app"
	"github.com/techies/streamify/internal/database"
	"github.com/techies/streamify/internal/models"
	"github.com/techies/streamify/internal/utils"
)

type SongService struct {
	BaseService
	cfg *app.AppConfig
}

func NewSongService(db *database.Queries, cfg *app.AppConfig) *SongService {
	return &SongService{
		BaseService: NewBaseService(db),
		cfg:         cfg,
	}
}

type CreateSongParams struct {
	Title    string    `validate:"required"`
	ArtistID uuid.UUID `validate:"required"`
	AlbumID  *uuid.UUID
	Genre    *string
	URL      string `validate:"required,url"`
	Duration int    `validate:"required,gte=0"`
}

// CreateSong creates a new song in the database
func (s *SongService) CreateSong(ctx context.Context, params CreateSongParams) (models.SongResponse, *utils.AppError) {

	result, err := s.DB.CreateSong(ctx, database.CreateSongParams{
		Title:    params.Title,
		ArtistID: uuid.NullUUID{UUID: params.ArtistID, Valid: true},
		AlbumID:  utils.ToNullUUID(params.AlbumID),
		Genre:    utils.ToNullString(params.Genre),
		Url:      sql.NullString{String: params.URL, Valid: true},
		Duration: sql.NullInt32{Int32: int32(params.Duration), Valid: true},
	})
	if err != nil {
		return models.SongResponse{}, &utils.AppError{
			Code:    500,
			Message: "Failed to create song",
			Err:     err,
		}
	}

	return *models.MapSongResponse(&result), nil
}

// SoftDeleteSong marks a song as deleted in the database
func (s *SongService) SoftDeleteSong(ctx context.Context, songID uuid.UUID) *utils.AppError {
	err := s.DB.SoftDeleteSong(ctx, songID)
	if err != nil {
		if err == sql.ErrNoRows {
			return &utils.AppError{
				Code:    404,
				Message: "Song not found",
				Err:     err,
			}
		}
		return &utils.AppError{
			Code:    500,
			Message: "Failed to delete song",
			Err:     err,
		}
	}
	return nil
}

// GetSong retrieves a song by its ID
func (s *SongService) GetSong(ctx context.Context, songID uuid.UUID) (models.SongResponse, *utils.AppError) {
	song, err := s.DB.GetSongByID(ctx, songID)
	if err != nil {
		if err == sql.ErrNoRows {
			return models.SongResponse{}, &utils.AppError{
				Code:    404,
				Message: "Song not found",
				Err:     err,
			}
		}
		return models.SongResponse{}, &utils.AppError{
			Code:    500,
			Message: "Failed to fetch song",
			Err:     err,
		}
	}

	return *models.MapSongResponse(&song), nil
}

type ListSongsParams struct {
	Limit    int32
	Offset   int32
	Title    string
	ArtistID string
	Genre    string
}

type ListSongsResponse struct {
	Songs []models.SongResponse
	Total int64
}

// ListSongs retrieves a list of songs based on the provided parameters
func (s *SongService) ListSongs(ctx context.Context, params ListSongsParams) (ListSongsResponse, *utils.AppError) {
	arg := database.ListSongsParams{
		Limit:          params.Limit,
		Offset:         params.Offset,
		SearchTitle:    sql.NullString{String: params.Title, Valid: params.Title != ""},
		SearchArtistID: utils.ToNullUUIDFromString(params.ArtistID),
		SearchGenre:    sql.NullString{String: params.Genre, Valid: params.Genre != ""},
	}

	dbSongs, err := s.DB.ListSongs(ctx, arg)
	if err != nil {
		return ListSongsResponse{}, &utils.AppError{
			Code:    500,
			Message: "Failed to fetch songs",
			Err:     err,
		}
	}

	total, err := s.DB.CountSongs(ctx, database.CountSongsParams{
		SearchTitle:    arg.SearchTitle,
		SearchArtistID: arg.SearchArtistID,
		SearchGenre:    arg.SearchGenre,
	})
	if err != nil {
		return ListSongsResponse{}, &utils.AppError{
			Code:    500,
			Message: "Failed to count songs",
			Err:     err,
		}
	}

	var songResponses []models.SongResponse
	for _, song := range dbSongs {
		songResponses = append(songResponses, *models.MapSongResponse(&song))
	}

	return ListSongsResponse{
		Songs: songResponses,
		Total: total,
	}, nil
}

// UpdateSong updates the details of an existing song

type UpdateSongParams struct {
	Title    *string
	ArtistID *uuid.UUID
	AlbumID  *uuid.UUID
	Genre    *string
	URL      *string
	Duration *int
}

func (s *SongService) UpdateSong(
	ctx context.Context,
	songID uuid.UUID,
	params UpdateSongParams,
) (models.SongResponse, *utils.AppError) {

	updatedSong, err := s.DB.UpdateSong(ctx, database.UpdateSongParams{
		ID:       songID,
		Title:    utils.ToNullString(params.Title),
		ArtistID: utils.ToNullUUID(params.ArtistID),
		AlbumID:  utils.ToNullUUID(params.AlbumID),
		Genre:    utils.ToNullString(params.Genre),
		Url:      utils.ToNullString(params.URL),
		Duration: utils.ToNullInt32(params.Duration),
	})
	if err != nil {
		if err == sql.ErrNoRows {
			return models.SongResponse{}, &utils.AppError{
				Code:    404,
				Message: "Song not found",
				Err:     err,
			}
		}
		return models.SongResponse{}, &utils.AppError{
			Code:    500,
			Message: "Failed to update song",
			Err:     err,
		}
	}

	return *models.MapSongResponse(&updatedSong), nil
}
