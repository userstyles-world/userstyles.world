package api

import (
	"github.com/gofiber/fiber/v2"

	"userstyles.world/modules/config"
	"userstyles.world/modules/search"
)

func GetSearchResult(c *fiber.Ctx) error {
	q := c.Params("query")

	// TODO: Add support for customizing parameters.
	styles, _, err := search.FindStylesByText(q, "", 1, config.AppPageMaxItems)
	if err != nil {
		return c.JSON(fiber.Map{"data": "no styles found"})
	}

	return c.JSON(fiber.Map{"data": styles})
}
