// Package core provides base endpoints.
package core

import (
	"github.com/gofiber/fiber/v2"

	"userstyles.world/handlers/jwt"
	jwtware "userstyles.world/handlers/jwt"
	"userstyles.world/handlers/middleware"
)

// Routes provides routes for Fiber's router.
func Routes(app *fiber.App) {
	r := app.Group("/")
	r.Get("/", Home)
	r.Get("/proxy", Proxy)
	r.Get("/search", Search)
	r.Get("/docs/*", GetDocs)
	r.Get("/modlog", middleware.Alert, GetModLog)
	r.Get("/link/:site", GetLinkedSite)
	r.Get("/security-policy", Redirect("/docs/security"))
	r.Get("/sitemap.xml", GetSiteMap)
	r.Get("/monitor/*", jwtware.Protected, Monitor)
	r.Get("/dashboard", jwtware.Protected, Dashboard)
	r.Get("/changelog", changelogPage)

	r = app.Group("/changelog", jwt.Admin)
	r.Get("/create", createChangelogPage)
	r.Post("/create", createChangelogForm)
	r.Get("/:id/edit", editChangelogPage)
	r.Post("/:id/edit", editChangelogForm)
	r.Get("/:id/delete", deleteChangelogPage)
	r.Post("/:id/delete", deleteChangelogForm)
}
