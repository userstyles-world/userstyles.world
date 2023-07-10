package review

import (
	"github.com/gofiber/fiber/v2"

	"userstyles.world/handlers/jwt"
)

func deletePage(c *fiber.Ctx) error {
	u, _ := jwt.User(c)
	id := c.Params("id")

	return c.Render("style/deletereview", fiber.Map{
		"Title":    "Delete review",
		"User":     u,
		"ID":       nil,
		"ReviewID": id,
	})
}

func deleteForm(c *fiber.Ctx) error {
	u, _ := jwt.User(c)
	id := c.Params("id")

	return c.Render("style/deletereview", fiber.Map{
		"Title":    "Deleted review",
		"User":     u,
		"ID":       nil,
		"ReviewID": id,
	})
}
