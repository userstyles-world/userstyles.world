package api

import (
	"github.com/gofiber/fiber/v2"

	"userstyles.world/search"
)

func GetSearchResult(c *fiber.Ctx) error {
	q := c.Params("query")

	styles, _, err := search.FindStylesByText(q)
	if err != nil {
		return c.JSON(fiber.Map{"data": "no styles found"})
	}

	return c.JSON(fiber.Map{"data": styles})
}
