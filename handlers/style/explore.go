package style

import (
	"github.com/gofiber/fiber/v2"

	"userstyles.world/database"
	"userstyles.world/handlers/sessions"
	"userstyles.world/models"
)

func GetExplore(c *fiber.Ctx) error {
	s := sessions.State(c)
	u := s.Get("name")

	data, err := models.GetAllStyles(database.DB)
	if err != nil {
		return c.Render("err", fiber.Map{
			"Title": "Styles not found",
		})
	}

	return c.Render("explore", fiber.Map{
		"Title":  "Explore",
		"Name":   u,
		"Styles": data,
	})
}
