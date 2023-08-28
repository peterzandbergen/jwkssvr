package main

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"golang.org/x/exp/slog"
)

func TestGetOnly(t *testing.T) {

	cases := []struct {
		name           string
		method         string
		expectedStatus int
	}{
		{
			name:           "GET",
			method:         http.MethodGet,
			expectedStatus: http.StatusInternalServerError,
		},
		{
			name:           "HEAD",
			method:         http.MethodHead,
			expectedStatus: http.StatusMethodNotAllowed,
		},
	}

	svr := newServer()
	svr.logger = slog.Default()
	svr.routes()

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			r := httptest.NewRequest(c.method, "/", nil)
			w := httptest.NewRecorder()
			svr.ServeHTTP(w, r)
			if w.Code != c.expectedStatus {
				t.Errorf("got %d, expected %d", w.Code, c.expectedStatus)
			}
		})
	}
}
