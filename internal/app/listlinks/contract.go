package listlinks

import (
	"errors"
	"time"
)

const (
	// DefaultLimit is used when callers do not specify a page size.
	DefaultLimit = 20
	// MaximumLimit protects repositories from unbounded page requests.
	MaximumLimit = 100
)

var (
	// ErrInvalidLimit indicates that the requested page size is unsupported.
	ErrInvalidLimit = errors.New("limit must be between 1 and 100")
	// ErrInvalidOffset indicates that the page starts before the collection.
	ErrInvalidOffset = errors.New("offset must be non-negative")
)

// ListLinksQuery selects one page from the link collection.
type ListLinksQuery struct {
	Limit  int
	Offset int
}

// ListLinksResult contains one page and its pagination metadata.
type ListLinksResult struct {
	Links  []LinkResult
	Total  int
	Limit  int
	Offset int
}

// LinkResult is the application representation of one shortened link.
type LinkResult struct {
	Code      string
	OriginURL string
	CreatedAt time.Time
	Visits    int
}

// ValidateLimit checks whether limit is an accepted page size.
func ValidateLimit(limit int) error {
	if limit < 1 || limit > MaximumLimit {
		return ErrInvalidLimit
	}
	return nil
}

// ValidateOffset checks whether offset is a valid collection position.
func ValidateOffset(offset int) error {
	if offset < 0 {
		return ErrInvalidOffset
	}
	return nil
}
