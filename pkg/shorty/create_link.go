package shorty

import (
	"bytes"
	"context"
	"encoding/json/v2"
	"fmt"
	"net/http"
)

// CreateLink creates a shortened URL and returns its code.
func (client *Client) CreateLink(ctx context.Context, originURL string) (string, error) {
	body, err := json.Marshal(struct {
		OriginURL string `json:"origin_url"`
	}{OriginURL: originURL})
	if err != nil {
		return "", fmt.Errorf("encode create-link request: %w", err)
	}

	request, err := client.request(ctx, http.MethodPost, "/links", bytes.NewReader(body))
	if err != nil {
		return "", err
	}
	request.Header.Set("Content-Type", "application/json")

	var output struct {
		Code string `json:"code"`
	}
	if err := client.doJSON(request, http.StatusCreated, &output); err != nil {
		return "", err
	}
	return output.Code, nil
}
