package api

import (
	"github.com/gofiber/fiber/v2"

	"userstyles.world/models"
)

func GetStyleStats(c *fiber.Ctx) error {
	id := c.Params("id")

	return c.JSON(fiber.Map{
		"total_views":    models.GetTotalViewsForStyle(id),
		"total_installs": models.GetTotalInstallsForStyle(id),
	})
}
