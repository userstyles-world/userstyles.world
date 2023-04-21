package core

import (
	"encoding/json"

	"github.com/gofiber/fiber/v2"

	"userstyles.world/models"
)

// FlagsMiddleware gets flags cookie to enable feature flags.
func FlagsMiddleware(c *fiber.Ctx) error {
	var flags models.Flags
	if err := json.Unmarshal([]byte(c.Cookies("flags")), &flags); err != nil {
		return err
	}

	c.Locals("flags", flags)

	return c.Next()
}
