package user

import (
	"log"

	"github.com/gofiber/fiber/v2"
	"userstyles.world/handlers/sessions"
)

func Logout(c *fiber.Ctx) error {
	if err := sessions.State(c).Destroy(); err != nil {
		log.Println("Wanted to destroy session, but caught error:", err)
	}

	return c.Redirect("/login", fiber.StatusSeeOther)
}
