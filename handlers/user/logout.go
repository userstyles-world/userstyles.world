package user

import (
	"log"

	"github.com/gofiber/fiber/v2"
)

func Logout(c *fiber.Ctx) error {
	sess, err := store.Get(c)
	if err != nil {
		log.Println(err)
	}
	sess.Destroy()

	return c.Redirect("/login", fiber.StatusFound)
}
