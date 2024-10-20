package style

import (
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

	i, err := c.ParamsInt("id")
	if err != nil || i < 1 {
		return c.Status(fiber.StatusBadRequest).Render("err", fiber.Map{
			"Title": "Invalid style ID",
			"User":  u,
		})
	}
	id := c.Params("id")

	// Check if style exists.
	s, err := models.GetStyleByID(i)
	if err != nil {
		return c.Render("err", fiber.Map{
			"Title": "Style not found",
			"User":  u,
		})
	}

	// Always redirect to correct slugged URL.
	if util.Slug(c.Params("name")) != s.Slug {
		return c.Redirect(s.Permalink(), fiber.StatusSeeOther)
	}

	args := fiber.Map{
		"User":       u,
		"Title":      s.Name,
		"Style":      s,
		"URL":        c.BaseURL() + "/style/" + id,
		"Canonical":  s.Permalink(),
		"RenderMeta": true,
	}

	// Upsert style views.
	if util.IsCrawler(string(c.Context().UserAgent())) {
		return c.Render("style/view", args)
	} else {
		cache.ViewStats.Add(c.IP() + " " + id)
	}

	if u.ID != s.UserID {
		if u.ID == 0 {
			args["CanReview"] = false
			args["CantReviewReason"] = "An account is required in order to review userstyles."
		} else {
			dur, ok := models.AbleToReview(u.ID, s.ID)
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

	reviews, err := models.FindAllForStyle(id)
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
