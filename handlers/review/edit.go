package review

import (
	"github.com/gofiber/fiber/v2"

	"userstyles.world/handlers/jwt"
)

func updatePage(c *fiber.Ctx) error {
	u, _ := jwt.User(c)
	id := c.Params("id")

	return c.Render("style/review", fiber.Map{
		"Title":    "Update review",
		"User":     u,
		"ID":       nil,
		"ReviewID": id,
	})
}

func updateForm(c *fiber.Ctx) error {
	u, _ := jwt.User(c)
	id := c.Params("id")

	return c.Render("style/review", fiber.Map{
		"Title":    "Update review",
		"User":     u,
		"ID":       nil,
		"ReviewID": id,
	})
}
