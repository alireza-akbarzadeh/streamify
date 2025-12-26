package songs

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/techies/streamify/internal/logger"
	"github.com/techies/streamify/internal/utils"
)

// DeleteSongHandler soft-deletes a song by its ID
// @Summary      Delete a song
// @Description  Soft-delete a song by its ID
// @Tags         Songs
// @Security     BearerAuth
// @Param        id   path      string  true  "Song ID (UUID)"
// @Success      204  "Song deleted successfully"
// @Failure      400  {object}  utils.ErrorResponse  "Invalid song ID format"
// @Failure      404  {object}  utils.ErrorResponse  "Song not found"
// @Failure      500  {object}  utils.ErrorResponse  "Failed to delete song"
// @Router       /api/v1/songs/{id} [delete]
func (h *SongHandler) DeleteSongHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	songIDParam := chi.URLParam(r, "id")
	songID, err := uuid.Parse(songIDParam)
	if err != nil {
		logger.Warn(ctx, "DeleteSongHandler: invalid song ID format", "song_id", songIDParam)
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid song ID format", nil)
		return
	}

	appErr := h.Service.SoftDeleteSong(ctx, songID)
	if appErr != nil {
		utils.RespondWithError(w, appErr.Code, appErr.Message, appErr.Err)
		return
	}

	logger.Info(ctx, "DeleteSongHandler: song deleted successfully", "song_id", songID)
	utils.RespondWithJSON(w, http.StatusNoContent, nil)
}
