package style

import (
	"fmt"
	"log"
	"strings"

	"github.com/gofiber/fiber/v2"

	"userstyles.world/database"
	"userstyles.world/handlers/jwt"
	"userstyles.world/models"
)

func GetStyle(c *fiber.Ctx) error {
	u, _ := jwt.User(c)
	id := c.Params("id")

	data, err := models.GetStyleByID(database.DB, id)
	if err != nil {
		return c.Render("err", fiber.Map{
			"Title": "Style not found",
			"User":  u,
		})
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
		"Title": data.Name,
		"User":  u,
		"Style": data,
		"Views": models.GetTotalViewsForStyle(database.DB, id),
		"Total": models.GetTotalInstallsForStyle(database.DB, id),
		"Week":  models.GetWeeklyInstallsForStyle(database.DB, id),
		"Url":   fmt.Sprintf("https://userstyles.world/style/%d", data.ID),
	})
}

func slugify(s string) string {
	s = strings.ReplaceAll(s, " ", "-")
	s = strings.ReplaceAll(s, ".", "-")
	s = strings.ToLower(s)

	return s
}

func GetStyleSlug(c *fiber.Ctx) error {
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
	slug := slugify(data.Name)

	// Always redirect to correct slugged URL.
	if name != slug {
		url := fmt.Sprintf("/style/%s/%s", id, slugify(data.Name))
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
		"Title": data.Name,
		"User":  u,
		"Style": data,
		"Views": models.GetTotalViewsForStyle(database.DB, id),
		"Total": models.GetTotalInstallsForStyle(database.DB, id),
		"Week":  models.GetWeeklyInstallsForStyle(database.DB, id),
		"Url":   fmt.Sprintf("https://userstyles.world/style/%d", data.ID),
		"Slug":  c.Path(),
	})
}
