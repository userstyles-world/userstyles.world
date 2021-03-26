package style

import (
	"log"
	"strings"

	"github.com/gofiber/fiber/v2"

	"userstyles.world/database"
	"userstyles.world/handlers/jwt"
	"userstyles.world/models"
)

func StyleEditGet(c *fiber.Ctx) error {
	u, _ := jwt.User(c)
	p := c.Params("id")

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
	u, _ := jwt.User(c)
	p, t := c.Params("id"), new(models.Style)

	q := models.Style{
		Name:        c.FormValue("name"),
		Description: c.FormValue("description"),
		Notes:       c.FormValue("notes"),
		Code:        c.FormValue("code"),
		Preview:     c.FormValue("preview"),
		Homepage:    c.FormValue("homepage"),
		License:     strings.TrimSpace(c.FormValue("license", "No License")),
		Category:    strings.TrimSpace(c.FormValue("category", "global")),
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
