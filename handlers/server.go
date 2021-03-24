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

<<<<<<< Caching-Branch
	app.Use(cache.New(cache.Config{
		Expiration:   5 * time.Minute,
		CacheControl: true,
		KeyGenerator: func(c *fiber.Ctx) string {
			return c.Path() + c.Get(fiber.HeaderAcceptEncoding)
		},
	}))
=======
>>>>>>> refactor: use custom middleware
	app.Use(compress.New())
	if config.DB != "dev.db" {
		app.Use(limiter.New())
	}
	app.Use(jwt.New())

	app.Get("/", core.Home)
	app.Get("/login", user.LoginGet)
	app.Post("/login", user.LoginPost)
	app.Get("/register", user.RegisterGet)
	app.Post("/register", user.RegisterPost)
	app.Get("/explore", style.GetExplore)
	app.Get("/style/:id", style.GetStyle)
	app.Get("/user/:name", user.Profile)

	app.Post("/logout", jwt.Protected, user.Logout)
	app.Get("/account", jwt.Protected, user.Account)
	app.Post("/account", jwt.Protected, user.EditAccount)
	app.Post("/style/:id", jwt.Protected, style.DeleteByID)
	app.Get("/add", jwt.Protected, style.StyleCreateGet)
	app.Post("/add", jwt.Protected, style.StyleCreatePost)
	app.Get("/import", jwt.Protected, style.StyleImportGet)
	app.Post("/import", jwt.Protected, style.StyleImportPost)
	app.Get("/edit/:id", jwt.Protected, style.StyleEditGet)
	app.Post("/edit/:id", jwt.Protected, style.StyleEditPost)
	app.Get("/monitor", jwt.Protected, core.Monitor)

	v1 := app.Group("/api")
	v1.Get("/style/:id.user.css", api.GetStyleSource)
	v1.Get("/style/:id", api.GetStyleDetails)
	v1.Get("/styles", api.GetStyleIndex)

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
