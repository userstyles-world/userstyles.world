package style

import (
	"fmt"

	"github.com/gofiber/fiber/v2"

	"userstyles.world/handlers/jwt"
	"userstyles.world/models"
	"userstyles.world/modules/charts"
	"userstyles.world/modules/log"
	"userstyles.world/utils/strings"
)

func GetStylePage(c *fiber.Ctx) error {
	u, _ := jwt.User(c)
	id, name := c.Params("id"), c.Params("name")

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
	slug := strings.SlugifyURL(data.Name)

	// Always redirect to correct slugged URL.
	if name != slug {
		url := fmt.Sprintf("/style/%s/%s", id, slug)
		return c.Redirect(url, fiber.StatusSeeOther)
	}

	go func(id, ip string) {
		// Count views.
		if _, err := models.AddStatsToStyle(id, ip, false); err != nil {
			log.Warn.Printf("Failed to add stats to style %s: %s\n", id, err.Error())
		}
	}(id, c.IP())

	// Get history data.
	history, err := new(models.History).GetStatsForStyle(id)
	if err != nil {
		log.Warn.Printf("Failed to get stats for style %s: %s", id, err.Error())
	}

	// Render graphs.
	var dailyHistory, totalHistory string
	if len(*history) > 2 {
		dailyHistory, totalHistory, err = charts.GetStatsHistory(*history)
		if err != nil {
			log.Warn.Printf("Failed to render history for style %s: %s\n", id, err.Error())
		}
	}

	reviews, err := new(models.Review).FindAllForStyle(id)
	if err != nil {
		log.Warn.Printf("Failed to get reviews for style %s: %v\n", id, err.Error())
	}

	return c.Render("style/view", fiber.Map{
		"Title":          data.Name,
		"User":           u,
		"Style":          data,
		"TotalViews":     models.GetTotalViewsForStyle(id),
		"TotalInstalls":  models.GetTotalInstallsForStyle(id),
		"WeeklyInstalls": models.GetWeeklyInstallsForStyle(id),
		"WeeklyUpdates":  models.GetWeeklyUpdatesForStyle(id),
		"Url":            fmt.Sprintf("https://userstyles.world/style/%d", data.ID),
		"Slug":           c.Path(),
		"BaseURL":        c.BaseURL(),
		"DailyHistory":   dailyHistory,
		"TotalHistory":   totalHistory,
		"Reviews":        reviews,
		"Canonical":      "style/" + id + "/" + slug,
	})
}
