package auth

import (
	"github.com/techies/streamify/internal/app"
	"github.com/techies/streamify/internal/service"
)

type Handler struct {
	App     *app.AppConfig
	Service *service.AuthService
}

func NewAuthHandler(app *app.AppConfig) *Handler {
	return &Handler{App: app}

}
