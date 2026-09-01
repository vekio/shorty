package config

import "testing"

func TestLoggerConfigValidate(t *testing.T) {
	for _, level := range []string{"debug", "info", "warn", "error"} {
		if err := (LoggerConfig{Level: level}).Validate(); err != nil {
			t.Errorf("level %q: %v", level, err)
		}
	}
	if err := (LoggerConfig{Level: "verbose"}).Validate(); err == nil {
		t.Fatal("unsupported logger level returned nil error")
	}
	for _, format := range []string{"", LoggerFormatJSON, LoggerFormatText} {
		if err := (LoggerConfig{Format: format, Level: "info"}).Validate(); err != nil {
			t.Errorf("format %q: %v", format, err)
		}
	}
	if err := (LoggerConfig{Format: "pretty", Level: "info"}).Validate(); err == nil {
		t.Fatal("unsupported logger format returned nil error")
	}
}
