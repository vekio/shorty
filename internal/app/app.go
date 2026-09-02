// Package app exposes the application's use-case handlers to input adapters.
package app

// Application groups the exposed use-case handlers. It contains no transport
// or infrastructure logic; CLI and HTTP may consume the same instance.
type Application struct {
	Commands Commands
	Queries  Queries
}

// Commands contains the application's state-changing handlers.
type Commands struct {
	CreateLink  CreateLinkHandler
	UpdateLink  UpdateLinkHandler
	DeleteLink  DeleteLinkHandler
	ResolveLink ResolveLinkHandler
}

// Queries contains the application's read-only handlers.
type Queries struct {
	GetLink   GetLinkHandler
	ListLinks ListLinksHandler
}
