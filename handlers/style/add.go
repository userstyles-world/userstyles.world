package style

import (
	"github.com/gofiber/fiber/v2"

	"userstyles.world/handlers/sessions"
)

func StyleCreateGet(c *fiber.Ctx) error {
	u := sessions.User(c)

	if sessions.State(c).Fresh() == true {
		c.Status(fiber.StatusUnauthorized)
		return c.Render("login", fiber.Map{
			"Title": "Login is required",
			"Error": "You must log in to add new userstyle.",
		})
	}

	return c.Render("add", fiber.Map{
		"Title": "Add userstyle",
		"Name":  u,
	})
}
