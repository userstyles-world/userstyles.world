package style

import (
	"github.com/gofiber/fiber/v2"

	"userstyles.world/config/database"
	"userstyles.world/handlers/jwt"
	"userstyles.world/models"
)

func Promote(c *fiber.Ctx) error {
	u, _ := jwt.User(c)
	p := c.Params("id")

	// Only moderator and above have permissions to promote styles.
	if u.Role < models.Moderator {
		return c.Render("err", fiber.Map{
			"Title": "You don't have enough permission for this.",
			"User":  u,
		})
	}

	style, err := models.GetStyleByID(p)
	if err != nil {
		return c.Render("err", fiber.Map{
			"Title": "Internal server error.",
			"User":  u,
		})
	}

	err = database.Conn.
		Model(models.Style{}).
		Where("id = ?", p).
		Update("featured", !style.Featured).
		Error

	if err != nil {
		return c.Render("err", fiber.Map{
			"Title": "Failed to promote a style",
			"User":  u,
		})
	}

	return c.Redirect("/style/"+p, fiber.StatusSeeOther)
}
