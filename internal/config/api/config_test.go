package api

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	shortyconfig "github.com/vekio/shorty/internal/config"
)

func TestConfigValidate(t *testing.T) {
	if err := Default().Validate(); err != nil {
		t.Fatalf("default configuration is invalid: %v", err)
	}
	if err := (Config{}).Validate(); err == nil {
		t.Fatal("empty configuration returned nil error")
	}
}

func TestConfigValidateIdentifiesInvalidComponent(t *testing.T) {
	tests := []struct {
		name   string
		config Config
		prefix string
	}{
		{name: "address", config: Config{ShortURL: "https://sho.rt", Logger: shortyconfig.DefaultLoggerConfig()}, prefix: "address: "},
		{name: "short URL", config: Config{Address: ":8080", ShortURL: "/shorty", Logger: shortyconfig.DefaultLoggerConfig()}, prefix: "short URL "},
		{name: "logger", config: Config{Address: ":8080", ShortURL: "https://sho.rt", Logger: shortyconfig.LoggerConfig{Level: "verbose"}}, prefix: "logger: "},
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
		[]byte("address: 127.0.0.1:9080\nshort_url: https://sho.rt\nlogger:\n  format: text\n  level: debug\n  add_source: true\n"),
		0o600,
	); err != nil {
		t.Fatalf("write configuration: %v", err)
	}

	got, err := Load()
	if err != nil {
		t.Fatalf("Load() error = %v", err)
	}
	want := Config{
		Address:  "127.0.0.1:9080",
		ShortURL: "https://sho.rt",
		Logger: shortyconfig.LoggerConfig{
			Format:    shortyconfig.LoggerFormatText,
			Level:     "debug",
			AddSource: true,
		},
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
