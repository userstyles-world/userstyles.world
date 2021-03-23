package handlers

import (
	"log"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cache"
	"github.com/gofiber/fiber/v2/middleware/compress"
	"github.com/gofiber/fiber/v2/middleware/filesystem"
	"github.com/gofiber/fiber/v2/middleware/limiter"
	"github.com/gofiber/template/html"
	"github.com/markbates/pkger"

	"userstyles.world/config"
	"userstyles.world/handlers/api"
	"userstyles.world/handlers/core"
	"userstyles.world/handlers/jwt"
	"userstyles.world/handlers/style"
	"userstyles.world/handlers/user"
)

func Initialize() {
	engine := html.NewFileSystem(pkger.Dir("/views"), ".html")

	app := fiber.New(fiber.Config{
		Views:                 engine,
		DisableStartupMessage: true,
	})

	app.Use(cache.New(cache.Config{
		Expiration:   5 * time.Minute,
		CacheControl: true,
	}))
	app.Use(compress.New())
	if config.DB != "dev.db" {
		app.Use(limiter.New())
	}

	app.Get("/", core.Home)
	app.Get("/login", jwt.NoLoggedInUsers, user.LoginGet)
	app.Post("/login", user.LoginPost)

	app.Post("/logout", jwt.Protected, user.Logout)
	app.Get("/account", jwt.Protected, user.Account)
	app.Post("/account", jwt.Protected, user.EditAccount)

	app.Get("/register", jwt.NoLoggedInUsers, user.RegisterGet)
	app.Post("/register", user.RegisterPost)

	app.Get("/user/:name", jwt.Everyone, user.Profile)

	app.Get("/explore", jwt.Everyone, style.GetExplore)
	app.Get("/style/:id", jwt.Everyone, style.GetStyle)
	app.Post("/style/:id", jwt.Protected, style.DeleteByID)
	app.Get("/add", jwt.Protected, style.StyleCreateGet)
	app.Post("/add", jwt.Protected, style.StyleCreatePost)
	app.Get("/import", jwt.Protected, style.StyleImportGet)
	app.Post("/import", jwt.Protected, style.StyleImportPost)
	app.Get("/edit/:id", jwt.Protected, style.StyleEditGet)
	app.Post("/edit/:id", jwt.Protected, style.StyleEditPost)

	app.Get("/api/style/:id.user.css", api.GetStyleSource)
	app.Get("/api/style/:id", api.GetStyleDetails)
	app.Get("/api/styles", api.GetStyleIndex)

	app.Get("/monitor", jwt.Protected, core.Monitor)

	if config.DB == "prod.dev" {
		app.Use(limiter.New())
	}
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
