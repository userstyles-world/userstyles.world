package api

import (
	"fmt"

	"github.com/gofiber/fiber/v2"

	"userstyles.world/modules/storage"
)

type badgeSchema struct {
	SchemaVersion int    `json:"schemaVersion"`
	Label         string `json:"label"`
	Message       string `json:"message"`
	Color         string `json:"color"`
}

func GetStyleStats(c *fiber.Ctx) error {
	id, t := c.Params("id"), c.Params("type")

	if t == "" {
		return c.JSON(fiber.Map{
			"total_views":    storage.GetTotalViews(id),
			"total_installs": storage.GetTotalInstalls(id),
		})
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
			storage.GetWeeklyViews(id),
			storage.GetTotalViews(id),
		)
	case "installs":
		badge.Message = fmt.Sprintf(
			"%d weekly, %d total installs",
			storage.GetWeeklyInstalls(id),
			storage.GetTotalInstalls(id),
		)
	default:
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"data": "Error: Invalid type parameter",
		})
	}

	return c.JSON(badge)
}
