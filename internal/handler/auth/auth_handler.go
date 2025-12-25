package auth

import "github.com/techies/streamify/internal/app"

type AuthHandler struct {
	App *app.AppConfig
}

func NewAuthHandler(app *app.AppConfig) *AuthHandler {
	return &AuthHandler{App: app}

}
