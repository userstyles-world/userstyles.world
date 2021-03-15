package api

import (
	"github.com/gofiber/fiber/v2"

	"userstyles.world/database"
	"userstyles.world/models"
)

func GetStyleDetails(c *fiber.Ctx) error {
	id := c.Params("id")

	style, err := models.GetStyleSourceCodeAPI(database.DB, id)
	if err != nil {
		return c.JSON(fiber.Map{
			"data": "style not found",
		})
	}

	return c.JSON(fiber.Map{
		"data": style,
	})
}
