package style

import (
	"fmt"
	"log"
	"regexp"

	"github.com/gofiber/fiber/v2"
	"github.com/vednoc/go-usercss-parser"

	"userstyles.world/database"
	"userstyles.world/handlers/sessions"
	"userstyles.world/models"
)

func StyleCreateGet(c *fiber.Ctx) error {
	u := sessions.User(c)

	if sessions.State(c).Fresh() {
		c.Status(fiber.StatusUnauthorized)
		return c.Render("login", fiber.Map{
			"Title": "Login is required",
			"Error": "You must log in to add new userstyle.",
		})
	}

	return c.Render("add", fiber.Map{
		"Title":  "Add userstyle",
		"User":   u,
		"Method": "add",
	})
}

func StyleCreatePost(c *fiber.Ctx) error {
	u := sessions.User(c)

	if sessions.State(c).Fresh() {
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

	// Check if source code is a link.
	r, err := regexp.Compile(`^https?://.*\.user\.(css|styl|less)$`)
	if err != nil {
		return c.Render("err", fiber.Map{
			"Title": "Internal server error",
		})
	}

	// Redirect to external userstyle.
	if r.MatchString(s.Code) {
		ext, err := usercss.ParseFromURL(s.Code)
		if err != nil {
			return c.Render("err", fiber.Map{
				"Title": "Failed to fetch external userstyle",
				"User":  u,
			})
		}

		// Check if external userstyle is valid.
		if valid, _ := usercss.BasicMetadataValidation(ext); !valid {
			return c.Render("err", fiber.Map{
				"Title": "Failed to validate external userstyle",
				"User":  u,
			})
		}

		s.Code = ext.SourceCode
	} else {
		form := usercss.ParseFromString(c.FormValue("code"))
		valid, errs := usercss.BasicMetadataValidation(form)
		if !valid {
			return c.Render("add", fiber.Map{
				"Title": "Add userstyle",
				"User":  u,
				"Style": s,
				"UCSS":  errs,
			})
		}
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
