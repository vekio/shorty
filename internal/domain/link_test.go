package domain

import (
	"errors"
	"strings"
	"testing"
	"time"
)

func TestNewCreatesLink(t *testing.T) {
	before := time.Now().UTC()
	link, err := New("  https://example.com/docs?q=go  ")
	after := time.Now().UTC()
	if err != nil {
		t.Fatalf("New() error = %v", err)
	}
	originURL := link.OriginURL()
	if got := originURL.String(); got != "https://example.com/docs?q=go" {
		t.Errorf("OriginURL() = %q", got)
	}
	if got := link.Code(); len(got) != 12 {
		t.Errorf("Code() length = %d, want 12", len(got))
	}
	if strings.ContainsAny(link.Code(), "+/=") {
		t.Errorf("Code() = %q, want unpadded URL-safe base64", link.Code())
	}
	if link.CreatedAt().Before(before) || link.CreatedAt().After(after) {
		t.Errorf("CreatedAt() = %v, want between %v and %v", link.CreatedAt(), before, after)
	}
	if link.CreatedAt().Location() != time.UTC {
		t.Errorf("CreatedAt() location = %v, want UTC", link.CreatedAt().Location())
	}
	if link.Visits() != 0 {
		t.Errorf("Visits() = %d, want 0", link.Visits())
	}
}

func TestNewRejectsInvalidOriginURL(t *testing.T) {
	tests := []struct {
		name string
		url  string
		want error
	}{
		{name: "empty", url: "  ", want: ErrOriginURLRequired},
		{name: "relative", url: "/docs", want: ErrOriginURLInvalid},
		{name: "missing host", url: "https:///docs", want: ErrOriginURLInvalid},
		{name: "unsupported scheme", url: "ftp://example.com/file", want: ErrOriginURLInvalid},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			_, err := New(test.url)
			if !errors.Is(err, test.want) {
				t.Errorf("New() error = %v, want %v", err, test.want)
			}
		})
	}
}

func TestRegisterVisitMutatesOnlyReceiver(t *testing.T) {
	link, err := New("https://example.com")
	if err != nil {
		t.Fatalf("New() error = %v", err)
	}
	copy := link
	copy.RegisterVisit()
	if copy.Visits() != 1 {
		t.Errorf("copy Visits() = %d, want 1", copy.Visits())
	}
	if link.Visits() != 0 {
		t.Errorf("original Visits() = %d, want 0", link.Visits())
	}
}
