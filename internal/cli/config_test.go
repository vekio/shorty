package cli

import "testing"

func TestDefaultConfigIsValid(t *testing.T) {
	if err := DefaultConfig().Validate(); err != nil {
		t.Fatalf("default configuration is invalid: %v", err)
	}
}

func TestConfigValidateRejectsInvalidValues(t *testing.T) {
	tests := []struct {
		name   string
		config Config
	}{
		{name: "server URL", config: Config{ServerURL: "/api/v1"}},
		{name: "API key whitespace", config: Config{ServerURL: "https://shorty.example", APIKey: " shorty_token"}},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if err := test.config.Validate(); err == nil {
				t.Fatal("Validate() returned nil error")
			}
		})
	}
}
