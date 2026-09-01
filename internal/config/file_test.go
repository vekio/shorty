package config

import (
	"path/filepath"
	"strings"
	"testing"

	vekconfig "github.com/vekio/config"
)

type fileTestConfig struct{}

func (fileTestConfig) Validate() error { return nil }

func TestSetPathFromEnv(t *testing.T) {
	file, err := vekconfig.NewYAMLConfigFile[fileTestConfig]("shorty", "config.yml")
	if err != nil {
		t.Fatalf("NewYAMLConfigFile() error = %v", err)
	}
	want := filepath.Join(t.TempDir(), "local.yml")
	t.Setenv("SHORTY_TEST_CONFIG_FILE", want)

	if err := SetPathFromEnv(file, "SHORTY_TEST_CONFIG_FILE"); err != nil {
		t.Fatalf("SetPathFromEnv() error = %v", err)
	}
	if file.Path() != want {
		t.Errorf("Path() = %q, want %q", file.Path(), want)
	}
}

func TestSetPathFromEnvIgnoresEmptyVariable(t *testing.T) {
	file, err := vekconfig.NewYAMLConfigFile[fileTestConfig]("shorty", "config.yml")
	if err != nil {
		t.Fatalf("NewYAMLConfigFile() error = %v", err)
	}
	want := file.Path()
	t.Setenv("SHORTY_TEST_CONFIG_FILE", "")

	if err := SetPathFromEnv(file, "SHORTY_TEST_CONFIG_FILE"); err != nil {
		t.Fatalf("SetPathFromEnv() error = %v", err)
	}
	if file.Path() != want {
		t.Errorf("Path() = %q, want conventional path %q", file.Path(), want)
	}
}

func TestSetPathFromEnvReportsInvalidPath(t *testing.T) {
	file, err := vekconfig.NewYAMLConfigFile[fileTestConfig]("shorty", "config.yml")
	if err != nil {
		t.Fatalf("NewYAMLConfigFile() error = %v", err)
	}
	t.Setenv("SHORTY_TEST_CONFIG_FILE", " ")

	err = SetPathFromEnv(file, "SHORTY_TEST_CONFIG_FILE")
	if err == nil || !strings.Contains(err.Error(), "SHORTY_TEST_CONFIG_FILE") {
		t.Errorf("SetPathFromEnv() error = %v, want environment variable name", err)
	}
}

var _ vekconfig.Validatable = fileTestConfig{}
