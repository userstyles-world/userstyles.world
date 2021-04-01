package style

import (
	"github.com/gofiber/fiber/v2"

	"userstyles.world/database"
	"userstyles.world/handlers/jwt"
	"userstyles.world/models"
)

func StylePromote(c *fiber.Ctx) error {
	u, _ := jwt.User(c)
	p := c.Params("id")

	// TODO: Make it possible to remove promotion.
	err := database.DB.
		Model(models.Style{}).
		Where("id = ?", p).
		Update("featured", true).
		Error

	if err != nil {
		c.Render("err", fiber.Map{
			"Title": "Failed to promote a style",
			"User":  u,
		})
	}

	return c.Redirect("/style/"+p, fiber.StatusSeeOther)
}
