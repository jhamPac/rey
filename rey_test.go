package rey

import (
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
	t.Run("testing / path", func(t *testing.T) {
		data := "A long time ago in a galaxy far, far away"
		s := Server(&SpyStore{data, false})

		req := httptest.NewRequest(http.MethodGet, "/", nil)
		rr := httptest.NewRecorder()

		s.ServeHTTP(rr, req)

		if rr.Body.String() != data {
			t.Errorf(`got "%s", expected "%s"`, rr.Body.String(), data)
		}
	})

}
