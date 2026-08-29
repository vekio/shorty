package root

// ResolveLinkInput binds the shortened code from the root path.
type ResolveLinkInput struct {
	Code string `path:"code" json:"-"`
}

// ResolveLinkOutput exposes the redirect destination as an HTTP header.
type ResolveLinkOutput struct {
	Location string `header:"Location" json:"-"`
}
