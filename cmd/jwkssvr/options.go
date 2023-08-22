package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"net/http"
	"os"
	"reflect"
	"strings"

	"golang.org/x/exp/slog"
)

const (
	EnvPort      = "PORT"
	EnvIssuer    = "ISSUER"
	EnvJwksUri   = "JWKS_URI"
	EnvLogFormat = "LOG_FORMAT"
	EnvLogLevel  = "LOG_LEVEL"

	FlagPort      = "port"
	FlagIssuer    = "issuer"
	FlagJwksUri   = "jwks-uri"
	FlagLogFormat = "log-format"
	FlagLogLevel  = "log-level"
	FlagDryRun    = "dry-run"

	DiscoverySuffix = ".well-known/openid-configuration"
)

type options struct {
	Port       string // PORT, default 8080
	IssuerURL  string // ISSUER, default empty
	RemoteJWKS string // JWKS_URI, default https://fibi2.acc.belastingdienst.nl/pf/JWKS
	LogLevel   string // LOG_LEVEL, default info
	LogFormat  string // LOG_FORMAT, default text
	DryRun     bool
}

func (o *options) parseFlags(args []string) {
	fs := flag.NewFlagSet("default", flag.ContinueOnError)
	port := fs.String(FlagPort, "", "port number")
	jwksUri := fs.String(FlagJwksUri, "", "remote jwks uri")
	logLevel := fs.String(FlagLogLevel, "", "log level [info | warn | error | debug]")
	logFormat := fs.String(FlagLogFormat, "", "log format [ text | json ]")
	issuer := fs.String("issuer", "", "issuer, when a valid URL this will be used to discover the jwksUri")
	fs.BoolVar(&o.DryRun, FlagDryRun, false, "")
	fs.Parse(args)
	if *port != "" {
		o.Port = *port
	}
	if *issuer != "" {
		o.IssuerURL = *issuer
	}
	if *jwksUri != "" {
		o.RemoteJWKS = *jwksUri
	}
	if *logLevel != "" {
		o.LogLevel = *logLevel
	}
	if *logFormat != "" {
		o.LogFormat = *logFormat
	}
}

func getEnv(key, def string) string {
	res := def
	if val := os.Getenv(key); val != "" {
		res = val
	}
	return res
}

func getLoglevel(level string) slog.Level {
	switch strings.ToLower(level) {
	case "debug":
		return slog.LevelDebug
	case "warn":
		return slog.LevelWarn
	case "error":
		return slog.LevelError
	default:
		return slog.LevelInfo
	}
}

func getLogger(format string, level string) *slog.Logger {
	var lh slog.Handler
	switch strings.ToLower(format) {
	case "json":
		lh = slog.NewJSONHandler(os.Stderr, &slog.HandlerOptions{})
	default:
		lh = slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{
			Level: getLoglevel(level),
		})
	}
	return slog.New(lh)
}

func (o *options) getFromEnv() {
	o.Port = getEnv(EnvPort, "8080")
	o.RemoteJWKS = getEnv(EnvJwksUri, "")
	o.LogLevel = getEnv(EnvLogLevel, "info")
	o.LogFormat = getEnv(EnvLogFormat, "text")
	o.IssuerURL = getEnv(EnvIssuer, "")
}

func NewOptions() *options {
	o := &options{}
	o.getFromEnv()
	o.parseFlags(os.Args[1:])
	return o
}

func (o *options) String() string {
	return fmt.Sprintf("JWKS_URI=%s LOG_LEVEL=%s LOG_FORMAT=%s dry-run=%t PORT=%s", o.RemoteJWKS, o.LogLevel, o.LogFormat, o.DryRun, o.Port)
}

func (o *options) discoverJWKSUri() error {
	if o.RemoteJWKS != "" ||  o.IssuerURL == "" {
		return nil
	}

	url := o.IssuerURL + "/" + DiscoverySuffix
	resp, err := http.Get(url)
	if err != nil {
		return err
	}

	var discovery map[string]any
	if err := json.NewDecoder(resp.Body).Decode(&discovery); err != nil {
		return err
	}
	// Find the jwksUri entry
	{
		e := discovery["jwks_uri"]
		if e == nil {
			return fmt.Errorf("cannot find jwks_uri entry in response from %s", url)
		}
		es, ok := e.(string)
		if !ok {
			return fmt.Errorf("jwks_uri entry is not a string: %s", reflect.TypeOf(e).String())
		}
		o.RemoteJWKS = es
		return nil
	}
}