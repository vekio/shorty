package pages

import "github.com/vekio/shorty/internal/web/components"

// DashboardData contains the installation details rendered by the admin page.
type DashboardData struct {
	WorkspaceName string
	APIKeys       components.APIKeyPanelData
}
