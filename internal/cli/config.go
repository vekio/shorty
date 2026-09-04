package cli

import (
	"fmt"
	"strings"

	shortyconfig "github.com/vekio/shorty/internal/config"
)

// Config contains the connection settings used by the Shorty CLI.
type Config struct {
	ServerURL string `json:"server_url" yaml:"server_url"`
	APIKey    string `json:"api_key,omitempty" yaml:"api_key,omitempty"`
}

// DefaultConfig returns settings for a locally running Shorty server.
func DefaultConfig() Config {
	return Config{ServerURL: "http://localhost:8080"}
}

// Validate verifies the CLI connection settings.
func (config Config) Validate() error {
	if err := shortyconfig.ValidateHTTPURL("server URL", config.ServerURL); err != nil {
		return err
	}
	if strings.TrimSpace(config.APIKey) != config.APIKey {
		return fmt.Errorf("API key cannot contain surrounding whitespace")
	}
	return nil
}
