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
		case "created":
			return s[i].CreatedAt.Unix() > s[j].CreatedAt.Unix()
		case "updated":
			return s[i].UpdatedAt.Unix() > s[j].UpdatedAt.Unix()
		case "installs":
			return s[i].Installs > s[j].Installs
		case "views":
			return s[i].Views > s[j].Views
		default:
			return s[i].CreatedAt.Unix() > s[j].CreatedAt.Unix()
		}
	})

	return c.Render("search", fiber.Map{
		"Title":  "Search",
		"User":   u,
		"Styles": s,
		"Value":  q,
		"Root":   c.OriginalURL() == "/search",
		"Sort":   fv,
	})
}
