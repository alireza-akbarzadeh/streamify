package routes

import (
	"github.com/go-chi/chi/v5"
	"github.com/techies/streamify/internal/handler"
	"github.com/techies/streamify/internal/middleware"
)

func userRouter(h *handler.Handler) chi.Router {
	r := chi.NewRouter()

	r.Get("/", h.User.UserList)
	r.Get("/{id}", h.User.GetUser)
	r.Put("/{id}", h.User.UpdateProfile)
	r.With(middleware.AdminOnly).Post("/{id}/lock", h.Auth.LockUser)
	r.With(middleware.AdminOnly).Post("/{id}/unlock", h.Auth.UnLockUser)
	r.With(middleware.AdminOnly).Delete("/{id}", h.User.DeleteUser)
	r.With(middleware.AdminOnly).Delete("/old-soft-deleted", h.User.PermanentlyDeleteOldSoftDeletedUsers)

	return r
}
