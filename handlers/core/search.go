package core

import (
	"log"

	"github.com/gofiber/fiber/v2"

	"userstyles.world/handlers/jwt"
	"userstyles.world/search"
)

func Search(c *fiber.Ctx) error {
	u, _ := jwt.User(c)

	p := c.Params("query")
	log.Printf("%#+v", p)

	s, _ := search.SearchText(p)

	return c.Render("search", fiber.Map{
		"Title":  "Home",
		"User":   u,
		"Styles": s,
	})
}
