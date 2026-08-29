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

// ListLinksInput binds optional pagination parameters from the query string.
type ListLinksInput struct {
	Limit  int `query:"limit" json:"-" validate:"page_limit"`
	Offset int `query:"offset" json:"-" validate:"page_offset"`
}

// ListLinksOutput contains all links returned by the collection endpoint.
type ListLinksOutput struct {
	Links  []LinkOutput `json:"links"`
	Total  int          `json:"total"`
	Limit  int          `json:"limit"`
	Offset int          `json:"offset"`
}
