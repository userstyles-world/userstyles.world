package core

import (
	"net/http"
	"testing"

	"github.com/userstyles-world/fiber/v2"
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
	if res.Header.Get(string(headerCSP)) != string(valueCSPStrictForm) {
		t.Error("\nExpected ", string(valueCSPStrictForm), "\ngot ", res.Header.Get(string(headerCSP)))
	}
}

// Because of some unexepected behavior.
// In the CSP Policy on any page that starts with /api/oauth/style/.
// It will add an exception to the form-action CSP directive.
func TestCSPException(t *testing.T) {
	t.Parallel()

	app := fiber.New()
	app.Use(CSPMiddleware)
	app.Get("/api/oauth/style/new", func(c *fiber.Ctx) error {
		c.Response().Header.SetContentType(fiber.MIMETextHTMLCharsetUTF8)
		return c.SendString("<p>Hi, wanna add a style?</p>")
	})
	app.Get("/api/oauth/style/link", func(c *fiber.Ctx) error {
		c.Response().Header.SetContentType(fiber.MIMETextHTMLCharsetUTF8)
		return c.SendString("<p>Hi linking?</p>")
	})

	req, err := http.NewRequest("GET", "/api/oauth/style/new", nil)
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

	req, err = http.NewRequest("GET", "/api/oauth/style/link", nil)
	if err != nil {
		t.Error(err)
	}
	res, err = app.Test(req)
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
