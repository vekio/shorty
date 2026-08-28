package links

import "time"

// Empty represents an endpoint without request or response fields.
type Empty struct{}

// LinkByIDInput binds the link code from the route path.
type LinkByIDInput struct {
	Code string `path:"id" json:"-"`
}

// CreateLinkInput is the JSON representation accepted when creating a link.
type CreateLinkInput struct {
	OriginURL string `json:"origin_url" validate:"required,url"`
}

// CreateLinkOutput returns the new code and advertises its resource location.
type CreateLinkOutput struct {
	Location string `header:"Location" json:"-"`
	Code     string `json:"code"`
}

// LinkOutput is the public representation of a shortened link.
type LinkOutput struct {
	Code      string    `json:"code"`
	OriginURL string    `json:"origin_url"`
	CreatedAt time.Time `json:"created_at"`
	Visits    int       `json:"visits"`
}

// ListLinksOutput contains all links returned by the collection endpoint.
type ListLinksOutput struct {
	Links []LinkOutput `json:"links"`
}
