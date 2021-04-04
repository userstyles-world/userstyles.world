package api

import (
	"encoding/json"

	"github.com/gofiber/fiber/v2"
	"userstyles.world/search"
	"userstyles.world/utils"
)

func GetSearchResult(c *fiber.Ctx) error {
	searchQuery := c.Params("query")

	results, err := search.SearchText(searchQuery)
	if err != nil {
		return c.Render("err", fiber.Map{
			"Title": "Internal server error.",
		})
	}
	var StylesInfo string

	for _, hit := range results {
		json, err := json.Marshal(hit)
		if err != nil {
			return c.Render("err", fiber.Map{
				"Title": "Internal server error.",
			})
		}
		StylesInfo += utils.B2s(json)
	}
	return c.SendString(StylesInfo)
}
