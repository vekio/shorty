package bootstrap

import "testing"

func TestNewLinkRepository(t *testing.T) {
	if repository := newLinkRepository(); repository == nil {
		t.Fatal("newLinkRepository() returned nil")
	}
}
