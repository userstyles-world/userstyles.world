package core

import (
	"log"

	"github.com/gofiber/fiber/v2"
)

func NotFound(c *fiber.Ctx) error {
	s, err := store.Get(c)
	if err != nil {
		log.Println(err)
	}

	return c.Render("404", fiber.Map{
		"Name":  s.Get("name"),
		"Title": "Page not found",
	})
}
