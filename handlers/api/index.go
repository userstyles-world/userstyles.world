package api

import (
	"github.com/gofiber/fiber/v2"

	"userstyles.world/database"
	"userstyles.world/models"
)

func GetStyleIndex(c *fiber.Ctx) error {
	styles, err := models.GetAllStyles(database.DB)
	if err != nil {
		return c.JSON(fiber.Map{
			"data": "style not found",
		})
	}

	return c.JSON(fiber.Map{
		"data": styles,
	})
}
