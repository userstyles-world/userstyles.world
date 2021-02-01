package core

import (
	"log"

	"github.com/gofiber/fiber/v2"

	"userstyles.world/handlers/sessions"
)

var (
	store = sessions.GetStore()
)

func Home(c *fiber.Ctx) error {
	s, err := store.Get(c)
	if err != nil {
		log.Println(err)
	}

	log.Println(s, err, s.Get("name"))

	return c.Render("index", fiber.Map{
		"Name":  s.Get("name"),
		"Title": "Home",
	})
}
