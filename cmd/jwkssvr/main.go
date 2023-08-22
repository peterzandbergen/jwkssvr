package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"

	"golang.org/x/exp/slog"

	"jkwksvr/jwks"
)

const (
	JWKSPing = "https://fibi2.acc.belastingdienst.nl/pf/JWKS"
)

func filterWithAlg(k *jwks.JWK) bool {
	return k.Algorithm != ""
}

func getJWKS(url string, f ...func(j *jwks.JWK) bool) (*jwks.JWKS, error) {
	ff := func(_ *jwks.JWK) bool { return true }
	if len(f) > 1 {
		panic("max one f accepted")
	}
	if len(f) == 1 {
		ff = f[0]
	}

	b, err := getJWKSBytes(url)
	if err != nil {
		return nil, err
	}
	var from jwks.JWKS
	if err := json.Unmarshal(b, &from); err != nil {
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

func getJWKSBytes(url string) ([]byte, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	b, err := io.ReadAll(resp.Body)
	return b, err
}

func newHandleJWKS(l *slog.Logger) http.HandlerFunc {
	l = l.With(slog.String("handler", "HandleJWKS"))
	return func(w http.ResponseWriter, r *http.Request) {
		defer func(start time.Time) {
			l.Info(
				"served request",
				slog.Duration("duration", time.Since(start)),
				slog.String("remote-host", r.RemoteAddr),
				slog.String("path", r.URL.Path),
			)
		}(time.Now())
		startRetrieve := time.Now()
		ks, err := getJWKS(JWKSPing, filterWithAlg)
		endRetrieve := time.Since(startRetrieve)
		l.Info(
			"retrieved jwks",
			slog.Duration("duration", endRetrieve),
		)
		if err != nil {
			l.Error("internal error", slog.String("err", err.Error()))
			http.Error(w,
				fmt.Errorf("error retrieving jwks from %s: %w", JWKSPing, err).Error(),
				http.StatusInternalServerError)
			return
		}
		json.NewEncoder(w).Encode(ks)
	}
}

func main() {
	opts := NewOptions()

	// Build the logger
	logger := getLogger(opts.LogFormat, opts.LogLevel)
	logger = logger.With(slog.String("app", os.Args[0]))

	// Determine the jwksURI if not set
	if opts.RemoteJWKS == "" {
		opts.discoverJWKSUri()
	}

	if opts.DryRun {
		logger := logger.WithGroup("options")
		logger.Info(
			"dry run",
			slog.String("PORT", opts.Port),
			slog.String("LOG_FORMAT", opts.LogFormat),
			slog.String("LOG_LEVEL", opts.LogLevel),
			slog.String("JWKS_URI", opts.RemoteJWKS),
			slog.String("ISSUER", opts.IssuerURL),
		)
		// fmt.Println(opts.String())
		return
	}

	// Run Server
	logger.Info(
		"starting server",
		slog.String("port", opts.Port),
		slog.String("jwksUri", opts.RemoteJWKS),
	)
	http.ListenAndServe(":"+opts.Port, newHandleJWKS(logger))
}
