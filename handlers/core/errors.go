package core

import (
	"github.com/gofiber/fiber/v2"
	"userstyles.world/handlers/sessions"
)

func NotFound(c *fiber.Ctx) error {
	u := sessions.User(c)

	return c.Render("err", fiber.Map{
		"User":  u,
		"Title": "Page not found",
	})
}
