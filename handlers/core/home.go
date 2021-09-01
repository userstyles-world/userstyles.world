package core

import (
	"github.com/userstyles-world/fiber/v2"

	"userstyles.world/handlers/jwt"
	"userstyles.world/models"
)

func Home(c *fiber.Ctx) error {
	u, _ := jwt.User(c)

	// Skip stats if user is logged in.
	// TODO: Combine this with a new dashboard.
	var stats *models.SiteStats
	if u.ID == 0 {
		stats = models.GetHomepageStatistics()
	}

	styles, err := models.GetAllFeaturedStyles()
	if err != nil {
		return c.Render("core/home", fiber.Map{
			"Title":  "Home",
			"User":   u,
			"Styles": nil,
		})
	}

	return c.Render("core/home", fiber.Map{
		"Title":  "Home",
		"User":   u,
		"Styles": styles,
		"Stats":  stats,
	})
}
