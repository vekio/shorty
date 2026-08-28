package api

import "time"

type Empty struct{}

type LinkByIDInput struct {
	Code string `path:"id" json:"-" validate:"required"`
}

type CreateLinkInput struct {
	OriginURL string `json:"origin_url" validate:"required,url"`
}

type CreateLinkOutput struct {
	Location string `header:"Location" json:"-"`
	Code     string `json:"code"`
}

type LinkOutput struct {
	Code      string    `json:"code"`
	OriginURL string    `json:"origin_url"`
	CreatedAt time.Time `json:"created_at"`
	Visits    int       `json:"visits"`
}

type ListLinksOutput struct {
	Links []LinkOutput `json:"links"`
}
