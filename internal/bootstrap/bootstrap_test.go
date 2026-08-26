package bootstrap

import (
	"testing"

	"github.com/vekio/shorty/internal/app/createlink"
	"github.com/vekio/shorty/internal/app/getlink"
	"github.com/vekio/shorty/internal/app/visitlink"
)

func TestNewWiresHandlersToSharedRepository(t *testing.T) {
	application := New()
	if application.Commands.CreateLink == nil || application.Commands.VisitLink == nil || application.Queries.GetLink == nil {
		t.Fatal("New() returned an application with nil handlers")
	}

	created, err := application.Commands.CreateLink.Handle(t.Context(), createlink.CreateLinkCommand{URL: "https://example.com"})
	if err != nil {
		t.Fatalf("create link: %v", err)
	}
	if _, err := application.Commands.VisitLink.Handle(t.Context(), visitlink.VisitLinkCommand{Code: created.Code}); err != nil {
		t.Fatalf("visit link: %v", err)
	}
	found, err := application.Queries.GetLink.Handle(t.Context(), getlink.GetLinkQuery{Code: created.Code})
	if err != nil {
		t.Fatalf("get link: %v", err)
	}
	if found.Code != created.Code || found.Visits != 1 {
		t.Errorf("found = %#v, want created link with one visit", found)
	}
}
