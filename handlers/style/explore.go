package style

import (
	"sort"

	"github.com/gofiber/fiber/v2"

	"userstyles.world/database"
	"userstyles.world/handlers/jwt"
	"userstyles.world/models"
)

func GetExplore(c *fiber.Ctx) error {
	u, _ := jwt.User(c)

	s, err := models.GetAllAvailableStyles(database.DB)
	if err != nil {
		return c.Render("err", fiber.Map{
			"Title": "Styles not found",
			"User":  u,
		})
	}

	sort.Slice(s, func(i, j int) bool {
		switch c.Query("sort") {
		case "created":
			return s[i].CreatedAt.Unix() > s[j].CreatedAt.Unix()
		case "updated":
			return s[i].UpdatedAt.Unix() > s[j].UpdatedAt.Unix()
		case "installs":
			return s[i].Installs > s[j].Installs
		case "views":
			return s[i].Views > s[j].Views
		default:
			return s[i].Installs > s[j].Installs
		}
	})

	return c.Render("explore", fiber.Map{
		"Title":  "Explore",
		"User":   u,
		"Styles": s,
	})
}
