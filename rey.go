package rey

import (
	"context"
	"net/http"
)

// Store for data
type Store interface {
	Fetch(ctx context.Context) (string, error)
	Cancel()
}

// Server creates a simple server to test
func Server(store Store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

	}
}
