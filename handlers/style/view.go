package style

import (
	"fmt"
	"strings"

	"github.com/gofiber/fiber/v2"

	"userstyles.world/handlers/jwt"
	"userstyles.world/models"
	"userstyles.world/modules/cache"
	// "userstyles.world/modules/charts"
	"userstyles.world/modules/log"
	"userstyles.world/modules/util"
)

func GetStylePage(c *fiber.Ctx) error {
	u, _ := jwt.User(c)
	id := c.Params("id")

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
	slug := util.SlugifyURL(data.Name)

	// Always redirect to correct slugged URL.
	if c.Params("name") != slug {
		url := fmt.Sprintf("/style/%s/%s", id, slug)
		return c.Redirect(url, fiber.StatusSeeOther)
	}

	args := fiber.Map{
		"User":       u,
		"Title":      data.Name,
		"Style":      data,
		"URL":        c.BaseURL() + c.Path(),
		"Canonical":  "style/" + id + "/" + slug,
		"RenderMeta": true,
	}

	// Upsert style views.
	ua := strings.ToLower(string(c.Context().UserAgent()))
	if strings.Contains(ua, "bot") {
		log.Info.Printf("Ignored a bot on style %s: %q\n", id, ua)
		return c.Render("style/view", args)
	} else {
		cache.ViewStats.Add(c.IP() + " " + id)
	}

    /*
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
	args["DailyHistory"] = dailyHistory
	args["TotalHistory"] = totalHistory
    */

	reviews, err := new(models.Review).FindAllForStyle(id)
	if err != nil {
		log.Warn.Printf("Failed to get reviews for style %s: %v\n", id, err.Error())
	}
	args["Reviews"] = reviews

	// Get stats.
	stats := models.GetStyleStatistics(id)
	args["Stats"] = stats

	return c.Render("style/view", args)
}
