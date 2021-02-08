package style

import (
	"github.com/gofiber/fiber/v2"

	"userstyles.world/database"
	"userstyles.world/models"
)

func GetExplore(c *fiber.Ctx) error {

	data, err := models.GetAllStyles(database.DB)
	if err != nil {
		return c.Render("err", fiber.Map{
			"Title": "Styles not found",
		})
	}

	return c.Render("explore", fiber.Map{
		// "Name":   u,
		"Title":  "Explore",
		"Styles": data,
	})
}
