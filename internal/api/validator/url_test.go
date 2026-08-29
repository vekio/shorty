package validator

import (
	"errors"
	"testing"

	"github.com/vekio/shorty/internal/domain"
)

func TestValidateURL(t *testing.T) {
	tests := []struct {
		name    string
		value   string
		wantErr error
	}{
		{name: "empty", value: "  ", wantErr: domain.ErrOriginURLRequired},
		{name: "HTTP", value: "http://example.com"},
		{name: "HTTPS", value: "https://example.com/docs?q=go"},
		{name: "relative", value: "/docs", wantErr: domain.ErrOriginURLInvalid},
		{name: "unsupported scheme", value: "ftp://example.com/file", wantErr: domain.ErrOriginURLInvalid},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			err := validateURL(test.value)
			if !errors.Is(err, test.wantErr) {
				t.Errorf("validateURL(%q) error = %v, want %v", test.value, err, test.wantErr)
			}
		})
	}
}
