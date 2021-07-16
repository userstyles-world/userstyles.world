package style

import (
	"github.com/gofiber/fiber/v2"

	"userstyles.world/handlers/jwt"
)

func ReviewGet(c *fiber.Ctx) error {
	u, _ := jwt.User(c)
	id := c.Params("id")

	return c.Render("style/review", fiber.Map{
		"Title": "Review style",
		"User":  u,
		"ID":    id,
	})
}

func ReviewPost(c *fiber.Ctx) error {
	return nil
}
