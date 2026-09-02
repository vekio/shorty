package domain

import (
	"errors"
	"strings"
	"testing"
	"time"
)

func TestNewLinkCreatesLink(t *testing.T) {
	before := time.Now().UTC()
	link, err := NewLink("  https://example.com/docs?q=go  ")
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

func TestNewLinkRejectsInvalidOriginURL(t *testing.T) {
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
			_, err := NewLink(test.url)
			if !errors.Is(err, test.want) {
				t.Errorf("New() error = %v, want %v", err, test.want)
			}
		})
	}
}

func TestValidateOriginURL(t *testing.T) {
	if err := ValidateOriginURL("https://example.com/docs"); err != nil {
		t.Errorf("ValidateOriginURL() error = %v", err)
	}
	if err := ValidateOriginURL("not-a-url"); !errors.Is(err, ErrOriginURLInvalid) {
		t.Errorf("ValidateOriginURL() error = %v, want %v", err, ErrOriginURLInvalid)
	}
}

func TestDisallowOriginHostRejectsThisShortyInstance(t *testing.T) {
	policy, err := DisallowOriginHost("https://sho.rt")
	if err != nil {
		t.Fatalf("DisallowOriginHost() error = %v", err)
	}

	for _, originURL := range []string{
		"https://sho.rt",
		"https://SHO.RT/r/already-shortened",
	} {
		_, err := NewLink(originURL, policy)
		if !errors.Is(err, ErrOriginURLSelfReference) {
			t.Errorf("NewLink(%q) error = %v, want %v", originURL, err, ErrOriginURLSelfReference)
		}
	}
}

func TestDisallowOriginHostAllowsAnotherAuthority(t *testing.T) {
	policy, err := DisallowOriginHost("http://localhost:3000")
	if err != nil {
		t.Fatalf("DisallowOriginHost() error = %v", err)
	}

	if _, err := NewLink("http://localhost:8080/docs", policy); err != nil {
		t.Errorf("NewLink() error = %v, want another port to be accepted", err)
	}
}

func TestDisallowOriginHostRejectsInvalidPublicURL(t *testing.T) {
	if _, err := DisallowOriginHost("/relative"); !errors.Is(err, ErrOriginURLInvalid) {
		t.Errorf("DisallowOriginHost() error = %v, want wrapped %v", err, ErrOriginURLInvalid)
	}
}

func TestRegisterVisitMutatesOnlyReceiver(t *testing.T) {
	link, err := NewLink("https://example.com")
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

func TestRestoreLinkRehydratesPersistedState(t *testing.T) {
	createdAt := time.Date(2026, time.September, 1, 9, 30, 0, 0, time.FixedZone("test", 2*60*60))
	link, err := RestoreLink("abc123", "https://example.com/docs", createdAt, 7)
	if err != nil {
		t.Fatalf("RestoreLink() error = %v", err)
	}
	if link.Code() != "abc123" || link.Visits() != 7 || !link.CreatedAt().Equal(createdAt) {
		t.Errorf("RestoreLink() = %#v", link)
	}
	if link.CreatedAt().Location() != time.UTC {
		t.Errorf("CreatedAt() location = %v, want UTC", link.CreatedAt().Location())
	}
}

func TestRestoreLinkRejectsInvalidPersistedState(t *testing.T) {
	for _, test := range []struct {
		name      string
		code      string
		originURL string
		createdAt time.Time
		visits    int
	}{
		{name: "code", originURL: "https://example.com", createdAt: time.Now()},
		{name: "origin URL", code: "abc", originURL: "invalid", createdAt: time.Now()},
		{name: "creation time", code: "abc", originURL: "https://example.com"},
		{name: "visits", code: "abc", originURL: "https://example.com", createdAt: time.Now(), visits: -1},
	} {
		t.Run(test.name, func(t *testing.T) {
			if _, err := RestoreLink(test.code, test.originURL, test.createdAt, test.visits); err == nil {
				t.Fatal("RestoreLink() returned nil error")
			}
		})
	}
}

func TestChangeOriginURLAppliesValidationAndPolicies(t *testing.T) {
	link, err := NewLink("https://example.com/old")
	if err != nil {
		t.Fatalf("NewLink() error = %v", err)
	}
	policy, err := DisallowOriginHost("https://sho.rt")
	if err != nil {
		t.Fatalf("DisallowOriginHost() error = %v", err)
	}
	if err := link.ChangeOriginURL("https://example.com/new", policy); err != nil {
		t.Fatalf("ChangeOriginURL() error = %v", err)
	}
	if originURL := link.OriginURL(); originURL.String() != "https://example.com/new" {
		t.Errorf("OriginURL() = %q", originURL.String())
	}
	if err := link.ChangeOriginURL("https://sho.rt/r/existing", policy); !errors.Is(err, ErrOriginURLSelfReference) {
		t.Errorf("ChangeOriginURL() error = %v, want %v", err, ErrOriginURLSelfReference)
	}
}
