package user

import (
	"github.com/gofiber/fiber/v2"

	"userstyles.world/modules/log"
)

func loginGet(c *fiber.Ctx) error {
	log.Spam.Printf("kind=loginGet ip=%q ua=%q",
		c.IP(), c.Context().UserAgent())
	return nil
}

func loginPost(c *fiber.Ctx) error {
	email := c.FormValue("email")
	password := c.FormValue("password")
	log.Spam.Printf("kind=loginPost ip=%q email=%q password=%d ua=%q",
		c.IP(), email, len(password), c.Context().UserAgent())
	return nil
}

func registerGet(c *fiber.Ctx) error {
	log.Spam.Printf("kind=registerGet ip=%q ua=%q",
		c.IP(), c.Context().UserAgent())
	return nil
}

func registerPost(c *fiber.Ctx) error {
	username := c.FormValue("username")
	password := c.FormValue("password")
	email := c.FormValue("email")
	log.Spam.Printf("kind=registerPost ip=%q un=%q email=%q password=%d ua=%q",
		c.IP(), username, email, len(password), c.Context().UserAgent())
	return nil
}
