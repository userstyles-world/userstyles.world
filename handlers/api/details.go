package api

import (
	"github.com/userstyles-world/fiber/v2"

	"userstyles.world/models"
)

func GetStyleDetails(c *fiber.Ctx) error {
	id := c.Params("id")

	style, err := models.GetStyleSourceCodeAPI(id)
	if err != nil {
		return c.JSON(fiber.Map{
			"data": "style not found",
		})
	}

	return c.JSON(fiber.Map{
		"data": style,
	})
}
