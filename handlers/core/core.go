// Package core provides base endpoints.
package core

import (
	"github.com/gofiber/fiber/v2"

	jwtware "userstyles.world/handlers/jwt"
)

// Routes provides routes for Fiber's router.
func Routes(app *fiber.App) {
	r := app.Group("/")
	r.Get("/", Home)
	r.Get("/proxy", Proxy)
	r.Get("/search", Search)
	r.Get("/docs/:document?", GetDocs)
	r.Get("/modlog", GetModLog)
	r.Get("/link/:site", GetLinkedSite)
	r.Get("/security-policy", Redirect("/docs/security"))
	r.Get("/sitemap.xml", GetSiteMap)
	r.Get("/monitor/*", jwtware.Protected, Monitor)
	r.Get("/dashboard", jwtware.Protected, Dashboard)
}
