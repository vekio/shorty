package cli

import (
	"bytes"
	"context"
	"strings"
	"testing"

	vekconfig "github.com/vekio/config"
	"github.com/vekio/shorty/pkg/shorty"
)

type linkCreatorStub struct {
	input shorty.CreateLinkRequest
}

func (stub *linkCreatorStub) CreateLink(
	_ context.Context,
	input shorty.CreateLinkRequest,
) (shorty.CreateLinkResult, error) {
	stub.input = input
	return shorty.CreateLinkResult{Code: "abc123"}, nil
}

func TestNewCreatesShortyRootWithConfigCommand(t *testing.T) {
	t.Setenv("XDG_CONFIG_HOME", t.TempDir())
	command, err := New()
	if err != nil {
		t.Fatalf("New() error = %v", err)
	}

	var output bytes.Buffer
	command.Writer = &output
	if err := command.Run(context.Background(), []string{"shorty", "--help"}); err != nil {
		t.Fatalf("Run() error = %v", err)
	}

	if command.Name != "shorty" {
		t.Errorf("Name = %q", command.Name)
	}
	for _, expected := range []string{"Manage links on a Shorty server", "setup", "config", "links", "--config"} {
		if !strings.Contains(output.String(), expected) {
			t.Errorf("help does not contain %q:\n%s", expected, output.String())
		}
	}
}

func TestSetupCreatesConfigurationAndRequestsAPIKey(t *testing.T) {
	configFile, err := vekconfig.NewYAMLConfigFile[Config]("shorty-cli-test", "config.yml")
	if err != nil {
		t.Fatal(err)
	}
	configPath := t.TempDir() + "/config.yml"
	if err := configFile.SetPath(configPath); err != nil {
		t.Fatal(err)
	}

	command := newRootCommand(configFile, DefaultConfig(), func(Config) (linkCreator, error) {
		t.Fatal("setup created an API client")
		return nil, nil
	})
	var output bytes.Buffer
	command.Writer = &output
	if err := command.Run(t.Context(), []string{"shorty", "setup"}); err != nil {
		t.Fatalf("Run() error = %v", err)
	}

	config, err := configFile.Load()
	if err != nil {
		t.Fatalf("Load() error = %v", err)
	}
	if config != DefaultConfig() {
		t.Errorf("config = %#v, want %#v", config, DefaultConfig())
	}
	if !strings.Contains(output.String(), configPath) ||
		!strings.Contains(output.String(), "set api_key") {
		t.Errorf("output = %q", output.String())
	}
}

func TestCreateLinkCommandUsesCLIConfiguration(t *testing.T) {
	configFile, err := vekconfig.NewYAMLConfigFile[Config]("shorty-cli-test", "config.yml")
	if err != nil {
		t.Fatal(err)
	}
	configPath := t.TempDir() + "/config.yml"
	if err := configFile.SetPath(configPath); err != nil {
		t.Fatal(err)
	}
	wantConfig := Config{ServerURL: "https://shorty.example", APIKey: "shorty_token"}
	if err := configFile.Save(wantConfig); err != nil {
		t.Fatal(err)
	}

	creator := &linkCreatorStub{}
	var gotConfig Config
	command := newRootCommand(configFile, DefaultConfig(), func(config Config) (linkCreator, error) {
		gotConfig = config
		return creator, nil
	})
	var output bytes.Buffer
	command.Writer = &output
	if err := command.Run(t.Context(), []string{
		"shorty", "--config", configPath, "links", "create", "https://example.com/docs",
	}); err != nil {
		t.Fatalf("Run() error = %v", err)
	}

	if got := output.String(); got != "abc123\n" {
		t.Errorf("output = %q", got)
	}
	if gotConfig != wantConfig {
		t.Errorf("client config = %#v, want %#v", gotConfig, wantConfig)
	}
	if creator.input.OriginURL != "https://example.com/docs" {
		t.Errorf("create input = %#v", creator.input)
	}
}
