package style

import (
	"fmt"

	"github.com/gofiber/fiber/v2"

	"userstyles.world/handlers/jwt"
	"userstyles.world/models"
	"userstyles.world/modules/cache"
	"userstyles.world/modules/storage"

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
	slug := util.Slug(data.Name)

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
	if util.IsCrawler(string(c.Context().UserAgent())) {
		return c.Render("style/view", args)
	} else {
		cache.ViewStats.Add(c.IP() + " " + id)
	}

	if u.ID != data.UserID {
		if u.ID == 0 {
			args["CanReview"] = false
			args["CantReviewReason"] = "An account is required in order to review userstyles."
		} else {
			dur, ok := models.AbleToReview(u.ID, data.ID)
			if !ok {
				args["CanReview"] = false
				args["CantReviewReason"] = "You can review this style again in " + dur
			} else {
				args["CanReview"] = true
			}
		}
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

	stats, err := storage.GetStyleStats(id)
	if err != nil {
		log.Database.Printf("Failed to get stats: %s\n", err)
	}
	args["Stats"] = stats

	return c.Render("style/view", args)
}
