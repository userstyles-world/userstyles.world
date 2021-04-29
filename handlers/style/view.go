package style

import (
	"fmt"
	"log"
	"regexp"
	"strings"

	"github.com/gofiber/fiber/v2"

	"userstyles.world/database"
	"userstyles.world/handlers/jwt"
	"userstyles.world/models"
)

func slugify(s string) string {
	re := regexp.MustCompile(`[a-zA-Z0-9]+`)

	// Extract valid characters.
	parts := re.FindAllString(s, -1)
	fmt.Printf("parts: %#+v\n", parts)

	joined := strings.Join(parts, "-")
	fmt.Printf("joined: %#+v\n", joined)

	s = strings.ToLower(joined)
	fmt.Printf("output: %#+v\n", s)

	return s
}

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
