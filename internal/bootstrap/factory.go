package bootstrap

import (
	"fmt"
	"io"
	"log/slog"
	"os"

	"github.com/vekio/shorty/internal/app/ports"
	shortyconfig "github.com/vekio/shorty/internal/config"
	"github.com/vekio/shorty/internal/infra/memory"
)

// NewLogger creates the process logger using its configured output format.
func NewLogger(config shortyconfig.LoggerConfig) (*slog.Logger, error) {
	return newLogger(config, os.Stdout)
}

func newLogger(config shortyconfig.LoggerConfig, output io.Writer) (*slog.Logger, error) {
	var level slog.Level
	if err := level.UnmarshalText([]byte(config.Level)); err != nil {
		return nil, err
	}
	options := &slog.HandlerOptions{
		Level:     level,
		AddSource: config.AddSource,
	}
	var handler slog.Handler
	switch config.Format {
	case "", shortyconfig.LoggerFormatJSON:
		handler = slog.NewJSONHandler(output, options)
	case shortyconfig.LoggerFormatText:
		handler = slog.NewTextHandler(output, options)
	default:
		return nil, fmt.Errorf("invalid logger format %q", config.Format)
	}
	return slog.New(handler), nil
}

// NewLinkRepository creates the persistence adapter used by the API process.
func NewLinkRepository() ports.LinkRepository {
	return memory.NewLinkRepository()
}
