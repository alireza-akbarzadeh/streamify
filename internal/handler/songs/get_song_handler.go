package songs

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/techies/streamify/internal/logger"
	"github.com/techies/streamify/internal/models"
	"github.com/techies/streamify/internal/service"
	"github.com/techies/streamify/internal/utils"
)

type SongListResponse struct {
	Songs   []models.SongResponse `json:"songs"`
	Total   int64                 `json:"total"`
	Limit   int32                 `json:"limit"`
	Offset  int32                 `json:"offset"`
	HasMore bool                  `json:"has_more"`
}

// SongListHandler returns a paginated list of songs
// @Summary      List songs
// @Description  Get a paginated list of songs with optional filters
// @Tags         Songs
// @Accept       json
// @Produce      json
// @Param        limit     query     int     false  "Max results per page (default 20, max 100)"
// @Param        offset    query     int     false  "Offset for pagination (default 0)"
// @Param        title     query     string  false  "Filter by song title"
// @Param        artist_id query     string  false  "Filter by artist UUID"
// @Param        genre     query     string  false  "Filter by genre"
// @Success      200       {object}  SongListResponse
// @Failure      500       {object}  utils.ErrorResponse
// @Router       /api/v1/songs [get]
func (h *SongHandler) SongListHandler(w http.ResponseWriter, r *http.Request) {
	limit := utils.ParseInt(r.URL.Query().Get("limit"), 20)
	if limit > 100 {
		limit = 100
	}
	offset := utils.ParseInt(r.URL.Query().Get("offset"), 0)
	result, appErr := h.Service.ListSongs(r.Context(), service.ListSongsParams{
		Limit:    int32(limit),
		Offset:   int32(offset),
		Title:    r.URL.Query().Get("title"),
		ArtistID: r.URL.Query().Get("artist_id"),
		Genre:    r.URL.Query().Get("genre"),
	})
	if appErr != nil {
		utils.RespondWithError(w, appErr.Code, appErr.Message, appErr.Err)
		return
	}

	utils.RespondWithJSON(w, http.StatusOK, SongListResponse{
		Songs:   result.Songs,
		Total:   result.Total,
		Limit:   int32(limit),
		Offset:  int32(offset),
		HasMore: int64(offset+limit) < result.Total,
	})
}

// GetSongHandler returns a song by its ID
// @Summary      Get song by ID
// @Description  Get a single song by its UUID
// @Tags         Songs
// @Accept       json
// @Produce      json
// @Param        id   path      string  true  "Song ID (UUID)"
// @Success      200  {object}  models.SongResponse
// @Failure      400  {object}  utils.ErrorResponse  "Invalid song ID format"
// @Failure      404  {object}  utils.ErrorResponse  "Song not found"
// @Failure      500  {object}  utils.ErrorResponse  "Failed to fetch song"
// @Router       /api/v1/songs/{id} [get]
func (h *SongHandler) GetSongHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	songIDParam := chi.URLParam(r, "id")
	songID, err := uuid.Parse(songIDParam)
	if err != nil {
		logger.Warn(ctx, "GetSongByIDHandler: invalid song ID format", "song_id", songIDParam)
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid song ID format", nil)
		return
	}

	result, appErr := h.Service.GetSong(ctx, songID)
	if appErr != nil {
		utils.RespondWithError(w, appErr.Code, appErr.Message, appErr.Err)
		return
	}
	utils.RespondWithJSON(w, http.StatusOK, result)
}
