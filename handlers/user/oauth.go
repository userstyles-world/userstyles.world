package user

import (
	"log"

	"github.com/gofiber/fiber/v2"
	"userstyles.world/handlers/jwt"
	"userstyles.world/utils"
)

func AuthLoginGet(c *fiber.Ctx) error {
	if u, ok := jwt.User(c); ok {
		log.Printf("User %d has set session, redirecting.", u.ID)
		c.Redirect("/account", fiber.StatusSeeOther)
	}

	oathType := c.Params("type")
	redirectURI := ""

	switch oathType {
	case "github":
		redirectURI = utils.GithubMakeURL(c.BaseURL())
	}

	return c.Redirect(redirectURI, fiber.StatusSeeOther)
}
