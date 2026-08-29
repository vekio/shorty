package deletelink

// DeleteLinkCommand identifies the link to remove.
type DeleteLinkCommand struct {
	Code string
}

// DeleteLinkResult represents a successful deletion without response data.
type DeleteLinkResult struct{}
