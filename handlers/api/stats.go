package api

import (
	"fmt"

	"github.com/userstyles-world/fiber/v2"

	"userstyles.world/models"
)

type badgeSchema struct {
	SchemaVersion int    `json:"schemaVersion"`
	Label         string `json:"label"`
	Message       string `json:"message"`
	Color         string `json:"color"`
}

func fullStats(c *fiber.Ctx, id string) error {
	return c.JSON(fiber.Map{
		"total_views":    models.GetTotalViewsForStyle(id),
		"total_installs": models.GetTotalInstallsForStyle(id),
	})
}

func GetStyleStats(c *fiber.Ctx) error {
	id, t := c.Params("id"), c.Params("type")

	if t == "" {
		return fullStats(c, id)
	}

	badge := badgeSchema{
		SchemaVersion: 1,
		Label:         "UserStyles.world",
		Color:         "#679cd0",
	}

	switch c.Params("type") {
	case "views":
		badge.Message = fmt.Sprintf(
			"%d weekly, %d total views",
			models.GetWeeklyViewsForStyle(id),
			models.GetTotalViewsForStyle(id),
		)
	case "installs":
		badge.Message = fmt.Sprintf(
			"%d weekly, %d total installs",
			models.GetWeeklyInstallsForStyle(id),
			models.GetTotalInstallsForStyle(id),
		)
	default:
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"data": "Error: Invalid type parameter",
		})
	}

	return c.JSON(badge)
}
