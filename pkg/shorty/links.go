package shorty

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

const maximumResponseSize = 1 << 20

// CreateLinkRequest contains the destination for a new short link.
type CreateLinkRequest struct {
	OriginURL string `json:"origin_url"`
}

// CreateLinkResult identifies a newly created short link.
type CreateLinkResult struct {
	Code     string
	Location string
}

type createLinkResponse struct {
	Code string `json:"code"`
}

// CreateLink creates a short link through the management API.
func (client *Client) CreateLink(
	ctx context.Context,
	input CreateLinkRequest,
) (CreateLinkResult, error) {
	body, err := json.Marshal(input)
	if err != nil {
		return CreateLinkResult{}, fmt.Errorf("encode create-link request: %w", err)
	}
	request, err := http.NewRequestWithContext(
		ctx,
		http.MethodPost,
		client.endpoint("/api/v1/links"),
		bytes.NewReader(body),
	)
	if err != nil {
		return CreateLinkResult{}, fmt.Errorf("create create-link request: %w", err)
	}
	request.Header.Set("Accept", "application/json")
	request.Header.Set("Authorization", "Bearer "+client.apiKey)
	request.Header.Set("Content-Type", "application/json")

	response, err := client.httpClient.Do(request)
	if err != nil {
		return CreateLinkResult{}, fmt.Errorf("create link: %w", err)
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusCreated {
		return CreateLinkResult{}, decodeAPIError(response)
	}

	var output createLinkResponse
	if err := decodeResponse(response.Body, &output); err != nil {
		return CreateLinkResult{}, fmt.Errorf("decode create-link response: %w", err)
	}
	if output.Code == "" {
		return CreateLinkResult{}, fmt.Errorf("decode create-link response: code is missing")
	}
	return CreateLinkResult{
		Code:     output.Code,
		Location: response.Header.Get("Location"),
	}, nil
}

func decodeAPIError(response *http.Response) error {
	problem := Problem{
		Title:  http.StatusText(response.StatusCode),
		Status: response.StatusCode,
		Detail: http.StatusText(response.StatusCode),
	}
	if err := decodeResponse(response.Body, &problem); err != nil {
		return &APIError{StatusCode: response.StatusCode, Problem: problem}
	}
	if problem.Status == 0 {
		problem.Status = response.StatusCode
	}
	return &APIError{StatusCode: response.StatusCode, Problem: problem}
}

func decodeResponse(body io.Reader, destination any) error {
	return json.NewDecoder(io.LimitReader(body, maximumResponseSize)).Decode(destination)
}
