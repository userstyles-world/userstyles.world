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

	return c.Render("core/explore", fiber.Map{
		"Title":  "Explore",
		"User":   u,
		"Styles": s,
		"Sort":   fv,
	})
}
