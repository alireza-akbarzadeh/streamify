package routes

import (
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/httprate"
	"github.com/techies/streamify/internal/app"
	"github.com/techies/streamify/internal/handler"
)

func authRouter(h *handler.Handler, cfg *app.AppConfig) chi.Router {
	r := chi.NewRouter()

	// Stricter rate limit for Auth only
	r.Use(httprate.LimitByIP(5, time.Minute))

	r.Get("/verify", h.Auth.VerifyToken)

	r.Post("/register", h.Auth.Register)
	r.Post("/login", h.Auth.Login)
	r.Post("/refresh", h.Token.RefreshToken)
	r.Post("/logout", h.Token.Logout)
	r.Post("/logout-all", h.Token.LogoutAllDevices)

	return r
}
