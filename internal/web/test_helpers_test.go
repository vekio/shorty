package web

import (
	"context"
	"io"
	"log/slog"
	"net/http"

	shortysdk "github.com/vekio/shorty/pkg/shorty"
)

type linkAPIStub struct {
	createdURL         string
	createCode         string
	createErr          error
	page               shortysdk.LinkPage
	listOptions        shortysdk.ListOptions
	listPage           func(shortysdk.ListOptions) shortysdk.LinkPage
	listErr            error
	deletedCode        string
	deleteErr          error
	resolvedCode       string
	resolveDestination string
	resolveErr         error
}

func (stub *linkAPIStub) DeleteLink(_ context.Context, code string) error {
	stub.deletedCode = code
	if stub.deleteErr == nil {
		links := stub.page.Links[:0]
		for _, link := range stub.page.Links {
			if link.Code != code {
				links = append(links, link)
			}
		}
		stub.page.Links = links
		stub.page.Total = len(links)
	}
	return stub.deleteErr
}

func (stub *linkAPIStub) CreateLink(_ context.Context, originURL string) (string, error) {
	stub.createdURL = originURL
	return stub.createCode, stub.createErr
}

func (stub *linkAPIStub) ListLinks(_ context.Context, options shortysdk.ListOptions) (shortysdk.LinkPage, error) {
	stub.listOptions = options
	if stub.listPage != nil {
		return stub.listPage(options), stub.listErr
	}
	return stub.page, stub.listErr
}

func (stub *linkAPIStub) ResolveLink(_ context.Context, code string) (string, error) {
	stub.resolvedCode = code
	return stub.resolveDestination, stub.resolveErr
}

func newTestWeb(api *linkAPIStub) http.Handler {
	logger := slog.New(slog.NewTextHandler(io.Discard, nil))
	return New(api, logger, "https://sho.rt")
}
