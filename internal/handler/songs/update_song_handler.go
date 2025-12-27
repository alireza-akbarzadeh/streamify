package songs

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/techies/streamify/internal/models"
	"github.com/techies/streamify/internal/service"
	"github.com/techies/streamify/internal/utils"
)

type UpdateSongRequest struct {
	Title    *string `json:"title" validate:"omitempty"`
	ArtistID *string `json:"artist_id" validate:"omitempty,uuid4"`
	AlbumID  *string `json:"album_id" validate:"omitempty,uuid4"`
	Genre    *string `json:"genre" validate:"omitempty"`
	URL      *string `json:"url" validate:"omitempty,url"`
	Duration *int    `json:"duration" validate:"omitempty,gte=0"`
}

// UpdateSongHandler updates a song's details by its ID
// @Summary      Update a song
// @Description  Update a song's details by its ID
// @Tags         Songs
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id    path      string                     true  "Song ID (UUID)"
// @Param        song  body      models.UpdateSongRequest   true  "Song update payload"
// @Success      200   {object}  models.SongResponse
// @Failure      400   {object}  utils.ErrorResponse
// @Failure      404   {object}  utils.ErrorResponse
// @Failure      500   {object}  utils.ErrorResponse
// @Router       /api/v1/songs/{id} [put]
func (h *SongHandler) UpdateSongHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	songIDStr := chi.URLParam(r, "id")
	songID, err := uuid.Parse(songIDStr)
	if err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid song ID format", err)
		return
	}

	var req models.UpdateSongRequest
	if err := utils.ParseJSON(w, r, &req); err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid request payload", err)
		return
	}

	validate := validator.New()
	if err := validate.Struct(req); err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "Validation failed", err)
		return
	}
	albumID, err := utils.ParseUUIDPtr(req.AlbumID)
	if err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid album_id", err)
		return
	}

	artistID, err := utils.ParseUUIDPtr(req.ArtistID)
	if err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid artist_id", err)
		return
	}
	song, appErr := h.Service.UpdateSong(ctx, songID, service.UpdateSongParams{
		Title:    req.Title,
		ArtistID: artistID,
		AlbumID:  albumID,
		Genre:    req.Genre,
		URL:      req.URL,
		Duration: req.Duration,
	})

	if appErr != nil {
		utils.RespondWithError(w, appErr.Code, appErr.Message, appErr.Err)
		return
	}

	utils.RespondWithJSON(w, http.StatusOK, song)
}
