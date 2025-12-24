package app

import (
	"database/sql"
	"errors"
	"log"
	"net/http"
	"os"
	"time"

	_ "github.com/lib/pq"
	"github.com/techies/streamify/internal/database"
	"github.com/techies/streamify/internal/utils"
)

type AppConfig struct {
	DB             *database.Queries
	Conn           *sql.DB
	Server         *http.Server
	JWTSecret      string
	FrontendURL    string
	AllowedOrigins []string
}

func New() (*AppConfig, error) {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	dbURL := os.Getenv("DB_URL")
	if dbURL == "" {
		return nil, errors.New("DB_URL is required")
	}
	jwt := os.Getenv("JWT_SECRET")
	if jwt == "" {
		return nil, errors.New("JWT_SECRET is required")
	}

	conn, err := sql.Open("postgres", dbURL)
	if err != nil {
		return nil, err
	}

	// Verify DB early
	if err := conn.Ping(); err != nil {
		return nil, err
	}

	return &AppConfig{
		DB:          database.New(conn),
		Conn:        conn,
		JWTSecret:   jwt,
		FrontendURL: utils.GetEnvString("FRONTEND_URL", "http://localhost:3000"),
		Server: &http.Server{
			Addr:         ":" + port,
			ReadTimeout:  10 * time.Second,
			WriteTimeout: 10 * time.Second,
			IdleTimeout:  60 * time.Second,
		},
		AllowedOrigins: []string{
			"http://localhost:3000",
			"https://yourdomain.com",
		},
	}, nil
}

func (a *AppConfig) Close() {
	if a.Conn == nil {
		return
	}

	if err := a.Conn.Close(); err != nil {
		log.Printf("app shutdown: failed to close db connection: %v", err)
	}
}
