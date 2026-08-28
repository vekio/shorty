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
	ErrOriginURLRequired = errors.New("origin URL is required")
	ErrOriginURLInvalid  = errors.New("origin URL must be an absolute HTTP or HTTPS URL")
)

// Link is the shortened representation of an absolute HTTP(S) URL.
type Link struct {
	code      string
	originURL url.URL
	createdAt time.Time
	visits    int
}

// New creates a Link with a URL-safe, cryptographically random code.
func New(rawOriginURL string) (Link, error) {
	originURL, err := parseOriginURL(rawOriginURL)
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

// ValidateOriginURL checks whether rawOriginURL can be used as a link target.
func ValidateOriginURL(rawOriginURL string) error {
	_, err := parseOriginURL(rawOriginURL)
	return err
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
