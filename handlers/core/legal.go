package core

import (
	"os"

	"github.com/gofiber/fiber/v2"
)

func GetLegal(c *fiber.Ctx) error {
	document := c.Params("document")

	switch document {
	case "privacy_policy":
		content, err := os.ReadFile("privacy_policy.md")
		if err != nil {
			c.Render("err", fiber.Map{
				"Title": "Couldn't load the document.",
			})
		}
		return c.Render("legal", fiber.Map{
			"content": string(content),
		})
	case "terms_of_services":
		content, err := os.ReadFile("terms_of_services.md")
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
