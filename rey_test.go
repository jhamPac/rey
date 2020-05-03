package rey

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

type SpyStore struct {
	response string
	t        *testing.T
}

func (s *SpyStore) Fetch(ctx context.Context) (string, error) {
	data := make(chan string, 1)

	go func() {
		var result string
		for _, c := range s.response {
			select {
			case <-ctx.Done():
				s.t.Log("spy store got cancelled")
				return
			default:
				time.Sleep(10 * time.Millisecond)
				result += string(c)
			}
		}
		data <- result
	}()

	select {
	case <-ctx.Done():
		return "", ctx.Err()
	case res := <-data:
		return res, nil
	}
}

// func TestHandler(t *testing.T) {
// 	data := "A long time ago in a galaxy far, far away"

// 	t.Run("testing / path", func(t *testing.T) {
// 		store := &SpyStore{data, false, t}
// 		s := Server(store)

// 		req := httptest.NewRequest(http.MethodGet, "/", nil)
// 		rr := httptest.NewRecorder()

// 		s.ServeHTTP(rr, req)

// 		if rr.Body.String() != data {
// 			t.Errorf(`got "%s", expected "%s"`, rr.Body.String(), data)
// 		}

// 		if store.cancelled {
// 			t.Error("store should have not cancelled")
// 		}
// 	})

// 	t.Run("store should cancel work when cancelled", func(t *testing.T) {
// 		store := &SpyStore{data, false, t}
// 		s := Server(store)

// 		req := httptest.NewRequest(http.MethodGet, "/", nil)

// 		cancellingCtx, cancel := context.WithCancel(req.Context())
// 		time.AfterFunc(5*time.Millisecond, cancel)
// 		req = req.WithContext(cancellingCtx)

// 		rr := httptest.NewRecorder()

// 		s.ServeHTTP(rr, req)

// 		if !store.cancelled {
// 			t.Errorf("store was not cancelled")
// 		}
// 	})
// }

func TestServer(t *testing.T) {
	data := "Home, home on the range"

	t.Run("returns data from store", func(t *testing.T) {
		store := &SpyStore{data, t}
		s := Server(store)

		req := httptest.NewRequest(http.MethodGet, "/", nil)
		rr := httptest.NewRecorder()

		s.ServeHTTP(rr, req)

		got := rr.Body.String()

		if got != data {
			t.Errorf("got %q, expected %q", got, data)
		}
	})

	// t.Run("tells store to cancel work if request is cancelled", func(t *testing.T) {
	// 	store := &SpyStore{data, false, t}
	// 	s := Server(store)

	// 	req := httptest.NewRequest(http.MethodGet, "/", nil)

	// 	cancellingCtx, cancel := context.WithCancel(req.Context())
	// 	time.AfterFunc(5*time.Millisecond, cancel)
	// 	req = req.WithContext(cancellingCtx)

	// 	rr := httptest.NewRecorder()

	// 	s.ServeHTTP(rr, req)

	// 	store.assertWasCancelled()
	// })
}
