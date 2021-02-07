package user

import (
	"github.com/gofiber/fiber/v2"
	"userstyles.world/handlers/sessions"
)

func Logout(c *fiber.Ctx) error {
	s := sessions.State(c)
	s.Destroy()

	return c.Redirect("/login", fiber.StatusSeeOther)
}
