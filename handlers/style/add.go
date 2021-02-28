package style

import (
	"fmt"
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/vednoc/go-usercss-parser"

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
		Homepage:    c.FormValue("homepage"),
		Category:    c.FormValue("category"),
		Preview:     c.FormValue("preview"),
		Code:        c.FormValue("code"),
		UserID:      u.ID,
	}

	uc := usercss.ParseFromString(c.FormValue("code"))
	valid, errs := usercss.BasicMetadataValidation(uc)
	if !valid {
		return c.Render("add", fiber.Map{
			"Title":  "Add userstyle",
			"User":   u,
			"Style":  s,
			"Errors": errs,
		})
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
