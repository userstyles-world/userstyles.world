package core

import (
	"errors"

	"github.com/gofiber/fiber/v2"

	"userstyles.world/handlers/jwt"
	"userstyles.world/modules/search"
)

func Search(c *fiber.Ctx) error {
	u, _ := jwt.User(c)

	keyword := c.Query("q")
	if keyword == "" {
		return c.Render("core/search", fiber.Map{
			"Title": "Search",
			"User":  u,
		})
	}

	// TODO: Refactor [probably] using Pagination struct.
	size := 96
	kind := c.Query("sort")
	if kind != "" {
		size = 500
	}

	s, metrics, err := search.FindStylesByText(keyword, kind, size)
	if errors.Is(err, search.ErrSearchNoResults) {
		return c.
			Status(fiber.StatusNotFound).
			Render("core/search", fiber.Map{
				"Title":   "No results found",
				"User":    u,
				"Keyword": keyword,
				"Sort":    kind,
				"Error":   "No results found for <b>" + keyword + "</b>.",
			})
	} else if err != nil {
		return c.
			Status(fiber.StatusBadRequest).
			Render("core/search", fiber.Map{
				"User":  u,
				"Title": "Bad search request",
				"Error": "Bad search request.",
			})
	}

	return c.Render("core/search", fiber.Map{
		"Title":     "Search",
		"User":      u,
		"Styles":    s,
		"Keyword":   keyword,
		"Sort":      kind,
		"Canonical": "search",
		"Metrics":   metrics,
	})
}
