package style

import (
	"fmt"

	"github.com/gofiber/fiber/v2"

	"userstyles.world/database"
	"userstyles.world/handlers/sessions"
	"userstyles.world/models"
)

func DeleteByID(c *fiber.Ctx) error {
	u := sessions.User(c)
	p := c.Params("id")

	if sessions.State(c).Fresh() == true {
		c.Status(fiber.StatusUnauthorized)
		return c.Render("login", fiber.Map{
			"Title": "Login is required",
			"Error": "You must log in to add new userstyle.",
		})
	}

	s, err := models.GetStyleByID(database.DB, p)
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

	t := new(models.Style)
	err = database.DB.
		Debug().
		Delete(t, "styles.id = ?", p).
		Error

	if err != nil {
		fmt.Printf("Failed to delete style, err: %#+v\n", err)
		return c.Render("err", fiber.Map{
			"Title": "Internal server error",
			"User":  u,
		})
	}

	return c.Redirect("/account", fiber.StatusSeeOther)
}
