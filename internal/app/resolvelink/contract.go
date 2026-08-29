package resolvelink

// ResolveLinkCommand identifies the shortened link being followed.
type ResolveLinkCommand struct {
	Code string
}

// ResolveLinkResult contains the destination of the shortened link.
type ResolveLinkResult struct {
	OriginURL string
}
