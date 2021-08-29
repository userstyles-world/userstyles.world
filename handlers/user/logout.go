package user

import (
	"github.com/userstyles-world/fiber/v2"
)

func Logout(c *fiber.Ctx) error {
	c.ClearCookie(fiber.HeaderAuthorization)

	return c.Redirect("/login", fiber.StatusSeeOther)
}
