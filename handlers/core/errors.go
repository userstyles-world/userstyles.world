package core

import (
	"strings"

	"github.com/gofiber/fiber/v2"
	"userstyles.world/handlers/jwt"
)

func NotFound(c *fiber.Ctx) error {
	u, _ := jwt.User(c)

	// Handle 404 errors on API route, otherwise render error template.
	if strings.HasPrefix(c.OriginalURL(), "/api/") {
		return c.JSON(fiber.Map{
			"error": "bad endpoint",
		})
	} else {
		return c.Render("err", fiber.Map{
			"User":  u,
			"Title": "Page not found",
		})
	}
}
