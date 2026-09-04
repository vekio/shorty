package config

import (
	"os"
	"path/filepath"
	"testing"
)

func TestDefaultIsValid(t *testing.T) {
	if err := Default().Validate(); err != nil {
		t.Fatalf("default configuration is invalid: %v", err)
	}
}

func TestConfigValidateRejectsInvalidFields(t *testing.T) {
	tests := []struct {
		name   string
		config Config
	}{
		{name: "address", config: Config{ShortURL: "https://sho.rt", Logger: DefaultLoggerConfig(), Database: DatabaseConfig{Driver: DatabaseDriverMemory}}},
		{name: "short URL", config: Config{Address: ":8080", ShortURL: "/shorty", Logger: DefaultLoggerConfig(), Database: DatabaseConfig{Driver: DatabaseDriverMemory}}},
		{name: "logger", config: Config{Address: ":8080", ShortURL: "https://sho.rt", Logger: LoggerConfig{Level: "verbose"}, Database: DatabaseConfig{Driver: DatabaseDriverMemory}}},
		{name: "database driver", config: Config{Address: ":8080", ShortURL: "https://sho.rt", Logger: DefaultLoggerConfig(), Database: DatabaseConfig{Driver: "postgres"}}},
		{name: "sqlite path", config: Config{Address: ":8080", ShortURL: "https://sho.rt", Logger: DefaultLoggerConfig(), Database: DatabaseConfig{Driver: DatabaseDriverSQLite}}},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if err := test.config.Validate(); err == nil {
				t.Fatal("Validate() returned nil error")
			}
		})
	}
}

func TestLoadFromConfiguredFile(t *testing.T) {
	path := filepath.Join(t.TempDir(), configFileName)
	contents := []byte("address: ':9090'\nshort_url: https://sho.rt\nlogger:\n  format: text\n  level: debug\ndatabase:\n  driver: memory\n")
	if err := os.WriteFile(path, contents, 0o600); err != nil {
		t.Fatalf("write configuration: %v", err)
	}
	t.Setenv("SHORTY_CONFIG_FILE", path)

	config, err := Load()
	if err != nil {
		t.Fatalf("Load() error = %v", err)
	}
	if config.Address != ":9090" || config.ShortURL != "https://sho.rt" || config.Logger.Format != LoggerFormatText || config.Database.Driver != DatabaseDriverMemory {
		t.Errorf("Load() = %#v", config)
	}
}
