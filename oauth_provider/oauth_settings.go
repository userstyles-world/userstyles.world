package oauth_provider

import (
	"github.com/gofiber/fiber/v2"
	"userstyles.world/database"
	"userstyles.world/handlers/jwt"
	"userstyles.world/models"
)

func OAuthSettingsGet(c *fiber.Ctx) error {
	u, _ := jwt.User(c)
	id := c.Params("id")

	var method string
	var OAuth *models.APIOAuth
	var err error
	if id != "" {
		method = "edit"
		OAuth, err = models.GetOAuthByID(database.DB, id)
	} else {
		method = "add"
	}

	oauths, err := models.GetOAuthByUser(database.DB, u.Username)
	if err != nil {
		return c.Render("oauth_settings", fiber.Map{
			"Title":  "OAuth Settings",
			"User":   u,
			"Oauth":  OAuth,
			"OAuths": nil,
			"Method": method,
		})
	}

	return c.Render("oauth_settings", fiber.Map{
		"Title":  "OAuth Settings",
		"User":   u,
		"Oauth":  OAuth,
		"OAuths": oauths,
		"Method": method,
	})
}
