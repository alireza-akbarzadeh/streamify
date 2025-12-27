package routes

import (
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/go-chi/httprate"
	httpSwagger "github.com/swaggo/http-swagger"
	_ "github.com/techies/streamify/docs" // Add this blank import
	"github.com/techies/streamify/internal/app"
	"github.com/techies/streamify/internal/handler"
	internalMiddleware "github.com/techies/streamify/internal/middleware"
	"github.com/techies/streamify/internal/utils"
)

func SetupRoutes(h *handler.Handler, cfg *app.AppConfig) http.Handler {
	r := chi.NewRouter()

	r.NotFound(func(w http.ResponseWriter, r *http.Request) {
		utils.RespondWithError(w, http.StatusNotFound, "Route not found", nil)
	})
	r.MethodNotAllowed(func(w http.ResponseWriter, r *http.Request) {
		utils.RespondWithError(w, http.StatusMethodNotAllowed, "Method not allowed", nil)
	})

	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(60 * time.Second))
	r.Use(middleware.Compress(5))

	// Global Rate Limiting
	r.Use(httprate.LimitByIP(100, time.Minute))

	// CORS Configuration
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   cfg.AllowedOrigins, // Move these to your AppConfig
		AllowedMethods:   []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300,
	}))

	// Public Infrastructure
	// --- 4. PUBLIC / INFRA ROUTES ---
	r.Get("/health", h.HandleHealthz)
	r.Get("/swagger/*", httpSwagger.Handler(
		httpSwagger.AfterScript(`
		window.onload = function() {
			const timer = setInterval(function() {
				const ui = window.ui;
				if (ui) {
					clearInterval(timer);
					const originalRequestInterceptor = ui.getConfigs().requestInterceptor || (r => r);
					ui.getConfigs().requestInterceptor = (req) => {
						if (req.headers['Authorization'] && !req.headers['Authorization'].startsWith('Bearer ')) {
							req.headers['Authorization'] = 'Bearer ' + req.headers['Authorization'];
						}
						return originalRequestInterceptor(req);
					};
				}
			}, 100);
		};
		`),
	))

	// --- 5. VERSIONED API ---
	r.Route("/api/v1", func(r chi.Router) {
		// General API Rate Limiting
		r.Use(httprate.LimitByIP(200, time.Minute))

		// Authentication Domain
		r.Mount("/auth", authRouter(h, cfg))
		r.Mount("/songs", songRouter(h))
		// Protected Domain
		r.Group(func(r chi.Router) {
			r.Use(internalMiddleware.AuthMiddleware(h.App.DB, cfg.JWTSecret))
			r.Mount("/users", userRouter(h))
		})
	})

	return r
}
