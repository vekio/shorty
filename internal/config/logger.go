package config

import (
	"fmt"
	"log/slog"
)

const (
	LoggerFormatJSON = "json"
	LoggerFormatText = "text"
)

type LoggerConfig struct {
	Format    string `json:"format" yaml:"format"`
	Level     string `json:"level" yaml:"level"`
	AddSource bool   `json:"add_source" yaml:"add_source"`
}

func DefaultLoggerConfig() LoggerConfig {
	return LoggerConfig{Format: LoggerFormatJSON, Level: "info"}
}

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
