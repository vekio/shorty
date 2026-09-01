package bootstrap

import (
	"bytes"
	"context"
	"log/slog"
	"strings"
	"testing"

	shortyconfig "github.com/vekio/shorty/internal/config"
)

func TestNewLoggerUsesConfiguredLevel(t *testing.T) {
	logger, err := NewLogger(shortyconfig.LoggerConfig{Level: "warn"})
	if err != nil {
		t.Fatalf("NewLogger() error = %v", err)
	}
	if logger.Enabled(context.Background(), slog.LevelInfo) {
		t.Error("info logging is enabled at warn level")
	}
	if !logger.Enabled(context.Background(), slog.LevelWarn) {
		t.Error("warn logging is disabled at warn level")
	}
}

func TestNewLoggerUsesConfiguredFormat(t *testing.T) {
	tests := []struct {
		format string
		want   string
	}{
		{format: shortyconfig.LoggerFormatJSON, want: `"msg":"hello"`},
		{format: shortyconfig.LoggerFormatText, want: "msg=hello"},
	}

	for _, test := range tests {
		t.Run(test.format, func(t *testing.T) {
			var output bytes.Buffer
			logger, err := newLogger(shortyconfig.LoggerConfig{
				Format: test.format,
				Level:  "info",
			}, &output)
			if err != nil {
				t.Fatalf("newLogger() error = %v", err)
			}
			logger.Info("hello")
			if !strings.Contains(output.String(), test.want) {
				t.Errorf("output = %q, want it to contain %q", output.String(), test.want)
			}
		})
	}
}

func TestNewLoggerRejectsInvalidLevel(t *testing.T) {
	if _, err := NewLogger(shortyconfig.LoggerConfig{Level: "verbose"}); err == nil {
		t.Fatal("NewLogger() returned nil error")
	}
}

func TestNewLoggerRejectsInvalidFormat(t *testing.T) {
	if _, err := NewLogger(shortyconfig.LoggerConfig{Format: "pretty", Level: "info"}); err == nil {
		t.Fatal("NewLogger() returned nil error")
	}
}

func TestNewLinkRepository(t *testing.T) {
	if NewLinkRepository() == nil {
		t.Fatal("NewLinkRepository() returned nil")
	}
}
