package domain

import (
	"crypto/rand"
	"encoding/base64"
	"errors"
	"fmt"
	"net/url"
	"strings"
	"time"
)

const codeBytes = 9

var (
	ErrOriginURLRequired      = errors.New("origin URL is required")
	ErrOriginURLInvalid       = errors.New("origin URL must be an absolute HTTP or HTTPS URL")
	ErrOriginURLSelfReference = errors.New("origin URL cannot point to this Shorty instance")
	ErrLinkCodeRequired       = errors.New("link code is required")
	ErrLinkCreatedAtRequired  = errors.New("link creation time is required")
	ErrLinkVisitsInvalid      = errors.New("link visits cannot be negative")
)

// OriginURLPolicy applies a configurable business rule to a parsed origin URL.
type OriginURLPolicy func(url.URL) error

// Link is the shortened representation of an absolute HTTP(S) URL.
type Link struct {
	code      string
	originURL url.URL
	createdAt time.Time
	visits    int
}

// RestoreLink rebuilds a persisted Link while preserving its identity and state.
func RestoreLink(code string, rawOriginURL string, createdAt time.Time, visits int) (Link, error) {
	if code == "" {
		return Link{}, ErrLinkCodeRequired
	}
	originURL, err := parseOriginURL(rawOriginURL)
	if err != nil {
		return Link{}, err
	}
	if createdAt.IsZero() {
		return Link{}, ErrLinkCreatedAtRequired
	}
	if visits < 0 {
		return Link{}, ErrLinkVisitsInvalid
	}
	return Link{
		code:      code,
		originURL: originURL,
		createdAt: createdAt.UTC(),
		visits:    visits,
	}, nil
}

// NewLink creates a Link with a URL-safe, cryptographically random code.
func NewLink(rawOriginURL string, policies ...OriginURLPolicy) (Link, error) {
	originURL, err := validateOriginURL(rawOriginURL, policies...)
	if err != nil {
		return Link{}, err
	}

	code, err := generateCode()
	if err != nil {
		return Link{}, err
	}

	return Link{
		code:      code,
		originURL: originURL,
		createdAt: time.Now().UTC(),
	}, nil
}

// DisallowOriginHost creates a policy that rejects links targeting the public
// authority used by this Shorty instance, regardless of their path.
func DisallowOriginHost(rawPublicURL string) (OriginURLPolicy, error) {
	publicURL, err := parseOriginURL(rawPublicURL)
	if err != nil {
		return nil, fmt.Errorf("parse Shorty public URL: %w", err)
	}
	return func(originURL url.URL) error {
		if strings.EqualFold(originURL.Host, publicURL.Host) {
			return ErrOriginURLSelfReference
		}
		return nil
	}, nil
}

// ValidateOriginURL checks whether rawOriginURL can be used as a link target.
func ValidateOriginURL(rawOriginURL string) error {
	_, err := validateOriginURL(rawOriginURL)
	return err
}

func validateOriginURL(rawOriginURL string, policies ...OriginURLPolicy) (url.URL, error) {
	originURL, err := parseOriginURL(rawOriginURL)
	if err != nil {
		return url.URL{}, err
	}
	for _, policy := range policies {
		if policy == nil {
			continue
		}
		if err := policy(originURL); err != nil {
			return url.URL{}, err
		}
	}
	return originURL, nil
}

func parseOriginURL(rawURL string) (url.URL, error) {
	rawURL = strings.TrimSpace(rawURL)
	if rawURL == "" {
		return url.URL{}, ErrOriginURLRequired
	}

	parsed, err := url.ParseRequestURI(rawURL)
	if err != nil || parsed.Host == "" || (parsed.Scheme != "http" && parsed.Scheme != "https") {
		return url.URL{}, ErrOriginURLInvalid
	}

	return *parsed, nil
}

func generateCode() (string, error) {
	random := make([]byte, codeBytes)
	if _, err := rand.Read(random); err != nil {
		return "", fmt.Errorf("generate link code: %w", err)
	}

	return base64.RawURLEncoding.EncodeToString(random), nil
}

func (l Link) Code() string {
	return l.code
}

func (l Link) OriginURL() url.URL {
	return l.originURL
}

func (l Link) CreatedAt() time.Time {
	return l.createdAt
}

func (l Link) Visits() int {
	return l.visits
}

func (l *Link) RegisterVisit() {
	l.visits++
}

// ChangeOriginURL replaces the destination after applying the same rules used
// during link creation.
func (l *Link) ChangeOriginURL(rawOriginURL string, policies ...OriginURLPolicy) error {
	originURL, err := validateOriginURL(rawOriginURL, policies...)
	if err != nil {
		return err
	}
	l.originURL = originURL
	return nil
}
