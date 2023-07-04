package api

import (
	"github.com/gofiber/fiber/v2"
)

func GetSearchResult(c *fiber.Ctx) error {
	// TODO: Re-implement search using new search engine.
	return c.JSON(fiber.Map{"data": "temporarily disabled"})
}
