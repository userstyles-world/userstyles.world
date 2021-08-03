package api

import (
	"fmt"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"

	"userstyles.world/models"
	"userstyles.world/modules/cache"
	"userstyles.world/modules/log"
)

func convertToUSoFormat(s models.APIStyle) models.USoFormat {
	id := fmt.Sprintf("%d", s.ID) // Convert uint to string.

	var img string
	if s.Preview != "" {
		img = fmt.Sprintf("https://userstyles.world/api/style/preview/%d.webp", s.ID)
	}

	return models.USoFormat{
		ID:             s.ID,
		Name:           s.Name,
		Category:       fixCategory(s.Category),
		Username:       s.Username,
		Screenshot:     img,
		UpdatedAt:      s.UpdatedAt.Unix(),
		TotalInstalls:  models.GetTotalInstallsForStyle(id),
		WeeklyInstalls: models.GetWeeklyInstallsForStyle(id),
	}
}

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
			log.Warn.Fatal("API/index/uso-format err:", err.Error())
			return c.JSON(fiber.Map{
				"data": "styles not found",
			})
		}

		// TODO: Normalize categories on add/import/edit pages.
		for _, style := range *styles {
			style.Category = fixCategory(style.Category)
		}

		cache.Store.Set("index", styles, 10*time.Minute)
		if err := cache.SaveToDisk(cache.CachedIndex, *styles); err != nil {
			log.Warn.Println("Failed to cache USo-formatted index:", err)
		}

		goto Convert
	}

	return c.JSON(fiber.Map{
		"data": cached,
	})
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
