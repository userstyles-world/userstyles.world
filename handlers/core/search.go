package core

import (
	"sort"

	"github.com/gofiber/fiber/v2"

	"userstyles.world/handlers/jwt"
	"userstyles.world/search"
)

func Search(c *fiber.Ctx) error {
	u, _ := jwt.User(c)

	q := c.Query("q")
	s, _ := search.FindStylesByText(q)

	fv := c.Query("sort")
	sort.Slice(s, func(i, j int) bool {
		switch fv {
		case "newest":
			return s[i].CreatedAt.Unix() > s[j].CreatedAt.Unix()
		case "oldest":
			return s[i].CreatedAt.Unix() < s[j].CreatedAt.Unix()
		case "recentlyupdated":
			return s[i].UpdatedAt.Unix() > s[j].UpdatedAt.Unix()
		case "leastupdated":
			return s[i].UpdatedAt.Unix() < s[j].UpdatedAt.Unix()
		case "mostinstalls":
			return s[i].Installs > s[j].Installs
		case "leastinstalls":
			return s[i].Installs < s[j].Installs
		case "mostviews":
			return s[i].Views > s[j].Views
		case "leastviews":
			return s[i].Views < s[j].Views
		default:
			return s[i].CreatedAt.Unix() < s[j].CreatedAt.Unix()
		}
	})

	return c.Render("core/search", fiber.Map{
		"Title":  "Search",
		"User":   u,
		"Styles": s,
		"Value":  q,
		"Root":   c.OriginalURL() == "/search",
		"Sort":   fv,
	})
}
