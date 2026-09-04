// Package config defines the configuration for the single Shorty process.
package config

import (
	"fmt"
	"strings"

	vekconfig "github.com/vekio/config"
)

const configFileName = "config.yml"

// DatabaseDriver identifies a persistence adapter available to Shorty.
type DatabaseDriver string

const (
	DatabaseDriverMemory DatabaseDriver = "memory"
	DatabaseDriverSQLite DatabaseDriver = "sqlite"
)

// Config contains the HTTP server, public URL, and logger settings.
type Config struct {
	Address  string         `json:"address" yaml:"address"`
	ShortURL string         `json:"short_url" yaml:"short_url"`
	Logger   LoggerConfig   `json:"logger" yaml:"logger"`
	Database DatabaseConfig `json:"database" yaml:"database"`
}

// DatabaseConfig selects and configures Shorty's persistence adapter.
type DatabaseConfig struct {
	Driver DatabaseDriver `json:"driver" yaml:"driver"`
	Path   string         `json:"path" yaml:"path"`
}

// Default returns a development-ready Shorty configuration.
func Default() Config {
	return Config{
		Address:  ":8080",
		ShortURL: "http://localhost:8080",
		Logger:   DefaultLoggerConfig(),
		Database: DatabaseConfig{Driver: DatabaseDriverSQLite, Path: "data/shorty.db"},
	}
}

// Validate implements config.Validatable.
func (cfg Config) Validate() error {
	if err := ValidateListenAddress(cfg.Address); err != nil {
		return fmt.Errorf("address: %w", err)
	}
	if err := ValidateHTTPURL("short URL", cfg.ShortURL); err != nil {
		return err
	}
	if err := cfg.Logger.Validate(); err != nil {
		return fmt.Errorf("logger: %w", err)
	}
	if err := cfg.Database.Validate(); err != nil {
		return fmt.Errorf("database: %w", err)
	}
	return nil
}

// Validate checks that the selected database adapter has the settings it needs.
func (cfg DatabaseConfig) Validate() error {
	switch cfg.Driver {
	case DatabaseDriverMemory:
		return nil
	case DatabaseDriverSQLite:
		if cfg.Path == "" || strings.TrimSpace(cfg.Path) != cfg.Path {
			return fmt.Errorf("path must be non-empty for sqlite")
		}
		return nil
	default:
		return fmt.Errorf("unsupported driver %q: must be memory or sqlite", cfg.Driver)
	}
}

// Load reads the user's configuration or creates it with defaults.
func Load() (Config, error) {
	file, err := New[Config](configFileName)
	if err != nil {
		return Config{}, err
	}
	return file.LoadOrCreate(Default())
}

var _ vekconfig.Validatable = Config{}
