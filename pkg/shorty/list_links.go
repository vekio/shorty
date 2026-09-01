package shorty

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
)

// ListLinks returns one page of links.
func (client *Client) ListLinks(ctx context.Context, options ListOptions) (LinkPage, error) {
	query := make(url.Values)
	if options.Limit != 0 {
		query.Set("limit", fmt.Sprint(options.Limit))
	}
	if options.Offset != 0 {
		query.Set("offset", fmt.Sprint(options.Offset))
	}
	path := "/links"
	if encoded := query.Encode(); encoded != "" {
		path += "?" + encoded
	}
	request, err := client.request(ctx, http.MethodGet, path, nil)
	if err != nil {
		return LinkPage{}, err
	}
	var page LinkPage
	if err := client.doJSON(request, http.StatusOK, &page); err != nil {
		return LinkPage{}, err
	}
	return page, nil
}
