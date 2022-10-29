package style

import (
	"fmt"
	"strings"

	"github.com/gofiber/fiber/v2"

	"userstyles.world/handlers/jwt"
	"userstyles.world/models"
	"userstyles.world/modules/cache"
	"userstyles.world/modules/charts"
	"userstyles.world/modules/log"
	"userstyles.world/utils/strutils"
)

func GetStylePage(c *fiber.Ctx) error {
	u, _ := jwt.User(c)
	id, name := strings.Clone(c.Params("id")), c.Params("name")

	// Check if style exists.
	data, err := models.GetStyleByID(id)
	if err != nil {
		return c.Render("err", fiber.Map{
			"Title": "Style not found",
			"User":  u,
		})
	}

	// Create slugged URL.
	// TODO: Refactor after GetStyleByID switches away from APIStyle.
	slug := strutils.SlugifyURL(data.Name)

	// Always redirect to correct slugged URL.
	if name != slug {
		url := fmt.Sprintf("/style/%s/%s", id, slug)
		return c.Redirect(url, fiber.StatusSeeOther)
	}

	// Upsert style views.
	ua := string(c.Context().UserAgent())
	if strings.Contains(ua, "Mastodon") || strings.Contains(ua, "Pleroma") {
		log.Info.Printf("Ignored Fediverse bot on style %s: %q\n", id, ua)
	} else {
		cache.ViewStats.Add(c.IP() + " " + id)
	}

	// Get history data.
	history, err := models.GetStyleHistory(id)
	if err != nil {
		log.Warn.Printf("Failed to get stats for style %s: %s", id, err.Error())
	}

	// Render graphs.
	var dailyHistory, totalHistory string
	if len(history) > 2 {
		dailyHistory, totalHistory, err = charts.GetStatsHistory(history)
		if err != nil {
			log.Warn.Printf("Failed to render history for style %s: %s\n", id, err.Error())
		}
	}

	reviews, err := new(models.Review).FindAllForStyle(id)
	if err != nil {
		log.Warn.Printf("Failed to get reviews for style %s: %v\n", id, err.Error())
	}

	// Get stats.
	stats := models.GetStyleStatistics(id)

	return c.Render("style/view", fiber.Map{
		"Title":        data.Name,
		"User":         u,
		"Style":        data,
		"Stats":        stats,
		"URL":          c.BaseURL() + c.Path(),
		"DailyHistory": dailyHistory,
		"TotalHistory": totalHistory,
		"Reviews":      reviews,
		"Canonical":    "style/" + id + "/" + slug,
		"RenderMeta":   true,
	})
}
