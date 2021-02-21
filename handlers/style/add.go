package style

import (
	"log"

	"github.com/gofiber/fiber/v2"

	"userstyles.world/database"
	"userstyles.world/handlers/sessions"
	"userstyles.world/models"
)

func StyleCreateGet(c *fiber.Ctx) error {
	u := sessions.User(c)

	if sessions.State(c).Fresh() == true {
		c.Status(fiber.StatusUnauthorized)
		return c.Render("login", fiber.Map{
			"Title": "Login is required",
			"Error": "You must log in to add new userstyle.",
		})
	}

	return c.Render("add", fiber.Map{
		"Title": "Add userstyle",
		"Name":  u,
	})
}

func StyleCreatePost(c *fiber.Ctx) error {
	u := sessions.User(c)

	if sessions.State(c).Fresh() == true {
		c.Status(fiber.StatusUnauthorized)
		return c.Render("login", fiber.Map{
			"Title": "Login is required",
			"Error": "You must log in to add new userstyle.",
		})
	}

	s := models.Style{
		Name:        c.FormValue("name"),
		Summary:     c.FormValue("summary"),
		Description: c.FormValue("description"),
		Preview:     c.FormValue("preview"),
		Code:        c.FormValue("code"),
		UserID:      u.ID,
	}

	err := database.DB.
		Debug().
		Create(&s).
		Error

	if err != nil {
		log.Println("Style creation failed, err:", err)
		return c.Render("err", fiber.Map{
			"Title": "Internal server error.",
			"User":  u,
		})
	}

	return c.Redirect("/account", fiber.StatusSeeOther)
}
