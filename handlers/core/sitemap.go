package core

import (
	"github.com/userstyles-world/fiber/v2"
	"userstyles.world/modules/sitemap"
)

func GetSiteMap(c *fiber.Ctx) error {
	// Check if sitemap is already generated
	if sitemap.SiteMapCache == nil {
		err := sitemap.UpdateSitemapCache()
		if err != nil {
			return c.Status(500).SendString("Sitemap generation failed")
		}
	}
	c.Response().Header.SetContentType(fiber.MIMEApplicationXMLCharsetUTF8)
	return c.Send(sitemap.SiteMapCache)
}
