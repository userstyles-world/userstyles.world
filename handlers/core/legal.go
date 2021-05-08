package core

import (
	"os"

	"github.com/gofiber/fiber/v2"
)

func GetLegal(c *fiber.Ctx) error {
	document := c.Params("document")

	switch document {
	case "privacy-policy":
		content, err := os.ReadFile("privacy-policy.md")
		if err != nil {
			c.Render("err", fiber.Map{
				"Title": "Couldn't load the document.",
			})
		}
		return c.Render("legal", fiber.Map{
			"content": string(content),
		})
	case "terms-of-service":
		content, err := os.ReadFile("terms-of-service.md")
		if err != nil {
			c.Render("err", fiber.Map{
				"Title": "Couldn't load the document.",
			})
		}
		return c.Render("legal", fiber.Map{
			"content": string(content),
		})
	}

	return c.Render("err", fiber.Map{
		"Title": "Couldn't load the document.",
	})
}
