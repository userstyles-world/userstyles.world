package core

import (
	"github.com/gofiber/fiber/v2"
	"userstyles.world/handlers/jwt"
)

func NotFound(c *fiber.Ctx) error {
	u, _ := jwt.User(c)

	return c.Render("err", fiber.Map{
		"User":  u,
		"Title": "Page not found",
	})
}
