package routes

import (
	"github.com/go-chi/chi/v5"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/techies/streamify/internal/handler"
)

func HealthRoutes(r chi.Router, h *handler.Handler) {
	r.Get("/healthz", h.HandleHealthz)
	r.Get("/readyz", h.HandleReadyz)
	r.Get("/metrics", promhttp.Handler().ServeHTTP)
}
