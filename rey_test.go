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
	t         *testing.T
}

func (s *SpyStore) Fetch() string {
	time.Sleep(100 * time.Millisecond)
	return s.response
}

func (s *SpyStore) Cancel() {
	s.cancelled = true
}

func (s *SpyStore) assertWasCancelled() {
	s.t.Helper()
	if !s.cancelled {
		s.t.Error("store was not told to cancel")
	}
}

func (s *SpyStore) assertWasNotCancelled() {
	s.t.Helper()
	if s.cancelled {
		s.t.Error("store was told to cancel")
	}
}

func TestHandler(t *testing.T) {
	data := "A long time ago in a galaxy far, far away"

	t.Run("testing / path", func(t *testing.T) {
		store := &SpyStore{data, false, t}
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
		store := &SpyStore{data, false, t}
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

func TestServer(t *testing.T) {
	data := "Home, home on the range"

	t.Run("returns data from store", func(t *testing.T) {
		store := &SpyStore{data, false, t}
		s := Server(store)

		req := httptest.NewRequest(http.MethodGet, "/", nil)
		rr := httptest.NewRecorder()

		s.ServeHTTP(rr, req)

		got := rr.Body.String()

		if got != data {
			t.Errorf("got %q, expected %q", got, data)
		}
	})

	t.Run("tells store to cancel work if request is cancelled", func(t *testing.T) {
		store := &SpyStore{data, false, t}
		s := Server(store)

		req := httptest.NewRequest(http.MethodGet, "/", nil)

		cancellingCtx, cancel := context.WithCancel(req.Context())
		time.AfterFunc(5*time.Millisecond, cancel)
		req = req.WithContext(cancellingCtx)

		rr := httptest.NewRecorder()

		s.ServeHTTP(rr, req)

		store.assertWasCancelled()
	})
}
