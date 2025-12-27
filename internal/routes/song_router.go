package routes

import (
	"github.com/go-chi/chi/v5"
	"github.com/techies/streamify/internal/handler"
)

func songRouter(h *handler.Handler) chi.Router {
	r := chi.NewRouter()

	// Song endpoints
	r.Get("/", h.Song.SongListHandler)
	r.Post("/", h.Song.CreateSongHandler)
	r.Get("/{id}", h.Song.GetSongHandler)
	r.Put("/{id}", h.Song.UpdateSongHandler)
	r.Delete("/{id}", h.Song.DeleteSongHandler)

	// (Optional) admin-only cleanup for songs, if needed
	// r.With(middleware.AdminOnly).Delete("/old-soft-deleted", h.Song.PermanentlyDeleteOldSoftDeletedSongs)

	return r
}
