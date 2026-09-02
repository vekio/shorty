package getlink

import "time"

type GetLinkQuery struct {
	Code string
}

type GetLinkResult struct {
	Code      string
	OriginURL string
	CreatedAt time.Time
	Visits    int
}
