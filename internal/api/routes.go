package api

import (
	"net/http"

	"github.com/vekio/shorty/internal/app"
	"github.com/vekio/shorty/pkg/amigo"
)

// New builds the HTTP API and registers all routes.
func New(application app.Application) *amigo.App {
	httpAPI := amigo.New()
	httpAPI.SetMaxBodyBytes(1 << 20)
	httpAPI.MapErrors(mapApplicationError)
	handlers := newHandlers(application)

	httpAPI.POST("/links", handlers.CreateLink,
		amigo.Status(http.StatusCreated),
	)
	httpAPI.GET("/links/{id}", handlers.GetLink)
	httpAPI.POST("/links/{id}/visit", handlers.VisitLink,
		amigo.Status(http.StatusNoContent),
	)
	return httpAPI
}
