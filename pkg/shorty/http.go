package shorty

import (
	"context"
	"encoding/json/v2"
	"fmt"
	"io"
	"net/http"
	"net/url"
)

func (client *Client) request(
	ctx context.Context,
	method string,
	path string,
	body io.Reader,
) (*http.Request, error) {
	reference, err := url.Parse(path)
	if err != nil {
		return nil, fmt.Errorf("build API URL: %w", err)
	}
	endpoint := client.baseURL.ResolveReference(reference)
	request, err := http.NewRequestWithContext(ctx, method, endpoint.String(), body)
	if err != nil {
		return nil, fmt.Errorf("build API request: %w", err)
	}
	request.Header.Set("Accept", "application/json, application/problem+json")
	if ownerID := ownerFromContext(ctx); ownerID != "" {
		request.Header.Set(ownerHeader, ownerID)
	}
	return request, nil
}

func (client *Client) doJSON(request *http.Request, expectedStatus int, destination any) error {
	response, err := client.httpClient.Do(request)
	if err != nil {
		return fmt.Errorf("send API request: %w", err)
	}
	defer response.Body.Close()
	if response.StatusCode != expectedStatus {
		return decodeProblem(response)
	}
	if err := json.UnmarshalRead(response.Body, destination); err != nil {
		return fmt.Errorf("decode API response: %w", err)
	}
	return nil
}

func decodeProblem(response *http.Response) error {
	var problem ProblemError
	if err := json.UnmarshalRead(response.Body, &problem); err != nil {
		return fmt.Errorf("API returned HTTP %d with an invalid problem response: %w", response.StatusCode, err)
	}
	// The transport status is authoritative if a malformed server response
	// disagrees with its JSON body.
	problem.Status = response.StatusCode
	return &problem
}
