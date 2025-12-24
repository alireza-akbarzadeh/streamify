package handler

import (
	"context"
	"net/http"
	"time"

	"github.com/techies/streamify/internal/utils"
)

// /healthz → liveness (NO dependencies)
func (h *Handler) HandleHealthz(w http.ResponseWriter, r *http.Request) {
	utils.RespondWithJSON(w, http.StatusOK, map[string]string{
		"status": "ok",
	})
}

// /readyz → readiness (checks dependencies)
func (h *Handler) HandleReadyz(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), time.Second)
	defer cancel()

	if err := h.app.Conn.PingContext(ctx); err != nil {
		utils.RespondWithJSON(w, http.StatusServiceUnavailable, map[string]string{
			"status": "not ready",
		})
		return
	}

	utils.RespondWithJSON(w, http.StatusOK, map[string]string{
		"status": "ready",
	})
}
