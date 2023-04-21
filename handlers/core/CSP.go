package core

import (
	"strings"

	"github.com/gofiber/fiber/v2"
)

var (
	headerCSP = []byte(fiber.HeaderContentSecurityPolicy)

	valueCSPStrictForm = append(valueCSP, []byte(" form-action 'self';")...)
	valueCSP           = []byte("default-src 'none'; font-src 'self'; img-src 'self'; script-src 'self'; style-src 'self' 'unsafe-inline'; frame-ancestors 'none'; upgrade-insecure-requests; base-uri 'none'; object-src 'none'; worker-src 'none'; child-src 'none'; frame-src 'none'; connect-src 'self';")
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

	switch {
	// Disable CSP for our Monitor service.
	case strings.HasPrefix(c.Path(), "/monitor"):
		return nil

	// Chromium doesn't allow redirects after form submissions, so we disable
	// form-action rule on /api/oauth/style/* endpoints.
	case strings.HasPrefix(c.Path(), "/api/oauth/style/"):
		c.Response().Header.SetCanonical(headerCSP, valueCSP)

	// Apply strict CSP to all other endpoints.
	default:
		c.Response().Header.SetCanonical(headerCSP, valueCSPStrictForm)
	}

	return nil
}
