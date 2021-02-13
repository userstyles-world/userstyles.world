package style

import (
	"github.com/gofiber/fiber/v2"

	"userstyles.world/database"
	"userstyles.world/handlers/sessions"
	"userstyles.world/models"
)

func GetStyle(c *fiber.Ctx) error {
	u := sessions.User(c)

	data, err := models.GetStyleByID(database.DB, c.Params("id"))
	if err != nil {
		return c.Render("err", fiber.Map{
			"Title": "Style not found",
		})
	}

	return c.Render("style", fiber.Map{
		"Title": data.Name,
		"User":  u,
		"Style": data,
	})
}
