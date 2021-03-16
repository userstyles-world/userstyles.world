package user

import (
	"github.com/gofiber/fiber/v2"

	"userstyles.world/database"
	"userstyles.world/handlers/sessions"
	"userstyles.world/models"
)

func Account(c *fiber.Ctx) error {
	u := sessions.User(c)

	if sessions.State(c).Fresh() {
		c.Status(fiber.StatusUnauthorized)
		return c.Render("login", fiber.Map{
			"Title": "Login is required",
			"Error": "You must log in to see account page.",
		})
	}

	styles, err := models.GetStylesByUser(database.DB, u.Username)
	if err != nil {
		return c.Render("err", fiber.Map{
			"Title": "Server error",
			"User":  u,
		})
	}

	return c.Render("account", fiber.Map{
		"Title":  "Account",
		"User":   u,
		"Styles": styles,
	})
}
