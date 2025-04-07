// Package core provides base endpoints.
package core

import (
	"github.com/gofiber/fiber/v2"

	"userstyles.world/handlers/jwt"
	jwtware "userstyles.world/handlers/jwt"
	"userstyles.world/handlers/middleware"
	"userstyles.world/models"
	"userstyles.world/modules/database"
)

// Routes provides routes for Fiber's router.
func Routes(app *fiber.App) {
	r := app.Group("/")
	r.Get("/", Home)
	r.Get("/proxy", Proxy)
	r.Get("/search", Search)
	r.Get("/docs/*", GetDocs)
	r.Get("/modlog", middleware.Alert, GetModLogs)
	r.Get("/modlog/:id", middleware.Alert, GetModLog)
	r.Get("/link/:site", GetLinkedSite)
	r.Get("/security-policy", Redirect("/docs/security"))
	r.Get("/sitemap.xml", GetSiteMap)
	r.Get("/monitor/*", jwtware.Protected, Monitor)
	r.Get("/dashboard", jwtware.Protected, Dashboard)
	r.Get("/changelog", changelogPage)

	r = app.Group("/changelog", jwt.Admin)
	r.Get("/create", createChangelogPage)
	r.Post("/create", createChangelogForm)

	r = app.Group("/changelog/:id", jwt.Admin, parseID)
	r.Get("/edit", editChangelogPage)
	r.Post("/edit", editChangelogForm)
	r.Get("/delete", deleteChangelogPage)
	r.Post("/delete", deleteChangelogForm)
}

func parseID(c *fiber.Ctx) error {
	i, err := c.ParamsInt("id")
	if err != nil || i < 1 {
		return c.Status(fiber.StatusBadRequest).Render("err", fiber.Map{
			"Title": "ID must be a positive number",
		})
	}

	cl, err := models.GetChangelog(database.Conn, i)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).Render("err", fiber.Map{
			"Title": "Failed to get the changelog",
		})
	}
	c.Locals("Changelog", cl)

	return c.Next()
}
