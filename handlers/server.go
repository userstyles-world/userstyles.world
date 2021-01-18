package handlers

import (
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/template/html"

	"userstyles.world/config"
	"userstyles.world/handlers/user"
)

func Initialize() {
	app := fiber.New(fiber.Config{
		Views: html.New("./views", ".html"),
	})

	app.Get("/", func(c *fiber.Ctx) error {
		return c.Render("index", fiber.Map{
			"Title": "UserStyles.world",
			"Body":  "Hello, World!",
		})
	})

	app.Get("/login", user.LoginGet)
	app.Post("/login", user.LoginPost)

	app.Get("/logout", user.Logout)

	app.Get("/register", user.RegisterGet)
	app.Post("/register", user.RegisterPost)

	log.Fatal(app.Listen(config.PORT))
}
