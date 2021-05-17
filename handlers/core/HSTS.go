package core

import "github.com/gofiber/fiber/v2"

var (
	header = []byte("Strict-Transport-Security")
	// Make it so that that the max-age is 2 years, it's high enough.
	// We also say that every subdomain has HTTPS://
	// And we say that we are preload(on browser's list).
	value = []byte("max-age=63072000; includeSubDomains; preload")
)

// HSTSMiddleware adds the HSTS Header
// See: https://developer.mozilla.org/en-US/docs/Web/HTTP/Headers/Strict-Transport-Security
func HSTSMiddleware(c *fiber.Ctx) error {
	c.Response().Header.SetCanonical(header, value)
	return c.Next()
}
