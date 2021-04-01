package style

import (
	"fmt"
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/vednoc/go-usercss-parser"

	"userstyles.world/database"
	"userstyles.world/handlers/jwt"
	"userstyles.world/models"
)

func StyleCreateGet(c *fiber.Ctx) error {
	u, _ := jwt.User(c)

	return c.Render("add", fiber.Map{
		"Title":  "Add userstyle",
		"User":   u,
		"Method": "add",
	})
}

func StyleCreatePost(c *fiber.Ctx) error {
	u, _ := jwt.User(c)

	s := &models.Style{
		Name:        c.FormValue("name"),
		Description: c.FormValue("description"),
		Notes:       c.FormValue("notes"),
		Homepage:    c.FormValue("homepage"),
		Category:    c.FormValue("category"),
		Preview:     c.FormValue("preview"),
		Code:        c.FormValue("code"),
		License:     c.FormValue("license", "No License"),
		UserID:      u.ID,
	}

	code := usercss.ParseFromString(c.FormValue("code"))
	valid, errs := usercss.BasicMetadataValidation(code)
	if !valid {
		return c.Render("add", fiber.Map{
			"Title": "Add userstyle",
			"User":  u,
			"Style": s,
			"UCSS":  errs,
		})
	}

	// Prevent adding multiples of the same style.
	err := models.CheckDuplicateStyleName(database.DB, s)
	if err != nil {
		return c.Render("err", fiber.Map{
			"Title": err,
			"User":  u,
		})
	}

	s, err = models.CreateStyle(database.DB, s)
	if err != nil {
		log.Println("Style creation failed, err:", err)
		return c.Render("err", fiber.Map{
			"Title": "Internal server error.",
			"User":  u,
		})
	}

	return c.Redirect(fmt.Sprintf("/style/%d", int(s.ID)), fiber.StatusSeeOther)
}
