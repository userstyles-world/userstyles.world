package api

import (
	"github.com/gofiber/fiber/v2"

	"userstyles.world/database"
	"userstyles.world/models"
)

func GetStyleSource(c *fiber.Ctx) error {
	id := c.Params("id")

	style, err := models.GetStyleSourceCodeAPI(database.DB, id)
	if err != nil {
		return c.JSON(fiber.Map{"data": "style not found"})
	}

	c.Set("Content-Type", "text/css")
	return c.SendString(style.Code)
}
