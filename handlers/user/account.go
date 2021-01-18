package user

import (
	"log"

	"github.com/gofiber/fiber/v2"
)

func Account(c *fiber.Ctx) error {
	s, err := store.Get(c)
	if err != nil {
		log.Println(err)
	}

	if s.Fresh() == true {
		c.Status(fiber.StatusFound)
		return c.Render("login", fiber.Map{
			"Error": "You must log in to see account page.",
		})
	}

	return c.Render("account", fiber.Map{
		"Name":  s.Get("name"),
		"Title": "UserStyles.world",
		"Body":  "Hello, World!",
	})
}
