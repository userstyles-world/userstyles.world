package user

import (
	"github.com/gofiber/fiber/v2"
)

func Logout(c *fiber.Ctx) error {
	c.ClearCookie(fiber.HeaderAuthorization)

	return c.Redirect("/signin", fiber.StatusSeeOther)
}
