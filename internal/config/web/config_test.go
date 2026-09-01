package web

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	shortyconfig "github.com/vekio/shorty/internal/config"
)

func TestConfigValidate(t *testing.T) {
	tests := []struct {
		name    string
		config  Config
		wantErr bool
	}{
		{name: "default", config: Default()},
		{name: "HTTPS", config: newConfig(":3000", "https://sho.rt", "https://api.sho.rt")},
		{name: "empty address", config: newConfig("", "https://sho.rt", "https://api.sho.rt"), wantErr: true},
		{name: "relative short URL", config: newConfig(":3000", "/shorty", "https://api.sho.rt"), wantErr: true},
		{name: "unsupported scheme", config: newConfig(":3000", "ftp://sho.rt", "https://api.sho.rt"), wantErr: true},
		{name: "short URL query", config: newConfig(":3000", "https://sho.rt?q=go", "https://api.sho.rt"), wantErr: true},
		{name: "short URL path", config: newConfig(":3000", "https://sho.rt/r", "https://api.sho.rt"), wantErr: true},
		{name: "relative API URL", config: newConfig(":3000", "https://sho.rt", "/api"), wantErr: true},
		{name: "API URL path", config: newConfig(":3000", "https://sho.rt", "https://api.sho.rt/v1"), wantErr: true},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			err := test.config.Validate()
			if (err != nil) != test.wantErr {
				t.Errorf("Validate() error = %v, wantErr %v", err, test.wantErr)
			}
		})
	}
}

func TestConfigValidateIdentifiesInvalidSharedComponent(t *testing.T) {
	tests := []struct {
		name   string
		config Config
		prefix string
	}{
		{name: "address", config: newConfig("", "https://sho.rt", "https://api.sho.rt"), prefix: "address: "},
		{name: "logger", config: Config{
			Address: ":3000", ShortURL: "https://sho.rt", APIURL: "https://api.sho.rt",
			Logger: shortyconfig.LoggerConfig{Level: "verbose"},
		}, prefix: "logger: "},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			err := test.config.Validate()
			if err == nil || !strings.HasPrefix(err.Error(), test.prefix) {
				t.Errorf("Validate() error = %v, want prefix %q", err, test.prefix)
			}
		})
	}
}

func TestLoadCreatesDefaultConfiguration(t *testing.T) {
	configDir := t.TempDir()
	t.Setenv("XDG_CONFIG_HOME", configDir)
	t.Setenv(configFileEnv, "")

	got, err := Load()
	if err != nil {
		t.Fatalf("Load() error = %v", err)
	}
	if got != Default() {
		t.Errorf("Load() = %#v, want %#v", got, Default())
	}
	path := filepath.Join(configDir, "shorty", configFileName)
	if _, err := os.Stat(path); err != nil {
		t.Errorf("configuration file %q: %v", path, err)
	}
}

func TestLoadReadsExistingConfiguration(t *testing.T) {
	configDir := t.TempDir()
	t.Setenv("XDG_CONFIG_HOME", configDir)
	t.Setenv(configFileEnv, "")
	directory := filepath.Join(configDir, "shorty")
	if err := os.MkdirAll(directory, 0o700); err != nil {
		t.Fatalf("create configuration directory: %v", err)
	}
	if err := os.WriteFile(
		filepath.Join(directory, configFileName),
		[]byte("address: 127.0.0.1:4000\nshort_url: https://sho.rt\napi_url: https://api.sho.rt\nlogger:\n  format: text\n  level: debug\n  add_source: true\n"),
		0o600,
	); err != nil {
		t.Fatalf("write configuration: %v", err)
	}

	got, err := Load()
	if err != nil {
		t.Fatalf("Load() error = %v", err)
	}
	want := newConfig("127.0.0.1:4000", "https://sho.rt", "https://api.sho.rt")
	want.Logger = shortyconfig.LoggerConfig{
		Format:    shortyconfig.LoggerFormatText,
		Level:     "debug",
		AddSource: true,
	}
	if got != want {
		t.Errorf("Load() = %#v, want %#v", got, want)
	}
}

func TestLoadUsesConfiguredFilePath(t *testing.T) {
	path := filepath.Join(t.TempDir(), configFileName)
	t.Setenv(configFileEnv, path)

	got, err := Load()
	if err != nil {
		t.Fatalf("Load() error = %v", err)
	}
	if got != Default() {
		t.Errorf("Load() = %#v, want %#v", got, Default())
	}
	if _, err := os.Stat(path); err != nil {
		t.Errorf("configuration file %q: %v", path, err)
	}
}

func newConfig(address string, shortURL string, apiURL string) Config {
	return Config{
		Address:  address,
		ShortURL: shortURL,
		APIURL:   apiURL,
		Logger:   shortyconfig.DefaultLoggerConfig(),
	}
}
