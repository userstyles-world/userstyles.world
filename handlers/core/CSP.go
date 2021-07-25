package core

import (
	"github.com/gofiber/fiber/v2"
)

var (
	headerCSP = []byte(fiber.HeaderContentSecurityPolicy)

	valueCSP = []byte("default-src 'none'; font-src https://fonts.imma.link; img-src 'self'; script-src 'self'; style-src 'self' 'unsafe-inline'; form-action 'self'; frame-ancestors 'none'; upgrade-insecure-requests; base-uri: 'none'; object-src: 'none'; worker-src: 'none'; child-src: 'none'; frame-src: 'none';")
)

// CSPMiddleware adds the CSP Header
// See: https://developer.mozilla.org/en-US/docs/Web/HTTP/Headers/Content-Security-Policy
func CSPMiddleware(c *fiber.Ctx) error {
	// Continue stack
	if err := c.Next(); err != nil {
		return err
	}
	// Check if the response is text/html
	if string(c.Response().Header.Peek(fiber.HeaderContentType)) != fiber.MIMETextHTMLCharsetUTF8 {
		return nil
	}
	c.Response().Header.SetCanonical(headerCSP, valueCSP)
	return nil
}
