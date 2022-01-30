// Package api provides API endpoints.
package api

import (
	"github.com/gofiber/fiber/v2"
)

// Routes provides routes for Fiber's router.
func Routes(app *fiber.App) {
	r := app.Group("/api", ParseAPIJWT)
	r.Head("/style/:id.user.css", GetStyleEtag)
	r.Get("/style/:id.user.css", GetStyleSource)
	r.Get("/style/:id", GetStyleDetails)
	r.Get("/style/preview/:id", GetPreviewScreenshot)
	r.Get("/style/stats/:id/:type?", GetStyleStats)
	r.Get("/index/:format?", GetStyleIndex)
	r.Get("/search/:query", GetSearchResult)
	r.Get("/callback/:rcode", CallbackGet)
	r.Get("/user", ProtectedAPI, UserGet)
	r.Get("/user/:identifier", SpecificUserGet)
	r.Get("/styles", ProtectedAPI, StylesGet)
	r.Post("/style/new", ProtectedAPI, NewStyle)
	r.Post("/style/:id", ProtectedAPI, StylePost)
	r.Delete("/style/:id", ProtectedAPI, DeleteStyle)
	r.Get("/style", ProtectedAPI, StyleGet)
}
