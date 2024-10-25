package api

import (
	"github.com/gofiber/fiber/v2"

	"userstyles.world/models"
)

func GetStyleDetails(c *fiber.Ctx) error {
	i, err := c.ParamsInt("id")
	if err != nil || i < 1 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "invalid userstyle ID",
		})
	}

	s, err := models.GetStyleSourceCodeAPI(i)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"message": "style not found",
		})
	}

	return c.JSON(fiber.Map{
		"data": s,
	})
}
