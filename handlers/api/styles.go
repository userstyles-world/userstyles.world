package api

import (
	"github.com/gofiber/fiber/v2"
	"userstyles.world/database"
	"userstyles.world/models"
	"userstyles.world/utils"
)

func StylesGet(c *fiber.Ctx) error {
	u, _ := APIUser(c)

	if !utils.Contains(u.Scopes, "style") {
		return c.Status(403).
			JSON(fiber.Map{
				"data": "You need the \"style\" scope to do this.",
			})
	}

	styles, err := models.GetStylesByUser(database.DB, u.Username)
	if err != nil {
		return c.Status(500).
			JSON(fiber.Map{
				"data": "Error: Couldn't find styles.",
			})
	}

	return c.JSON(fiber.Map{
		"data": styles,
	})

}
