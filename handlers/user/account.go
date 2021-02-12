package user

import (
	"github.com/gofiber/fiber/v2"

	"userstyles.world/database"
	"userstyles.world/handlers/sessions"
	"userstyles.world/models"
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

	styles, err := models.GetStylesByUser(database.DB, s.Get("name").(string))
	if err != nil {
		return c.Render("err", fiber.Map{
			"Name":  s.Get("name"),
			"Title": "Server error",
		})
	}
	return c.Render("account", fiber.Map{
		"Name":   s.Get("name"),
		"Title":  "Account",
		"Styles": styles,
	})
}
