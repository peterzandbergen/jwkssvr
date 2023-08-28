package main

import (
	"fmt"
	"net/http"
	"net/url"
	"os"

	"golang.org/x/exp/slog"
)

const (
	JWKSPing = "https://fibi2.acc.belastingdienst.nl/pf/JWKS"
)

var (
	Version = "v0.0.4"
	AppName = "jwkssvr"
)

func buildLogger(opts *options) *slog.Logger {
	logger := getLogger(opts.LogFormat, opts.LogLevel)
	logger = logger.With(
		slog.String("app", AppName),
		slog.String("version", Version),
	)
	return logger
}

func main() {
	opts := NewOptions()

	if opts.Version {
		fmt.Print(Version)
		return
	}

	// Build the logger
	logger := buildLogger(opts)

	// Determine the jwksURI if not set
	if opts.JWKSUri == "" {
		opts.discoverJWKSUri()
		logger.Info(
			"discovered jwks_uri from issuer",
			slog.String("issuer", opts.IssuerURL),
			slog.String("jwks_uri", opts.JWKSUri),
		)
	}

	// Exit if we have no jwks uri.
	if opts.JWKSUri == "" {
		logger.Error("jwksUri is empty", slog.String("use_default", JWKSPing))
		os.Exit(1)
	}
	// if opts.JWKSUri
	if _, err := url.Parse(opts.JWKSUri); err != nil {
		logger.Error("error in jkwsUri", slog.String("err", err.Error()))
		os.Exit(1)
	}
	logger.Info("no error in jkwsUri", slog.String("jwks_uri", opts.JWKSUri))

	if opts.DryRun {
		logger := logger.WithGroup("options")
		logger.Info(
			"dry run",
			slog.String("PORT", opts.Port),
			slog.String("LOG_FORMAT", opts.LogFormat),
			slog.String("LOG_LEVEL", opts.LogLevel),
			slog.String("JWKS_URI", opts.JWKSUri),
			slog.String("ISSUER", opts.IssuerURL),
		)
		return
	}

	// Create server
	svr := newServer()
	svr.jwksURI = opts.JWKSUri
	svr.logger = logger
	svr.routes()

	// Run Server
	logger.Info(
		"starting server",
		slog.String("port", opts.Port),
		slog.String("jwksUri", opts.JWKSUri),
	)

	http.ListenAndServe(":"+opts.Port, svr)
}
