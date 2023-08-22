package server

import (
	"jkwksvr/jwks"
	"log/slog"
	"net/http"
)

type Server struct {
	*http.Server
	mux *http.ServeMux

	Logger     *slog.Logger
	RemoteJwks string
	Filter     func(jwks *jwks.JWKS) bool

	// JWKS cache
	cache *JWKSCache
}

func New() *Server {
	s := &Server{
		Server: &http.Server{},
		mux:    http.NewServeMux(),
	}
	s.routes()
	return s
}


func (s *Server) handleJWKS() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		// Set the headers.
		// Write the serialized and modified JWKS list.
		if _, err := w.Write(s.cache.Bytes()); err != nil {
			s.Logger.ErrorContext(ctx, "error writing content", "err", err.Error())
		}
		http.Error(w, "not implemented yet", http.StatusNotImplemented)
	}
}

func (s *Server) routes() {
	s.mux.HandleFunc("/", s.handleJWKS())
}

func (s *Server) Handler() http.Handler {
	return s.Server.Handler
}
