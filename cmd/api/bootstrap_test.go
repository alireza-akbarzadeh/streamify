package main

import "testing"

func TestBootstrap_FailsGracefully(t *testing.T) {
	err := bootstrap()

	if err == nil {
		t.Fatal("expected error, got nil")
	}
}
