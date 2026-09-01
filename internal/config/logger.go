package config

import (
	"fmt"
	"log/slog"
)

const (
	LoggerFormatJSON = "json"
	LoggerFormatText = "text"
)

// LoggerConfig controls the output format, minimum level, and source details.
type LoggerConfig struct {
	Format    string `json:"format" yaml:"format"`
	Level     string `json:"level" yaml:"level"`
	AddSource bool   `json:"add_source" yaml:"add_source"`
}

// DefaultLoggerConfig returns the default production-friendly logger settings.
func DefaultLoggerConfig() LoggerConfig {
	return LoggerConfig{Format: LoggerFormatJSON, Level: "info"}
}

// Validate checks that Level is understood by slog.
func (config LoggerConfig) Validate() error {
	if config.Format != "" && config.Format != LoggerFormatJSON && config.Format != LoggerFormatText {
		return fmt.Errorf("invalid logger format %q: must be json or text", config.Format)
	}
	var level slog.Level
	if err := level.UnmarshalText([]byte(config.Level)); err != nil {
		return fmt.Errorf("invalid logger level %q: %w", config.Level, err)
	}
	return nil
}
