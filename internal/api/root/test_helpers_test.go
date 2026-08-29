package root

import (
	"context"
	"io"
	"log/slog"

	"github.com/vekio/amigo"
	"github.com/vekio/shorty/internal/app"
	"github.com/vekio/shorty/internal/app/resolvelink"
)

type resolveLinkHandlerStub struct {
	code   string
	result resolvelink.ResolveLinkResult
	err    error
}

func (stub *resolveLinkHandlerStub) Handle(
	_ context.Context,
	command resolvelink.ResolveLinkCommand,
) (resolvelink.ResolveLinkResult, error) {
	stub.code = command.Code
	return stub.result, stub.err
}

func newTestAPI(handler *resolveLinkHandlerStub) *amigo.API {
	logger := slog.New(slog.NewTextHandler(io.Discard, nil))
	httpAPI := amigo.New(amigo.WithLogger(logger))
	Register(httpAPI, app.Application{
		Commands: app.Commands{ResolveLink: handler},
	})
	return httpAPI
}
