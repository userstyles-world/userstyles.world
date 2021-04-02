package handlers

import (
	"html/template"
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/compress"
	"github.com/gofiber/fiber/v2/middleware/filesystem"
	"github.com/gofiber/fiber/v2/middleware/limiter"
	"github.com/gofiber/template/html"
	"github.com/markbates/pkger"
	"github.com/microcosm-cc/bluemonday"
	"github.com/russross/blackfriday/v2"

	"userstyles.world/config"
	"userstyles.world/handlers/api"
	"userstyles.world/handlers/core"
	"userstyles.world/handlers/jwt"
	"userstyles.world/handlers/style"
	"userstyles.world/handlers/user"
)

func Initialize() {
	engine := html.NewFileSystem(pkger.Dir("/views"), ".html")

	// TODO: Refactor this in a separate package.
	engine.AddFunc("Markdown", func(s string) template.HTML {
		gen := blackfriday.Run([]byte(s), blackfriday.WithExtensions(blackfriday.HardLineBreak))
		out := bluemonday.UGCPolicy().SanitizeBytes(gen)
		return template.HTML(out)
	})

	engine.AddFunc("GitCommit", func() template.HTML {
		return template.HTML(config.GIT_COMMIT)
	})

	app := fiber.New(fiber.Config{
		Views:                 engine,
		DisableStartupMessage: true,
	})

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
	app.Get("/verify/:key", user.VerifyGet)
	app.Get("/recover", user.RecoverGet)
	app.Post("/recover", user.RecoverPost)
	app.Get("/reset/:key", user.ResetGet)
	app.Post("/reset/:key", user.ResetPost)
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
	app.Post("/style/:id/promote", jwt.Protected, style.StylePromote)
	app.Get("/monitor", jwt.Protected, core.Monitor)

	v1 := app.Group("/api")
	v1.Get("/style/:id.user.css", api.GetStyleSource)
	v1.Get("/style/:id", api.GetStyleDetails)
	v1.Get("/styles", api.GetStyleIndex)

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
