package main

import (
	"encoding/json"
	"fmt"
	"io"
	"jkwksvr/jwks"
	"net/http"
	"time"

	"golang.org/x/exp/slog"
)

const (
	HeaderContentType = "Content-Type"

	ContentTypeJWKS = "application/json;charset=utf-8"
)

type jwksFilterFunc func(k *jwks.JWK) bool

func filterNone(_ *jwks.JWK) bool { return true }
func filterWithAlg(k *jwks.JWK) bool {
	return k.Algorithm != ""
}

func getJWKS(url string, f ...jwksFilterFunc) (*jwks.JWKS, error) {
	ff := filterNone
	if len(f) > 1 {
		panic("max one f accepted")
	}
	if len(f) == 1 {
		ff = f[0]
	}

	b, err := getBody(url)
	if err != nil {
		return nil, err
	}
	var from jwks.JWKS
	if err := json.NewDecoder(b).Decode(&from); err != nil {
		return nil, err
	}
	var to jwks.JWKS
	for _, j := range from.Keys {
		if ff(&j) {
			to.Keys = append(to.Keys, j)
		}
	}
	return &to, nil
}

type server struct {
	*http.ServeMux
	jwksURI string
	logger  *slog.Logger
}

func newServer() *server {
	return &server{
		ServeMux: http.NewServeMux(),
	}
}

func (s *server) routes() {
	s.HandleFunc("/", s.newHandleJWKS(filterWithAlg))
	s.HandleFunc("/fixed", s.newHandleJWKS(filterWithAlg))
	s.HandleFunc("/raw", s.newHandleJWKS(filterNone))

}

func (s *server) logWithHeaders(r *http.Request) *slog.Logger {
	var vals []any
	for k, h := range r.Header {
		if k == "Authorization" {
			continue
		}
		vals = append(vals, slog.String(k, h[0]))
	}
	return s.logger.WithGroup("headers").With(vals...)
}

func (s *server) newHandleJWKS(filter jwksFilterFunc) http.HandlerFunc {
	l := s.logger.With(slog.String("handler", "HandleJWKS"))
	return func(w http.ResponseWriter, r *http.Request) {
		l = s.logWithHeaders(r)

		defer func(start time.Time) {
			l.Info(
				"served request",
				slog.Duration("duration", time.Since(start)),
				slog.String("remote-host", r.RemoteAddr),
				slog.String("path", r.URL.Path),
			)
		}(time.Now())

		// Only serve GET
		if r.Method != http.MethodGet {
			status := http.StatusMethodNotAllowed
			http.Error(w, http.StatusText(status), status)
			return
		}

		startRetrieve := time.Now()
		ks, err := getJWKS(s.jwksURI, filter)
		endRetrieve := time.Since(startRetrieve)
		l.Info(
			"retrieved jwks",
			slog.Duration("duration", endRetrieve),
		)
		if err != nil {
			l.Error("internal error", slog.String("err", err.Error()))
			http.Error(w,
				fmt.Errorf("error retrieving jwks from %s: %w", s.jwksURI, err).Error(),
				http.StatusInternalServerError)
			return
		}
		w.Header().Set(HeaderContentType, ContentTypeJWKS)
		json.NewEncoder(w).Encode(ks)
	}
}

func getBody(url string) (io.Reader, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	if err != nil {
		return nil, err
	}
	return resp.Body, nil
}
