// pkg/log/logger.go
package log

import (
	"io"
	"os"
	"time"

	// FIXME: decouple logger from config package
	"bitka/config"
	"github.com/rs/zerolog"
)

// NewLogger returns a pointer to a zerolog.Logger configured for env.
// - development: human-friendly console writer
// - production: JSON to stdout
func NewLogger(cfg *config.Config) *zerolog.Logger {
	zerolog.TimeFieldFormat = time.RFC3339
	var writer io.Writer = os.Stdout

	if cfg.Env == "development" {
		// ConsoleWriter prints pretty output for local dev
		writer = zerolog.ConsoleWriter{Out: os.Stdout, TimeFormat: time.RFC3339}
	}

	logger := zerolog.New(writer).With().
		Str("service", "auth-service").
		Str("env", cfg.Env).
		Timestamp().
		Logger()

	// set global level
	level, err := zerolog.ParseLevel(cfg.LogLevel)
	if err != nil {
		level = zerolog.InfoLevel
	}
	zerolog.SetGlobalLevel(level)

	// return address of logger so callers can call methods on pointer
	return &logger
}
