// Package web composes the server-rendered Web process.
package web

import (
	"net/http"
	"time"

	"github.com/vekio/shorty/internal/bootstrap"
	webconfig "github.com/vekio/shorty/internal/config/web"
	"github.com/vekio/shorty/internal/httpmiddleware"
	webadapter "github.com/vekio/shorty/internal/web"
	shortysdk "github.com/vekio/shorty/pkg/shorty"
)

// New composes the complete HTML Web process and its Shorty API client.
func New(config webconfig.Config) (bootstrap.Runtime, error) {
	processLogger, err := bootstrap.NewLogger(config.Logger)
	if err != nil {
		return bootstrap.Runtime{}, err
	}
	client, err := shortysdk.NewClient(config.APIURL, &http.Client{Timeout: 5 * time.Second})
	if err != nil {
		return bootstrap.Runtime{}, err
	}
	handler := webadapter.New(client, processLogger, config.ShortURL)

	return bootstrap.Runtime{
		Handler: httpmiddleware.LogRequests(processLogger)(handler),
		Logger:  processLogger,
	}, nil
}
