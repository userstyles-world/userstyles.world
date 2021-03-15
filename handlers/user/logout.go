package user

import (
	"log"

	"github.com/gofiber/fiber/v2"
	"userstyles.world/handlers/sessions"
)

func Logout(c *fiber.Ctx) error {
	s := sessions.State(c)
	if err := s.Destroy(); err != nil {
		log.Println("Wanted to destroy session, but caught error:", err)
	}

	return c.Redirect("/login", fiber.StatusSeeOther)
}
