package core

import (
	"github.com/gofiber/fiber/v2"

	"userstyles.world/handlers/jwt"
	"userstyles.world/models"
)

func Home(c *fiber.Ctx) error {
	u, _ := jwt.User(c)

	// Skip stats if user is logged in.
	// TODO: Combine this with a new dashboard.
	var p *models.SiteStats
	if u.ID == 0 {
		p = models.GetHomepageStatistics()
	}

	styles, err := models.GetAllFeaturedStyles()
	if err != nil {
		return c.Render("index", fiber.Map{
			"Title":  "Home",
			"User":   u,
			"Styles": nil,
		})
	}

	return c.Render("index", fiber.Map{
		"Title":  "Home",
		"User":   u,
		"Styles": styles,
		"Params": p,
	})
}
