package user

import (
	"log"

	"github.com/gofiber/fiber/v2"
)

func Account(c *fiber.Ctx) error {
	sess, err := store.Get(c)
	if err != nil {
		log.Println(err)
	}

	return c.Render("account", fiber.Map{
		"Name":  sess.Get("name"),
		"Title": "UserStyles.world",
		"Body":  "Hello, World!",
	})
}
