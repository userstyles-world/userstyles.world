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

	app.Static("/", "./static")

	app.Get("/", func(c *fiber.Ctx) error {
		return c.Render("index", fiber.Map{
			"Title": "Home",
		})
	})

	app.Get("/login", user.LoginGet)
	app.Post("/login", user.LoginPost)

	app.Post("/logout", user.Logout)
	app.Get("/account", user.Account)

	app.Get("/register", user.RegisterGet)
	app.Post("/register", user.RegisterPost)

	log.Fatal(app.Listen(config.PORT))
}
