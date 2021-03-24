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

	public := app.Group("/", jwt.Everyone)
	public.Get("/", core.Home)
	public.Get("/login", jwt.NoLoggedInUsers, user.LoginGet)
	public.Post("/login", user.LoginPost)
	public.Get("/register", jwt.NoLoggedInUsers, user.RegisterGet)
	public.Post("/register", user.RegisterPost)
	public.Get("/explore", style.GetExplore)
	public.Get("/style/:id", style.GetStyle)
	public.Get("/user/:name", user.Profile)

	protected := app.Group("/", jwt.Protected)
	protected.Post("/logout", user.Logout)
	protected.Get("/account", user.Account)
	protected.Post("/style/:id", style.DeleteByID)
	protected.Get("/add", style.StyleCreateGet)
	protected.Post("/add", style.StyleCreatePost)
	protected.Get("/import", style.StyleImportGet)
	protected.Post("/import", style.StyleImportPost)
	protected.Get("/edit/:id", style.StyleEditGet)
	protected.Post("/edit/:id", style.StyleEditPost)
	protected.Get("/monitor", core.Monitor)

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
