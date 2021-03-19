package handlers

import (
	"log"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cache"
	"github.com/gofiber/fiber/v2/middleware/filesystem"
	"github.com/gofiber/fiber/v2/middleware/limiter"
	"github.com/gofiber/template/html"
	"github.com/markbates/pkger"

	"userstyles.world/config"
	"userstyles.world/handlers/api"
	"userstyles.world/handlers/core"
	"userstyles.world/handlers/style"
	"userstyles.world/handlers/user"
)

func Initialize() {
	engine := html.NewFileSystem(pkger.Dir("/views"), ".html")

	app := fiber.New(fiber.Config{
		Views:                 engine,
		DisableStartupMessage: true,
	})

	app.Get("/", core.Home)

	app.Get("/login", user.LoginGet)
	app.Post("/login", user.LoginPost)

	app.Post("/logout", user.Logout)
	app.Get("/account", user.Account)
	app.Post("/account", user.EditAccount)

	app.Get("/register", user.RegisterGet)
	app.Post("/register", user.RegisterPost)

	app.Get("/user/:name", user.Profile)

	app.Get("/explore", style.GetExplore)
	app.Get("/style/:id", style.GetStyle)
	app.Post("/style/:id", style.DeleteByID)
	app.Get("/add", style.StyleCreateGet)
	app.Post("/add", style.StyleCreatePost)
	app.Get("/import", style.StyleImportGet)
	app.Post("/import", style.StyleImportPost)
	app.Get("/edit/:id", style.StyleEditGet)
	app.Post("/edit/:id", style.StyleEditPost)

	app.Get("/api/style/:id.user.css", api.GetStyleSource)
	app.Get("/api/style/:id", api.GetStyleDetails)
	app.Get("/api/styles", api.GetStyleIndex)

	app.Get("/monitor", core.Monitor)

	app.Use(limiter.New())
	app.Use(cache.New(cache.Config{
		Expiration: 5 * time.Minute,
	}))

	// Allows assets to be reloaded in dev mode.
	// That means, they're not embedded into executable file.
	if config.DB == "dev.db" {
		app.Static("/", "/static")
	}

	app.Use("/", filesystem.New(filesystem.Config{
		Root: pkger.Dir("/static"),
	}))
	app.Use(core.NotFound)

	log.Fatal(app.Listen(config.PORT))
}
