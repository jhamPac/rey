package rey

import (
	"fmt"
	"net/http"
)

// Store for data
type Store interface {
	Fetch() string
	Cancel()
}

// Server creates a simple server to test
func Server(store Store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		store.Cancel()
		fmt.Fprint(w, store.Fetch())
	}
}
