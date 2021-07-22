package style

import (
	"github.com/gofiber/fiber/v2"

	"userstyles.world/handlers/jwt"
	"userstyles.world/models"
	"userstyles.world/modules/database"
	"userstyles.world/modules/log"
	"userstyles.world/search"
)

func DeleteGet(c *fiber.Ctx) error {
	u, _ := jwt.User(c)
	p := c.Params("id")

	s, err := models.GetStyleByID(p)
	if err != nil {
		return c.Render("err", fiber.Map{
			"Title": "Style not found",
			"User":  u,
		})
	}

	// Check if logged-in user matches style author.
	if u.ID != s.UserID {
		return c.Render("err", fiber.Map{
			"Title": "Users don't match",
			"User":  u,
		})
	}

	return c.Render("style/delete", fiber.Map{
		"Title": "Confirm deletion",
		"User":  u,
		"Style": s,
	})
}

func DeletePost(c *fiber.Ctx) error {
	u, _ := jwt.User(c)
	p := c.Params("id")

	s, err := models.GetStyleByID(p)
	if err != nil {
		return c.Render("err", fiber.Map{
			"Title": "Style not found",
			"User":  u,
		})
	}

	// Check if logged-in user matches style author.
	if u.ID != s.UserID {
		return c.Render("err", fiber.Map{
			"Title": "Users don't match",
			"User":  u,
		})
	}

	q := new(models.Style)
	if err = database.Conn.Delete(q, "styles.id = ?", p).Error; err != nil {
		log.Warn.Printf("Failed to delete style %d: %s\n", s.ID, err.Error())
		return c.Render("err", fiber.Map{
			"Title": "Internal server error",
			"User":  u,
		})
	}

	if err = search.DeleteStyle(s.ID); err != nil {
		log.Warn.Printf("Failed to delte style %d: %s", s.ID, err.Error())
	}

	return c.Redirect("/account", fiber.StatusSeeOther)
}
