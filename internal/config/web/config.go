// Package web defines configuration owned by the HTML Web process.
package web

import (
	"fmt"

	vekconfig "github.com/vekio/config"
	shortyconfig "github.com/vekio/shorty/internal/config"
)

const (
	configFileName = "config-web.yml"
	configFileEnv  = "SHORTY_WEB_CONFIG_FILE"
)

// Config contains the Web process, public URL, and upstream API settings.
type Config struct {
	Address  string                    `json:"address" yaml:"address"`
	ShortURL string                    `json:"short_url" yaml:"short_url"`
	APIURL   string                    `json:"api_url" yaml:"api_url"`
	Logger   shortyconfig.LoggerConfig `json:"logger" yaml:"logger"`
}

// Default returns a development-ready Web configuration.
func Default() Config {
	return Config{
		Address:  ":3000",
		ShortURL: "http://localhost:3000",
		APIURL:   "http://localhost:8080",
		Logger:   shortyconfig.DefaultLoggerConfig(),
	}
}

// Validate implements config.Validatable.
func (cfg Config) Validate() error {
	if err := shortyconfig.ValidateListenAddress(cfg.Address); err != nil {
		return fmt.Errorf("address: %w", err)
	}
	if err := shortyconfig.ValidateHTTPURL("short URL", cfg.ShortURL); err != nil {
		return err
	}
	if err := shortyconfig.ValidateHTTPURL("API URL", cfg.APIURL); err != nil {
		return err
	}
	if err := cfg.Logger.Validate(); err != nil {
		return fmt.Errorf("logger: %w", err)
	}
	return nil
}

// Load reads the user's Web configuration or creates it with defaults.
func Load() (Config, error) {
	file, err := shortyconfig.New[Config](configFileName)
	if err != nil {
		return Config{}, err
	}
	if err := shortyconfig.SetPathFromEnv(file, configFileEnv); err != nil {
		return Config{}, err
	}
	return file.LoadOrCreate(Default())
}

var _ vekconfig.Validatable = Config{}
