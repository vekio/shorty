package api

import "github.com/vekio/shorty/internal/app"

type handlers struct {
	createLink app.CreateLinkHandler
	getLink    app.GetLinkHandler
	listLinks  app.ListLinksHandler
	visitLink  app.VisitLinkHandler
}

func newHandlers(application app.Application) *handlers {
	return &handlers{
		createLink: application.Commands.CreateLink,
		getLink:    application.Queries.GetLink,
		listLinks:  application.Queries.ListLinks,
		visitLink:  application.Commands.VisitLink,
	}
}
