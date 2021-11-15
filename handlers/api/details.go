package api

import (
	"strconv"

	"github.com/gofiber/fiber/v2"

	"userstyles.world/models"
)

func GetStyleDetails(c *fiber.Ctx) error {
	i, err := c.ParamsInt("id")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "invalid userstyle ID",
		})
	}
	id := strconv.Itoa(i)

	style, err := models.GetStyleSourceCodeAPI(id)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"message": "style not found",
		})
	}

	return c.JSON(fiber.Map{
		"data": style,
	})
}
