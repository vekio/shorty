package config

import "testing"

func TestValidateListenAddress(t *testing.T) {
	tests := []struct {
		name    string
		address string
		wantErr bool
	}{
		{name: "all interfaces", address: ":8080"},
		{name: "hostname", address: "localhost:3000"},
		{name: "IPv6", address: "[::1]:8080"},
		{name: "empty", wantErr: true},
		{name: "missing port", address: "localhost", wantErr: true},
		{name: "invalid port", address: ":invalid", wantErr: true},
		{name: "port too high", address: ":65536", wantErr: true},
		{name: "surrounding whitespace", address: " :8080", wantErr: true},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			err := ValidateListenAddress(test.address)
			if (err != nil) != test.wantErr {
				t.Errorf("ValidateListenAddress(%q) error = %v, wantErr %v", test.address, err, test.wantErr)
			}
		})
	}
}
