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

	_, err := models.FindUserByName(database.DB, p)
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

	return c.Render("profile", fiber.Map{
		"Title":  "Profile",
		"User":   u,
		"Name":   p,
		"Styles": styles,
	})
}
