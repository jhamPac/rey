package rey

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

type SpyStore struct {
	response  string
	cancelled bool
}

func (s *SpyStore) Fetch() string {
	time.Sleep(100 * time.Millisecond)
	return s.response
}

func (s *SpyStore) Cancel() {
	s.cancelled = true
}

func TestHandler(t *testing.T) {
	data := "A long time ago in a galaxy far, far away"

	t.Run("testing / path", func(t *testing.T) {
		store := &SpyStore{data, false}
		s := Server(store)

		req := httptest.NewRequest(http.MethodGet, "/", nil)
		rr := httptest.NewRecorder()

		s.ServeHTTP(rr, req)

		if rr.Body.String() != data {
			t.Errorf(`got "%s", expected "%s"`, rr.Body.String(), data)
		}

		if store.cancelled {
			t.Error("store should have not cancelled")
		}
	})

	t.Run("store should cancel work when cancelled", func(t *testing.T) {
		store := &SpyStore{data, false}
		s := Server(store)

		req := httptest.NewRequest(http.MethodGet, "/", nil)

		cancellingCtx, cancel := context.WithCancel(req.Context())
		time.AfterFunc(5*time.Millisecond, cancel)
		req = req.WithContext(cancellingCtx)

		rr := httptest.NewRecorder()

		s.ServeHTTP(rr, req)

		if !store.cancelled {
			t.Errorf("store was not cancelled")
		}
	})
}
