package style

import (
	"sort"

	"github.com/gofiber/fiber/v2"

	"userstyles.world/handlers/jwt"
	"userstyles.world/models"
)

func GetExplore(c *fiber.Ctx) error {
	u, _ := jwt.User(c)

	s, err := models.GetAllAvailableStyles()
	if err != nil {
		return c.Render("err", fiber.Map{
			"Title": "Styles not found",
			"User":  u,
		})
	}

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
			return s[i].CreatedAt.Unix() > s[j].CreatedAt.Unix()
		}
	})

	return c.Render("core/explore", fiber.Map{
		"Title":  "Explore",
		"User":   u,
		"Styles": s,
		"Sort":   fv,
	})
}
