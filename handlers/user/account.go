package user

import (
	"github.com/gofiber/fiber/v2"
	"userstyles.world/handlers/sessions"
)

func Account(c *fiber.Ctx) error {
	s := sessions.State(c)

	if s.Fresh() == true {
		c.Status(fiber.StatusFound)

		return c.Render("login", fiber.Map{
			"Title": "Login is required",
			"Error": "You must log in to see account page.",
		})
	}

	return c.Render("account", fiber.Map{
		"Name":  s.Get("name"),
		"Title": "Account",
	})
}
