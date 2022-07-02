package core

import (
	"errors"
	"sort"

	"github.com/gofiber/fiber/v2"

	"userstyles.world/handlers/jwt"
	"userstyles.world/search"
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

	s, metrics, err := search.FindStylesByText(keyword)
	if errors.Is(err, search.ErrSearchNoResults) {
		return c.
			Status(fiber.StatusNotFound).
			Render("core/search", fiber.Map{
				"Title": "No results found",
				"Error": "No results found for <b>" + keyword + "</b>.",
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

	fv := c.Query("sort")

	var sortFunction func(i, j int) bool
	switch fv {
	case "newest":
		sortFunction = func(i, j int) bool { return s[i].CreatedAt.Unix() > s[j].CreatedAt.Unix() }
	case "oldest":
		sortFunction = func(i, j int) bool { return s[i].CreatedAt.Unix() < s[j].CreatedAt.Unix() }
	case "recentlyupdated":
		sortFunction = func(i, j int) bool { return s[i].UpdatedAt.Unix() > s[j].UpdatedAt.Unix() }
	case "leastupdated":
		sortFunction = func(i, j int) bool { return s[i].UpdatedAt.Unix() < s[j].UpdatedAt.Unix() }
	case "mostinstalls":
		sortFunction = func(i, j int) bool { return s[i].Installs > s[j].Installs }
	case "leastinstalls":
		sortFunction = func(i, j int) bool { return s[i].Installs < s[j].Installs }
	case "mostviews":
		sortFunction = func(i, j int) bool { return s[i].Views > s[j].Views }
	case "leastviews":
		sortFunction = func(i, j int) bool { return s[i].Views < s[j].Views }
	}
	if sortFunction != nil {
		sort.Slice(s, sortFunction)
	}

	return c.Render("core/search", fiber.Map{
		"Title":     "Search",
		"User":      u,
		"Styles":    s,
		"Keyword":   keyword,
		"Sort":      fv,
		"Canonical": "search",
		"Metrics":   metrics,
	})
}
