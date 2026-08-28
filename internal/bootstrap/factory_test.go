package bootstrap

import "testing"

func TestNewLogger(t *testing.T) {
	if logger := newLogger(); logger == nil {
		t.Fatal("newLogger() returned nil")
	}
}

func TestNewLinkRepository(t *testing.T) {
	if repository := newLinkRepository(); repository == nil {
		t.Fatal("newLinkRepository() returned nil")
	}
}
