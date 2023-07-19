package middleware

import (
	"github.com/gofiber/fiber/v2"

	"userstyles.world/handlers/jwt"
	"userstyles.world/modules/cache"
)

// Alert middleware shows alerts after page redirection.
func Alert(c *fiber.Ctx) error {
	u, _ := jwt.User(c)

	k := "alert " + u.Username
	msg, ok := cache.Store.Get(k)
	if !ok {
		return c.Next()
	}

	c.Locals("alert", msg)
	cache.Store.Delete(k)

	return c.Next()
}
