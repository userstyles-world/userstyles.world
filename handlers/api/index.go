package api

import (
	"os"
	"strings"

	"github.com/gofiber/fiber/v2"

	"userstyles.world/models"
	"userstyles.world/modules/cache"
	"userstyles.world/modules/log"
)

func fixCategory(cat string) string {
	if cat == "unset" {
		return "global"
	}
	cat = strings.ToLower(cat)

	if strings.HasSuffix(cat, ".com") || strings.HasSuffix(cat, ".org") {
		cat = strings.TrimSuffix(cat, ".com")
		cat = strings.TrimSuffix(cat, ".org")
		// Remove any subdomain
		// web.whatsapp -> whatsapp
		if strings.Count(cat, ".") >= 1 {
			cat = strings.Split(cat, ".")[1]
		}
	} else if strings.Count(cat, ".") >= 2 {
		cat = strings.Join(strings.Split(cat, ".")[1:], ".")
	}

	return cat
}

func getUSoIndex(c *fiber.Ctx) error {
Convert:
	cached, found := cache.Store.Get("index")
	if !found {
		styles := new(models.USoStyles)
		if err := styles.Query(); err != nil {
			log.Warn.Println("Failed to get styles for USo-formatted index:", err.Error())
			return c.JSON(fiber.Map{
				"data": "styles not found",
			})
		}

		// TODO: Normalize categories on add/import/edit pages.
		for _, style := range *styles {
			style.Category = fixCategory(style.Category)
		}

		// Save to disk and read it to avoid converting between types.
		if err := cache.SaveToDisk(cache.CachedIndex, fiber.Map{
			"data": *styles,
		}); err != nil {
			log.Warn.Println("Failed to cache USo-formatted index:", err)
			goto Convert
		}
		b, err := os.ReadFile(cache.CachedIndex)
		if err != nil {
			log.Warn.Println("Failed to read uso-format.json:", err)
			goto Convert
		}

		// Set cache for index endpoint.
		cache.Store.Set("index", b, 0)

		goto Convert
	}

	c.Set("Content-Type", "application/json")
	return c.Send(cached.([]byte))
}

func getFullIndex(c *fiber.Ctx) error {
	styles, err := models.GetAllStylesForIndexAPI()
	if err != nil {
		return c.JSON(fiber.Map{
			"data": "styles not found",
		})
	}

	return c.JSON(fiber.Map{
		"data": styles,
	})
}

func GetStyleIndex(c *fiber.Ctx) error {
	switch c.Params("format") {
	case "uso-format":
		return getUSoIndex(c)
	default:
		return getFullIndex(c)
	}
}
