package core

import (
	"encoding/json"

	"github.com/gofiber/fiber/v2"

	"userstyles.world/models"
)

// FlagsMiddleware gets flags cookie to enable feature flags.
func FlagsMiddleware(c *fiber.Ctx) error {
	cookie := c.Cookies("flags")
	if cookie == "" {
		return c.Next()
	}

	var f models.Flags
	if err := json.Unmarshal([]byte(cookie), &f); err != nil {
		return err
	}

	c.Locals("flags", f)

	return c.Next()
}
