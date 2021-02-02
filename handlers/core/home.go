package core

import (
	"log"

	"github.com/gofiber/fiber/v2"

	"userstyles.world/handlers/sessions"
)

func Home(c *fiber.Ctx) error {
	s := sessions.State(c)

	log.Println(s, s.Get("name"))

	return c.Render("index", fiber.Map{
		"Name":  s.Get("name"),
		"Title": "Home",
	})
}
