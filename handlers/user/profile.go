package user

import (
	"strings"

	"github.com/gofiber/fiber/v2"

	"userstyles.world/handlers/jwt"
	"userstyles.world/models"
)

func Profile(c *fiber.Ctx) error {
	u, _ := jwt.User(c)
	p := c.Params("name")

	profile, err := models.FindUserByName(p)
	if err != nil {
		return c.Render("err", fiber.Map{
			"Title": "User not found",
			"User":  u,
		})
	}

	// Always redirect to correct URL.
	if p != profile.Username {
		return c.Redirect("/user/"+strings.ToLower(p), fiber.StatusSeeOther)
	}

	styles, err := models.GetStylesByUser(p)
	if err != nil {
		return c.Render("err", fiber.Map{
			"User":  u,
			"Title": "Server error",
		})
	}

	// Render Account template if current user matches active session.
	/*if u.Username == p {
		return c.Render("user/account", fiber.Map{
			"Title":  "Account",
			"User":   u,
			"Styles": styles,
		})
	}*/

	return c.Render("user/profile", fiber.Map{
		"Title":     profile.Name() + "'s profile",
		"User":      u,
		"Profile":   profile,
		"Styles":    styles,
		"Canonical": "user/" + p,
	})
}
