package handler

import (
	"github.com/techies/streamify/internal/app"
	"github.com/techies/streamify/internal/handler/auth"
	"github.com/techies/streamify/internal/handler/token"
	"github.com/techies/streamify/internal/handler/users"
)

type Handler struct {
	App   *app.AppConfig
	Token *token.TokenHandler
	Auth  *auth.AuthHandler
	User  *users.UserHandler
}

func NewHandler(appConfig *app.AppConfig) *Handler {
	return &Handler{
		App:   appConfig,
		Token: token.NewTokenHandler(appConfig),
		Auth:  auth.NewAuthHandler(appConfig),
		User:  users.NewUserHandler(appConfig),
	}
}
