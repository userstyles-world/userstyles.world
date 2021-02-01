package core

import (
	"github.com/gofiber/fiber/v2"
)

func NotFound(c *fiber.Ctx) error {
	return c.Render("404", fiber.Map{
		"Title": "Page not found",
	})
}
