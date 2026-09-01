package shorty

import (
	"context"
	"net/http"
	"net/url"
)

// GetLink returns a link without registering a visit.
func (client *Client) GetLink(ctx context.Context, code string) (Link, error) {
	request, err := client.request(ctx, http.MethodGet, "/links/"+url.PathEscape(code), nil)
	if err != nil {
		return Link{}, err
	}
	var link Link
	if err := client.doJSON(request, http.StatusOK, &link); err != nil {
		return Link{}, err
	}
	return link, nil
}
