package main

import (
	"context"
	"errors"
	"fmt"
	"net"
	"net/http"
	"net/url"
	"os"
	"os/signal"
	"time"

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
		if err := opts.discoverJWKSUri(); err != nil {
			logger.Error("error discovering jwks_uri", slog.String("error", err.Error()))
			os.Exit(1)
		}
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
	if _, err := url.Parse(opts.JWKSUri); err != nil {
		logger.Error("error in jkwsUri", slog.String("err", err.Error()))
		os.Exit(1)
	}
	logger.Info("no error in jkwsUri", slog.String("jwks_uri", opts.JWKSUri))

	// Log the settings
	logger.WithGroup("options").Info(
		"using these options",
		slog.String("PORT", opts.Port),
		slog.String("LOG_FORMAT", opts.LogFormat),
		slog.String("LOG_LEVEL", opts.LogLevel),
		slog.String("JWKS_URI", opts.JWKSUri),
		slog.String("ISSUER", opts.IssuerURL),
	)

	if opts.DryRun {
		return
	}

	// Create server
	svr := newServer()
	svr.jwksURI = opts.JWKSUri
	svr.logger = logger
	svr.routes()

	// Set up signals.
	sigCtx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	server := &http.Server{
		Addr:         ":" + opts.Port,
		Handler:      svr,
		ReadTimeout:  1 * time.Second,
		WriteTimeout: 1 * time.Second,
		IdleTimeout:  2 * time.Second,
		BaseContext:  func(_ net.Listener) context.Context { return sigCtx },
	}

	// Start the server
	go func() {
		// Signal the wait that we stopped.
		defer cancel()
		// Run Server
		logger.Info(
			"starting server in go routine",
			slog.String("port", opts.Port),
			slog.String("jwksUri", opts.JWKSUri),
		)
		if err := server.ListenAndServe(); err != nil && ! errors.Is(err, http.ErrServerClosed) {
			logger.Error("server stopped with error", slog.String("err", err.Error()))
			return
		}
		logger.Debug("ListenAndServe stopped")
	}()

	logger.Debug("waiting for done")
	<-sigCtx.Done()
	logger.Debug("received done")
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()
	if err := server.Shutdown(ctx); err != nil {
		logger.Error("server shutdown error", slog.String("error", err.Error()))
	}
	logger.Info("server shutdown")
}
