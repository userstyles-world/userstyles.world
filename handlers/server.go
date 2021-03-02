package handlers

import (
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/monitor"
	"github.com/gofiber/template/html"

	"userstyles.world/config"
	"userstyles.world/handlers/api"
	"userstyles.world/handlers/core"
	"userstyles.world/handlers/style"
	"userstyles.world/handlers/user"
	"userstyles.world/utils"
)

func Initialize() {
	app := fiber.New(fiber.Config{
		Views:                 html.New("./views", ".html"),
		DisableStartupMessage: true,
	})

	app.Static("/", "./static")

	app.Get("/", core.Home)

	app.Get("/login", user.LoginGet)
	app.Post("/login", user.LoginPost)

	app.Post("/logout", user.Logout)
	app.Get("/account", user.Account)

	app.Get("/register", user.RegisterGet)
	app.Post("/register", user.RegisterPost)

	app.Get("/user/:name", user.Profile)

	app.Get("/explore", style.GetExplore)
	app.Get("/style/:id", style.GetStyle)
	app.Post("/style/:id", style.DeleteByID)
	app.Get("/add", style.StyleCreateGet)
	app.Post("/add", style.StyleCreatePost)

	app.Get("/api/style/:id.user.css", api.GetStyleSource)

	// Good luck landing on this route. (=
	app.Get("/monitor", core.Monitor)
	app.Get(utils.MonitorURL, monitor.New())

	app.Use(core.NotFound)

	log.Fatal(app.Listen(config.PORT))
}
