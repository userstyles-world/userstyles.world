package style

import (
	"fmt"

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

	return c.Render("style", fiber.Map{
		"Title": data.Name,
		"User":  u,
		"Style": data,
		"Total": models.GetTotalInstallsForStyle(database.DB, id),
		"Week":  models.GetWeeklyInstallsForStyle(database.DB, id),
		"Url":   fmt.Sprintf("https://userstyles.world/style/%d", data.ID),
	})
}
