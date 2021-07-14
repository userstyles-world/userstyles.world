package user

import (
	"log"

	"github.com/gofiber/fiber/v2"

	"userstyles.world/handlers/jwt"
	"userstyles.world/modules/oauthlogin"
)

func AuthLoginGet(c *fiber.Ctx) error {
	if u, ok := jwt.User(c); ok {
		log.Printf("User %d has set session, redirecting.", u.ID)
		return c.Redirect("/account", fiber.StatusSeeOther)
	}

	oauthType := c.Params("type")
	redirectURI := oauthlogin.OauthMakeURL(oauthType)

	return c.Redirect(redirectURI, fiber.StatusSeeOther)
}
