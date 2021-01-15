package main

import (
	"log"

	"userstyles.world/config"
	"userstyles.world/database"
	"userstyles.world/models"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/template/html"
)

func main() {
	database.Connect()
	database.Prepare()

	app := fiber.New(fiber.Config{
		Views: html.New("./views", ".html"),
	})

	app.Get("/", func(c *fiber.Ctx) error {
		return c.Render("index", fiber.Map{
			"Title": "UserStyles.world",
			"Body":  "Hello, World!",
		})
	})

	log.Fatal(app.Listen(config.PORT))
}
