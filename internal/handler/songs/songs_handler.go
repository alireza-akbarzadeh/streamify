package songs

import (
	"github.com/techies/streamify/internal/app"
	"github.com/techies/streamify/internal/service"
)

type SongHandler struct {
	App     *app.AppConfig
	Service *service.SongService
}

func NewSongHandler(app *app.AppConfig) *SongHandler {
	return &SongHandler{App: app}
}
