package songs

import (
	"errors"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/techies/streamify/internal/service"
	"github.com/techies/streamify/internal/utils"
)

// CreateSongHandler handles creating a new song
// @Summary      Create a new song
// @Description  Creates a new song in the system
// @Tags         Songs
// @Accept       json
// @Produce      json
// @Param        song  body      service.CreateSongParams  true  "Song data"
// @Success      201   {object}  models.SongResponse
// @Failure      400   {object}  utils.ErrorResponse  "Invalid input"
// @Failure      500   {object}  utils.ErrorResponse  "Failed to create song"
// @Router       /api/v1/songs [post]
func (h *SongHandler) CreateSongHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	var req service.CreateSongParams
	if err := utils.ParseJSON(w, r, &req); err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid request body", err)
		return
	}

	validate := validator.New()
	if err := validate.Struct(req); err != nil {
		if validationErrors, ok := err.(validator.ValidationErrors); ok {
			utils.RespondWithError(w, http.StatusBadRequest, "Validation failed", errors.New(validationErrors.Error()))
			return
		}
		utils.RespondWithError(w, http.StatusBadRequest, "Validation error", err)
		return
	}

	song, appErr := h.Service.CreateSong(ctx, req)
	if appErr != nil {
		utils.RespondWithError(w, appErr.Code, appErr.Message, appErr.Err)
		return
	}

	utils.RespondWithJSON(w, http.StatusCreated, song)
}
