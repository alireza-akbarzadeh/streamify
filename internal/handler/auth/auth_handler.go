package auth

import "github.com/techies/streamify/internal/app"

type AuthHandler struct {
	app *app.AppConfig
}

func NewAuthHandler(app *app.AppConfig) *AuthHandler {
	return &AuthHandler{app: app}

}
