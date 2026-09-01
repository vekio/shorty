package web

import (
	"net/http"
	"strconv"
)

const (
	defaultLinksPageLimit = 20
	maximumLinksPageLimit = 100
)

func paginationFromRequest(request *http.Request) (int, int) {
	query := request.URL.Query()
	limit := queryInteger(query.Get("limit"), defaultLinksPageLimit)
	if limit < 1 || limit > maximumLinksPageLimit {
		limit = defaultLinksPageLimit
	}
	offset := queryInteger(query.Get("offset"), 0)
	if offset < 0 {
		offset = 0
	}
	return limit, offset
}

func queryInteger(rawValue string, fallback int) int {
	if rawValue == "" {
		return fallback
	}
	value, err := strconv.Atoi(rawValue)
	if err != nil {
		return fallback
	}
	return value
}
