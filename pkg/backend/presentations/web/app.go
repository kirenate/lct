package web

import (
	"github.com/gofiber/fiber/v2"
	recover2 "github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/teadove/teasutils/fiber_utils"
	"main.go/services"
	"time"
)

type Presentation struct {
	service *services.Service
}

func NewPresentation(service *services.Service) *Presentation {
	return &Presentation{service: service}
}

func (r *Presentation) BuildApp() *fiber.App {
	app := fiber.New(fiber.Config{
		Immutable:    true,
		ErrorHandler: fiber_utils.ErrHandler(),
		BodyLimit:    10000000000,
	})
	app.Use(fiber_utils.MiddlewareLogger())
	app.Use(recover2.New(recover2.Config{EnableStackTrace: true}))
	app.Use(fiber_utils.MiddlewareCtxTimeout(29 * time.Second))

	app.Post("/api/documents/:documentId/pages", r.getDocumentPages)

	app.Get("/api/documents/:documentId", r.getSingleDocument)
	app.Get("/api/documents", r.getDocuments)
	app.Post("/api/documents", r.uploadDocument)

	app.Delete("/api/:documentId", r.deleteDocument)

	app.Get("/openapi.yaml", r.openapi)
	app.Get("/api/docs", r.swagger)

	return app
}
