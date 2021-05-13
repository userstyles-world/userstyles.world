package style

import (
	"fmt"
	"log"

	"github.com/gofiber/fiber/v2"

	"userstyles.world/database"
	"userstyles.world/handlers/jwt"
	"userstyles.world/models"
	"userstyles.world/utils"
)

func GetStylePage(c *fiber.Ctx) error {
	u, _ := jwt.User(c)
	id, name := c.Params("id"), c.Params("name")

	// Check if style exists.
	data, err := models.GetStyleByID(database.DB, id)
	if err != nil {
		return c.Render("err", fiber.Map{
			"Title": "Style not found",
			"User":  u,
		})
	}

	// Create slugged URL.
	slug := utils.SlugifyURL(data.Name)

	// Always redirect to correct slugged URL.
	if name != slug {
		url := fmt.Sprintf("/style/%s/%s", id, slug)
		return c.Redirect(url, fiber.StatusSeeOther)
	}

	// Count views.
	_, err = models.AddStatsToStyle(database.DB, id, c.IP(), false)
	if err != nil {
		log.Println("Failed to add stats to style, err:", err)
		return c.Render("err", fiber.Map{
			"Title": "Internal server error",
			"User":  u,
		})
	}

	return c.Render("style", fiber.Map{
		"Title":          data.Name,
		"User":           u,
		"Style":          data,
		"TotalViews":     models.GetTotalViewsForStyle(database.DB, id),
		"TotalInstalls":  models.GetTotalInstallsForStyle(database.DB, id),
		"WeeklyInstalls": models.GetWeeklyInstallsForStyle(database.DB, id),
		"WeeklyUpdates":  models.GetWeeklyUpdatesForStyle(database.DB, id),
		"Url":            fmt.Sprintf("https://userstyles.world/style/%d", data.ID),
		"Slug":           c.Path(),
	})
}
