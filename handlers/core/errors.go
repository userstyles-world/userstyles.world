package core

import (
	"github.com/gofiber/fiber/v2"
	"userstyles.world/handlers/sessions"
)

func NotFound(c *fiber.Ctx) error {
	s := sessions.State(c)

	return c.Render("404", fiber.Map{
		"Name":  s.Get("name"),
		"Title": "Page not found",
	})
}
