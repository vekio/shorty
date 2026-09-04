package components

import "github.com/vekio/shorty/internal/auth"

// APIKeyPanelData contains the API key state rendered as a full-page section
// or as an HTMX fragment.
type APIKeyPanelData struct {
	APIKeys  []auth.APIKey
	NewToken string
	Error    string
}
