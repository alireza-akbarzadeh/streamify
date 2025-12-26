package handler

import (
	"github.com/techies/streamify/internal/app"
	"github.com/techies/streamify/internal/handler/auth"
	"github.com/techies/streamify/internal/handler/songs"
	"github.com/techies/streamify/internal/handler/token"
	"github.com/techies/streamify/internal/handler/users"
	"github.com/techies/streamify/internal/service"
)

type Handler struct {
	App     *app.AppConfig
	Token   *token.TokenHandler
	Auth    *auth.Handler
	User    *users.UserHandler
	Song    *songs.SongHandler
	Service struct {
		Auth *service.AuthService
		User *service.UserService
		Song *service.SongService
	}
}

func NewHandler(appConfig *app.AppConfig) *Handler {
	authService := service.NewAuthService(appConfig.DB, appConfig)
	userService := service.NewUserService(appConfig.DB, appConfig)

	h := &Handler{
		App:   appConfig,
		Token: token.NewTokenHandler(appConfig),
		Auth:  auth.NewAuthHandler(appConfig),
		User:  users.NewUserHandler(appConfig),
		Song:  songs.NewSongHandler(appConfig),
	}
	h.Service.Auth = authService
	h.Service.User = userService

	// Pass services to handlers if needed or keep them accessible via h.Service
	h.Auth.Service = authService
	h.User.Service = userService
	h.Song.Service = h.Service.Song

	return h
}
