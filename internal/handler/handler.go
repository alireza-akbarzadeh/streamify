package handler

import (
	"github.com/techies/streamify/internal/app"
	"github.com/techies/streamify/internal/handler/auth"
	"github.com/techies/streamify/internal/handler/token"
	"github.com/techies/streamify/internal/handler/users"
)

type Handler struct {
	app   *app.AppConfig
	Token *token.TokenHandler
	Auth  *auth.AuthHandler
	User  *users.UserHandler
}

func NewHandler(app *app.AppConfig) *Handler {
	return &Handler{
		app:   app,
		Token: token.NewTokenHandler(app),
		Auth:  auth.NewAuthHandler(app),
		User:  users.NewUserHandler(app),
	}
}
