package rey

import (
	"context"
	"fmt"
	"net/http"
)

// Store for data
type Store interface {
	Fetch(ctx context.Context) (string, error)
}

// Server creates a simple server to test
func Server(store Store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		data, _ := store.Fetch(r.Context())
		fmt.Fprint(w, data)
	}
}
