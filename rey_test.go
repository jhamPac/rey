package rey

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

type StubStore struct {
	response string
}

func (s *StubStore) Fetch() string {
	return s.response
}

func TestHandler(t *testing.T) {
	data := "A long time ago in a galaxy far, far away"
	s := Server(&StubStore{data})

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rr := httptest.NewRecorder()

	s.ServeHTTP(rr, req)

	if rr.Body.String() != data {
		t.Errorf(`got "%s", expected "%s"`, rr.Body.String(), data)
	}

}
