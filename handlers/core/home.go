package core

import (
	"github.com/gofiber/fiber/v2"

	"userstyles.world/database"
	"userstyles.world/handlers/sessions"
	"userstyles.world/models"
)

func Home(c *fiber.Ctx) error {
	u := sessions.User(c)

	styles, err := models.GetAllFeaturedStyles(database.DB)
	if err != nil {
		return c.Render("index", fiber.Map{
			"Title":  "Home",
			"User":   u,
			"Styles": nil,
		})
	}

	return c.Render("index", fiber.Map{
		"Title":  "Home",
		"User":   u,
		"Styles": styles,
	})
}
