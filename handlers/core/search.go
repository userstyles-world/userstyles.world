package core

import (
	"sort"

	"github.com/userstyles-world/fiber/v2"

	"userstyles.world/handlers/jwt"
	"userstyles.world/search"
)

func Search(c *fiber.Ctx) error {
	u, _ := jwt.User(c)

	q := c.Query("q")
	s, metrics, _ := search.FindStylesByText(q)

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
		"Value":     q,
		"Root":      c.OriginalURL() == "/search",
		"Sort":      fv,
		"Canonical": "search",
		"Metrics":   metrics,
	})
}
