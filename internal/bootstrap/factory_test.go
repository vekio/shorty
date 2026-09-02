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
			logger, err := newLogger(shortyconfig.LoggerConfig{Format: test.format, Level: "info"}, &output)
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

func TestOpenPersistenceSelectsMemory(t *testing.T) {
	persistence, err := OpenPersistence(t.Context(), shortyconfig.DatabaseConfig{Driver: shortyconfig.DatabaseDriverMemory})
	if err != nil {
		t.Fatalf("OpenPersistence() error = %v", err)
	}
	if persistence.WorkspaceName != "default" || !strings.HasPrefix(persistence.WorkspaceID, "ws_") || persistence.Links == nil || persistence.APIKeys == nil {
		t.Fatalf("OpenPersistence() = %#v", persistence)
	}
	if err := persistence.Close(); err != nil {
		t.Errorf("Close() error = %v", err)
	}
}
