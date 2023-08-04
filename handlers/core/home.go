package core

import (
	"time"

	"github.com/gofiber/fiber/v2"

	"userstyles.world/handlers/jwt"
	"userstyles.world/modules/cache"
	"userstyles.world/modules/log"
	"userstyles.world/modules/storage"
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
	updatedStyles, err := storage.FindStyleCardsPaginated(1, 4, "styles.updated_at DESC")
	if err != nil {
		//idk
	}
	addedStyles, err := storage.FindStyleCardsPaginated(1, 4, "styles.created_at DESC")
	if err != nil {
		//idk
	}
	if !found {
		styles, err := storage.FindStyleCardsFeatured()
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
		"Title":         "Website themes and skins",
		"User":          u,
		"Styles":        featured,
		"UpdatedStyles": updatedStyles,
		"AddedStyles":   addedStyles,
		// "Stats":  stats,
	})
}
