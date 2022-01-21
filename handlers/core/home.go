package core

import (
	"time"

	"github.com/gofiber/fiber/v2"

	"userstyles.world/handlers/jwt"
	"userstyles.world/models"
	"userstyles.world/modules/cache"
	"userstyles.world/modules/log"
)

func Home(c *fiber.Ctx) error {
	u, _ := jwt.User(c)

	// Skip stats if user is logged in.
	// TODO: Combine this with a new dashboard.
	/*
			var stats *models.SiteStats
			if u.ID == 0 {
			Stats:
				cached, found := cache.Store.Get("siteStatistics")
				if !found {
					stats = models.GetHomepageStatistics()
					cache.Store.Set("siteStatistics", stats, 5*time.Minute)

				goto Stats
			}

			stats = cached.(*models.SiteStats)
		}
	*/

Styles:
	featured, found := cache.Store.Get("featuredStyles")
	if !found {
		styles, err := models.GetAllFeaturedStyles()
		if err != nil {
			log.Warn.Println("Couldn't get featured styles, due", err)
			return c.Render("core/home", fiber.Map{
				"Title":  "Home",
				"User":   u,
				"Styles": nil,
			})
		}
		cache.Store.Set("featuredStyles", styles, 5*time.Minute)

		goto Styles
	}

	return c.Render("core/home", fiber.Map{
		"Title":  "Home",
		"User":   u,
		"Styles": featured,
		// "Stats":  stats,
	})
}
