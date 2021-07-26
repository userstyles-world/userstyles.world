package core

import "github.com/gofiber/fiber/v2"

// Redirect function will return a function which will redirect to the given url.
func Redirect(url string) func(c *fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		return c.Redirect(url)
	}
}
