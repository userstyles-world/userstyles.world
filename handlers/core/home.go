package core

import (
	"github.com/gofiber/fiber/v2"

	"userstyles.world/database"
	"userstyles.world/handlers/sessions"
	"userstyles.world/models"
)

func Home(c *fiber.Ctx) error {
	s := sessions.State(c)

	styles, err := models.GetStyles(database.DB)
	if err != nil {
		return c.Render("index", fiber.Map{
			"Name":  s.Get("name"),
			"Title": "Home",
			"Styles": nil,
		})
	}

	return c.Render("index", fiber.Map{
		"Name":  s.Get("name"),
		"Title": "Home",
		"Styles": styles,
	})
}
