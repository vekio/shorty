package shorty

import (
	"context"
	"net/http"
	"net/url"
)

// ResolveLink registers a visit and returns its destination representation.
func (client *Client) ResolveLink(ctx context.Context, code string) (string, error) {
	request, err := client.request(
		ctx,
		http.MethodPost,
		"/links/"+url.PathEscape(code)+"/resolve",
		nil,
	)
	if err != nil {
		return "", err
	}
	var output struct {
		OriginURL string `json:"origin_url"`
	}
	if err := client.doJSON(request, http.StatusOK, &output); err != nil {
		return "", err
	}
	return output.OriginURL, nil
}

// RedirectURL returns the API navigation endpoint that redirects a code.
func (client *Client) RedirectURL(code string) string {
	path := "/r/" + url.PathEscape(code)
	reference, _ := url.Parse(path)
	return client.baseURL.ResolveReference(reference).String()
}
