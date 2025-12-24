package server

import (
	"context"
	"net/http"
	"testing"
	"time"
)

func TestRun_StartsAndShutsDown(t *testing.T) {
	srv := &http.Server{Addr: "127.0.0.1:0"}

	// Shut down server after short delay
	go func() {
		time.Sleep(100 * time.Millisecond)
		err := srv.Shutdown(context.Background()) // graceful shutdown
		if err != nil {
			t.Errorf("Shutdown returned error: %v", err)
		}
	}()

	err := Run(srv)
	if err != nil && err != http.ErrServerClosed {
		t.Fatalf("Run returned error: %v", err)
	}
}
