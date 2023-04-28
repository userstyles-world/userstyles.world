package core

import (
	"errors"

	"github.com/gofiber/fiber/v2"

	"userstyles.world/handlers/jwt"
	"userstyles.world/models"
	"userstyles.world/modules/config"
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

	page, err := models.IsValidPage(c.Query("page"))
	if err != nil || page < 1 {
		return c.Status(fiber.StatusBadRequest).Render("err", fiber.Map{
			"Title": "Invalid page size",
			"User":  u,
		})
	}

	sort := c.Query("sort")

	s, m, err := search.FindStylesByText(keyword, sort, page, config.AppPageMaxItems)
	if errors.Is(err, search.ErrSearchNoResults) {
		return c.
			Status(fiber.StatusNotFound).
			Render("core/search", fiber.Map{
				"Title":   "No results found",
				"User":    u,
				"Keyword": keyword,
				"Sort":    sort,
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

	p := models.NewPagination(page, m.Total, sort, c.Path())
	p.Query = keyword
	if p.OutOfBounds() {
		return c.Redirect(p.URL(p.Now), 302)
	}

	return c.Render("core/search", fiber.Map{
		"Title":     "Search",
		"User":      u,
		"Styles":    s,
		"Keyword":   keyword,
		"Sort":      p.Sort,
		"Canonical": "search",
		"Metrics":   m,
		"P":         p,
	})
}
