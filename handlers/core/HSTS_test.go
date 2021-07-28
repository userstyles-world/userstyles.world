package core

import (
	"net/http"
	"testing"

	"github.com/gofiber/fiber/v2"
)

func TestHSTSMiddleware(t *testing.T) {
	t.Parallel()

	app := fiber.New()
	app.Use(HSTSMiddleware)

	// Returns a 200
	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("")
	})

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
	if res.Header.Get(string(headerHSTS)) != string(valueHSTS) {
		t.Error("\nExpected ", string(valueCSP), "\ngot ", res.Header.Get(string(headerCSP)))
	}
}
