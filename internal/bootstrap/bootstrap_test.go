package bootstrap

import (
	"testing"

	"github.com/vekio/shorty/internal/app/createlink"
	"github.com/vekio/shorty/internal/app/deletelink"
	"github.com/vekio/shorty/internal/app/getlink"
	"github.com/vekio/shorty/internal/app/listlinks"
	"github.com/vekio/shorty/internal/app/resolvelink"
)

func TestNewWiresHandlersToSharedRepository(t *testing.T) {
	deps := New()
	application := deps.Application
	if deps.Logger == nil {
		t.Fatal("New() returned a nil logger")
	}
	if application.Commands.CreateLink == nil ||
		application.Commands.DeleteLink == nil ||
		application.Commands.ResolveLink == nil ||
		application.Queries.GetLink == nil ||
		application.Queries.ListLinks == nil {
		t.Fatal("New() returned an application with nil handlers")
	}

	created, err := application.Commands.CreateLink.Handle(t.Context(), createlink.CreateLinkCommand{URL: "https://example.com"})
	if err != nil {
		t.Fatalf("create link: %v", err)
	}
	resolved, err := application.Commands.ResolveLink.Handle(t.Context(), resolvelink.ResolveLinkCommand{Code: created.Code})
	if err != nil {
		t.Fatalf("resolve link: %v", err)
	}
	if resolved.OriginURL != "https://example.com" {
		t.Errorf("resolved origin URL = %q, want https://example.com", resolved.OriginURL)
	}
	found, err := application.Queries.GetLink.Handle(t.Context(), getlink.GetLinkQuery{Code: created.Code})
	if err != nil {
		t.Fatalf("get link: %v", err)
	}
	if found.Code != created.Code || found.Visits != 1 {
		t.Errorf("found = %#v, want created link with one visit", found)
	}
	listed, err := application.Queries.ListLinks.Handle(t.Context(), listlinks.ListLinksQuery{})
	if err != nil {
		t.Fatalf("list links: %v", err)
	}
	if len(listed.Links) != 1 || listed.Links[0].Code != created.Code {
		t.Errorf("listed = %#v, want created link", listed)
	}
	if _, err := application.Commands.DeleteLink.Handle(t.Context(), deletelink.DeleteLinkCommand{Code: created.Code}); err != nil {
		t.Fatalf("delete link: %v", err)
	}
	if _, err := application.Queries.GetLink.Handle(t.Context(), getlink.GetLinkQuery{Code: created.Code}); err == nil {
		t.Fatal("get deleted link returned nil error")
	}
}
