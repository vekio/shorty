package updatelink

import "time"

type UpdateLinkCommand struct {
	Code      string
	OriginURL string
}

type UpdateLinkResult struct {
	Code      string
	OriginURL string
	CreatedAt time.Time
	Visits    int
}
