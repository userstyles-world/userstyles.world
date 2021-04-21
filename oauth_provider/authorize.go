package oauth_provider

import (
	"github.com/gofiber/fiber/v2"
	"userstyles.world/handlers/jwt"
)

func AuthorizeGet(c *fiber.Ctx) error {
	u, _ := jwt.User(c)

	return c.Render("authorize", fiber.Map{
		"user": u,
	})
}
