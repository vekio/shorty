package shorty

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
)

// DeleteLink deletes a link by code.
func (client *Client) DeleteLink(ctx context.Context, code string) error {
	request, err := client.request(ctx, http.MethodDelete, "/links/"+url.PathEscape(code), nil)
	if err != nil {
		return err
	}
	response, err := client.httpClient.Do(request)
	if err != nil {
		return fmt.Errorf("send delete-link request: %w", err)
	}
	defer response.Body.Close()
	if response.StatusCode != http.StatusNoContent {
		return decodeProblem(response)
	}
	return nil
}
