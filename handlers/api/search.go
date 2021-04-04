package api

import (
	"strings"

	"github.com/gofiber/fiber/v2"
	"userstyles.world/search"
)

func GetSearchResult(c *fiber.Ctx) error {
	searchQuery := c.Params("query")

	results, err := search.SearchText(searchQuery)
	if err != nil {
		return c.Render("err", fiber.Map{
			"Title": "Internal server error.",
		})
	}
	return c.SendString(strings.Join(results, ", "))
}
