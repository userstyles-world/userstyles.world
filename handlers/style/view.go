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

	data, err := models.GetStyleByID(database.DB, c.Params("id"))
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
		"Url":   fmt.Sprintf("https://userstyles.world/style/%d", data.ID),
	})
}
