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
		data, err := store.Fetch(r.Context())

		if err != nil {
			return
		}

		fmt.Fprint(w, data)
	}
}
