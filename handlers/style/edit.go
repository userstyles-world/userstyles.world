package style

import (
	"log"

	"github.com/gofiber/fiber/v2"

	"userstyles.world/database"
	"userstyles.world/handlers/sessions"
	"userstyles.world/models"
)

func StyleEditGet(c *fiber.Ctx) error {
	u := sessions.User(c)
	p := c.Params("id")

	if sessions.State(c).Fresh() {
		c.Status(fiber.StatusUnauthorized)
		return c.Render("login", fiber.Map{
			"Title": "Login is required",
			"Error": "You must log in to edit a userstyle.",
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
			"Title": "User and style author don't match",
			"User":  u,
		})
	}

	return c.Render("add", fiber.Map{
		"Title":  "Edit userstyle",
		"Method": "edit",
		"User":   u,
		"Style":  s,
	})
}

func StyleEditPost(c *fiber.Ctx) error {
	u, p := sessions.User(c), c.Params("id")
	t := new(models.Style)

	if sessions.State(c).Fresh() {
		c.Status(fiber.StatusUnauthorized)
		return c.Render("login", fiber.Map{
			"Title": "Login is required",
			"Error": "You must log in to edit a userstyle.",
		})
	}

	q := models.Style{
		Name:        c.FormValue("name"),
		Description: c.FormValue("description"),
		Notes:       c.FormValue("notes"),
		Code:        c.FormValue("code"),
		Preview:     c.FormValue("preview"),
		Homepage:    c.FormValue("homepage"),
		Category:    c.FormValue("category"),
		UserID:      u.ID,
	}

	err := database.DB.
		Model(t).
		Where("id", p).
		Updates(q).
		Error

	if err != nil {
		log.Println("Updating style failed, err:", err)
		return c.Render("err", fiber.Map{
			"Title": "Internal server error.",
			"User":  u,
		})
	}

	return c.Redirect("/style/"+c.Params("id"), fiber.StatusSeeOther)
}
