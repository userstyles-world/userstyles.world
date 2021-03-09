package user

import (
	"github.com/gofiber/fiber/v2"

	"userstyles.world/database"
	"userstyles.world/handlers/sessions"
	"userstyles.world/models"
)

func Profile(c *fiber.Ctx) error {
	u := sessions.User(c)
	p := c.Params("name")

	user, err := models.FindUserByName(database.DB, p)
	if err != nil {
		return c.Render("err", fiber.Map{
			"Title": "User not found",
			"User":  u,
		})
	}

	styles, err := models.GetStylesByUser(database.DB, p)
	if err != nil {
		return c.Render("err", fiber.Map{
			"User":  u,
			"Title": "Server error",
		})
	}

	// Render Account template if current user matches active session.
	if u.Username == p {
		return c.Render("account", fiber.Map{
			"Title":  "Account",
			"User":   u,
			"Styles": styles,
		})
	}

	return c.Render("profile", fiber.Map{
		"Title":  "Profile",
		"User":   u,
		"Params": user,
		"Styles": styles,
	})
}
