package listlinks

import "time"

type ListLinksQuery struct{}

type ListLinksResult struct {
	Links []LinkResult
}

type LinkResult struct {
	Code      string
	OriginURL string
	CreatedAt time.Time
	Visits    int
}
