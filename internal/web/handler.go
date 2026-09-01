// Package web exposes Shorty's server-rendered HTML interface.
package web

import (
	"context"
	"net/url"
	"strings"

	shortysdk "github.com/vekio/shorty/pkg/shorty"
)

type shortyAPI interface {
	CreateLink(context.Context, string) (string, error)
	DeleteLink(context.Context, string) error
	ListLinks(context.Context, shortysdk.ListOptions) (shortysdk.LinkPage, error)
	ResolveLink(context.Context, string) (string, error)
}

func apiContext(ctx context.Context) context.Context {
	return shortysdk.WithOwner(ctx, browserSessionFromContext(ctx))
}

type handler struct {
	api            shortyAPI
	renderer       renderer
	shortURLPrefix string
}

func newHandler(api shortyAPI, shortURL string) *handler {
	return &handler{
		api:            api,
		renderer:       newRenderer(),
		shortURLPrefix: strings.TrimRight(shortURL, "/"),
	}
}

func (h *handler) shortURL(code string) string {
	if code == "" {
		return ""
	}
	return h.shortURLPrefix + "/r/" + url.PathEscape(code)
}
