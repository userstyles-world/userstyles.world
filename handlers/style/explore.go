package style

import (
	"github.com/gofiber/fiber/v2"

	"userstyles.world/database"
	"userstyles.world/handlers/jwt"
	"userstyles.world/models"
)

func GetExplore(c *fiber.Ctx) error {
	u := jwt.User(c)

	data, err := models.GetAllStyles(database.DB)
	if err != nil {
		return c.Render("err", fiber.Map{
			"Title": "Styles not found",
			"User":  u,
		})
	}

	return c.Render("explore", fiber.Map{
		"Title":  "Explore",
		"User":   u,
		"Styles": data,
	})
}
