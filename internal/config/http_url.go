package config

import (
	"fmt"
	"net/url"
	"strings"
)

// ValidateHTTPURL checks a public base URL used by a Shorty process.
func ValidateHTTPURL(name string, rawURL string) error {
	if rawURL == "" || strings.TrimSpace(rawURL) != rawURL {
		return fmt.Errorf("%s must be a non-empty absolute URL", name)
	}
	parsed, err := url.Parse(rawURL)
	if err != nil || parsed.Host == "" || (parsed.Scheme != "http" && parsed.Scheme != "https") {
		return fmt.Errorf("%s must be an absolute HTTP or HTTPS URL", name)
	}
	if parsed.User != nil || parsed.RawQuery != "" || parsed.Fragment != "" {
		return fmt.Errorf("%s cannot contain user information, query, or fragment", name)
	}
	if parsed.Path != "" && parsed.Path != "/" {
		return fmt.Errorf("%s cannot contain a path", name)
	}
	return nil
}
