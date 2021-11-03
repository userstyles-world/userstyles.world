package core

import (
	"io"
	"net/http"
	"testing"

	"github.com/gofiber/fiber/v2"
	"userstyles.world/modules/templates"
)

// Check if you get the correct 404 page.
// Depending on the request.
func Test404Normal(t *testing.T) {
	t.Parallel()

	app := fiber.New()
	app.Get("/", func(c *fiber.Ctx) error {
		c.Response().Header.SetContentType(fiber.MIMETextHTMLCharsetUTF8)
		return c.SendString("<p>Hi, wanna have some fun</p>")
	})
	app.Get("/api/index", func(c *fiber.Ctx) error {
		c.Response().Header.SetContentType(fiber.MIMEApplicationJSONCharsetUTF8)
		return c.SendString(`{"name":"fiber"}`)
	})
	app.Use(NotFound)

	req, err := http.NewRequest("GET", "/", nil)
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

	req, err = http.NewRequest("GET", "/api/index", nil)
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
}

func Test404Pages(t *testing.T) {
	t.Parallel()

	app := fiber.New(fiber.Config{
		Views: templates.New("../../views"),
	})
	app.Get("/", func(c *fiber.Ctx) error {
		c.Response().Header.SetContentType(fiber.MIMETextHTMLCharsetUTF8)
		return c.SendString("<p>Hi, wanna have some fun</p>")
	})
	app.Get("/api/index", func(c *fiber.Ctx) error {
		c.Response().Header.SetContentType(fiber.MIMEApplicationJSONCharsetUTF8)
		return c.SendString(`{"name":"fiber"}`)
	})
	app.Use(NotFound)

	req, err := http.NewRequest("GET", "/api/notfound", nil)
	if err != nil {
		t.Error(err)
	}
	res, err := app.Test(req)
	if err != nil {
		t.Error(err)
	}
	if res.StatusCode != fiber.StatusNotFound {
		t.Error("Expected 404, got ", res.StatusCode)
	}

	expected := `{"error":"bad endpoint"}`
	body := make([]byte, len(expected))
	_, err = res.Body.Read(body)
	if err != nil && err != io.EOF {
		t.Error(err)
	}
	if string(body) != expected {
		t.Error("Expected ", expected, " got ", string(body))
	}

	req, err = http.NewRequest("GET", "/notfound", nil)
	if err != nil {
		t.Error(err)
	}
	res, err = app.Test(req)
	if err != nil {
		t.Error(err)
	}
	if res.StatusCode != fiber.StatusNotFound {
		t.Error("Expected 404, got ", res.StatusCode)
	}
}
