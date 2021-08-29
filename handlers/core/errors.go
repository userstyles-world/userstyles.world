package core

import (
	"strings"

	"github.com/userstyles-world/fiber/v2"

	"userstyles.world/handlers/jwt"
)

func NotFound(c *fiber.Ctx) error {
	c.Status(404)
	// Handle 404 errors on API route, otherwise render error template.
	if strings.HasPrefix(c.OriginalURL(), "/api/") {
		return c.JSON(fiber.Map{
			"error": "bad endpoint",
		})
	}
	u, _ := jwt.User(c)
	return c.Render("err", fiber.Map{
		"User":      u,
		"Title":     "Page not found",
		"Canonical": "404",
	})
}
