package core

import (
	"net/http"
	"testing"

	"github.com/gofiber/fiber/v2"
)

func TestCSPMiddlewareCorrect(t *testing.T) {
	t.Parallel()
	app := fiber.New()
	app.Use(CSPMiddleware)

	// Returns an HTML response.
	app.Get("/html", func(c *fiber.Ctx) error {
		c.Response().Header.SetContentType(fiber.MIMETextHTMLCharsetUTF8)
		return c.SendString("<p>Hi</p>")
	})

	req, err := http.NewRequest("GET", "/html", nil)
	if err != nil {
		t.Error(err)
	}
	res, err := app.Test(req)
	if err != nil {
		t.Error(err)
	}
	if res.StatusCode != fiber.StatusOK {
		t.Error("Expected 200, got ", res.StatusCode)
	}
	if res.Header.Get(string(headerCSP)) != string(valueCSP) {
		t.Error("\nExpected ", string(valueCSP), "\ngot ", res.Header.Get(string(headerCSP)))
	}
}

func TestCSPMiddlewareInCorrect(t *testing.T) {
	t.Parallel()

	app := fiber.New()
	app.Use(CSPMiddleware)

	// Returns an json response.
	app.Get("/json", func(c *fiber.Ctx) error {
		c.Response().Header.SetContentType(fiber.MIMEApplicationJSONCharsetUTF8)
		return c.SendString(`{"name":"fiber"}`)
	})

	req, err := http.NewRequest("GET", "/json", nil)
	if err != nil {
		t.Error(err)
	}
	res, err := app.Test(req)
	if err != nil {
		t.Error(err)
	}
	if res.StatusCode != fiber.StatusOK {
		t.Error("Expected 200, got ", res.StatusCode)
	}
	if res.Header.Get(string(headerCSP)) != "" {
		t.Error("Expected empty CSP header, got ", res.Header.Get(string(headerCSP)))
	}
}
