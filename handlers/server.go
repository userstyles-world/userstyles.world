package handlers

import (
	"html/template"
	"log"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/compress"
	"github.com/gofiber/fiber/v2/middleware/filesystem"
	"github.com/gofiber/fiber/v2/middleware/limiter"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/template/html"
	"github.com/markbates/pkger"
	"github.com/microcosm-cc/bluemonday"
	"github.com/russross/blackfriday/v2"

	"userstyles.world/config"
	"userstyles.world/handlers/api"
	"userstyles.world/handlers/core"
	"userstyles.world/handlers/jwt"
	jwtware "userstyles.world/handlers/jwt"
	oauthprovider "userstyles.world/handlers/oauthProvider"
	"userstyles.world/handlers/style"
	"userstyles.world/handlers/user"
	"userstyles.world/utils"
)

// TODO: Refactor this as a separate package.
func renderEngine() *html.Engine {
	engine := html.NewFileSystem(pkger.Dir("/views"), ".html")

	engine.AddFunc("Markdown", func(s string) template.HTML {
		// Generate Markdown then sanitize it before returning HTML.
		gen := blackfriday.Run(
			[]byte(s),
			blackfriday.WithExtensions(blackfriday.HardLineBreak),
			blackfriday.WithExtensions(blackfriday.AutoHeadingIDs),
			blackfriday.WithExtensions(blackfriday.Extensions(blackfriday.CodeBlock)),
		)
		out := bluemonday.UGCPolicy().SanitizeBytes(gen)

		return template.HTML(out)
	})

	engine.AddFunc("GitCommit", func() template.HTML {
		if !config.Production {
			return template.HTML("dev")
		}

		return template.HTML(config.GIT_COMMIT)
	})

	engine.AddFunc("Date", func(time time.Time) template.HTML {
		return template.HTML(time.Format("January 02, 2006 15:04"))
	})

	if !config.Production {
		engine.Reload(true)
	}

	return engine
}

// Get proper IP depending on the environment.
func proxyHeader() (s string) {
	if config.Production {
		s = "X-Real-IP"
	}

	return s
}

func Initialize() {
	app := fiber.New(fiber.Config{
		Views:       renderEngine(),
		ProxyHeader: proxyHeader(),
		JSONEncoder: utils.JSONEncoder,
	})

	if !config.Production {
		app.Use(logger.New())
	}

	app.Use(core.HSTSMiddleware)
	app.Use(compress.New())
	if config.Production {
		app.Use(limiter.New(limiter.Config{Max: 300}))
	}
	app.Use(jwtware.New("user", jwtware.NormalJWTSigning))

	app.Get("/", core.Home)
	app.Get("/search", core.Search)
	app.Get("/login", user.LoginGet)
	app.Post("/login", user.LoginPost)
	app.Get("/register", user.RegisterGet)
	app.Post("/register", user.RegisterPost)
	app.Get("/oauth/:type", user.AuthLoginGet)
	app.Get("/verify/:key", user.VerifyGet)
	app.Get("/recover", user.RecoverGet)
	app.Post("/recover", user.RecoverPost)
	app.Get("/reset/:key", user.ResetGet)
	app.Post("/reset/:key", user.ResetPost)
	app.Get("/explore", style.GetExplore)
	app.Get("/style/:id/:name?", style.GetStylePage)
	app.Get("/user/:name", user.Profile)
	app.Get("~:name", user.Profile)
	app.Get("/docs/:document", core.GetDocs)

	app.Get("/logout", jwtware.Protected, user.Logout)
	app.Get("/account", jwtware.Protected, user.Account)
	app.Post("/account", jwtware.Protected, user.EditAccount)
	app.Get("/add", jwtware.Protected, style.CreateGet)
	app.Post("/add", jwtware.Protected, style.CreatePost)
	app.Get("/delete/:id", jwt.Protected, style.DeleteGet)
	app.Post("/delete/:id", jwt.Protected, style.DeletePost)
	app.Get("/import", jwtware.Protected, style.ImportGet)
	app.Post("/import", jwtware.Protected, style.ImportPost)
	app.Get("/edit/:id", jwtware.Protected, style.EditGet)
	app.Post("/edit/:id", jwtware.Protected, style.EditPost)
	app.Post("/style/:id/promote", jwtware.Protected, style.Promote)
	app.Get("/oauth_settings/:id?", jwtware.Protected, oauthprovider.OAuthSettingsGet)
	app.Post("/oauth_settings/:id?", jwtware.Protected, oauthprovider.OAuthSettingsPost)
	app.Get("/monitor", jwtware.Protected, core.Monitor)

	v1 := app.Group("/api", api.ParseAPIJWT)
	v1.Head("/style/:id.user.css", api.GetStyleEtag)
	v1.Get("/style/:id.user.css", api.GetStyleSource)
	v1.Get("/style/:id", api.GetStyleDetails)
	v1.Get("/style/preview/:id", api.GetPreviewScreenshot)
	v1.Get("/index/:format?", api.GetStyleIndex)
	v1.Get("/search/:query", api.GetSearchResult)
	v1.Get("/callback/:rcode", api.CallbackGet)
	v1.Get("/user", api.ProtectedAPI, api.UserGet)
	v1.Get("/user/:identifier", api.SpecificUserGet)
	v1.Get("/styles", api.ProtectedAPI, api.StylesGet)
	v1.Post("/style/new", api.ProtectedAPI, api.NewStyle)
	v1.Post("/style/:id", api.ProtectedAPI, api.StylePost)
	v1.Delete("/style/:id", api.ProtectedAPI, api.DeleteStyle)
	v1.Get("/style", api.ProtectedAPI, api.StyleGet)

	oauthV1 := app.Group("/api/oauth")
	oauthV1.Get("/authorize", jwtware.Protected, oauthprovider.AuthorizeGet)
	oauthV1.Get("/authorize_style", jwtware.Protected, oauthprovider.AuthorizeStyleGet)
	oauthV1.Post("/authorize_style", jwtware.Protected, oauthprovider.AuthorizeStylePost)
	oauthV1.Get("/authorize_style/new", jwtware.Protected, oauthprovider.AuthorizeStyleNewPost)
	oauthV1.Post("/authorize_style/new", jwtware.Protected, oauthprovider.AuthorizeStyleNewPost)
	oauthV1.Post("/authorize/:id/:token", jwtware.Protected, oauthprovider.AuthorizePost)
	oauthV1.Post("/access_token", oauthprovider.AccessTokenPost)

	// Allows assets to be reloaded in dev mode.
	// That means, they're not embedded into executable file.
	if !config.Production {
		app.Static("/", "/static")
	}

	app.Use("/", filesystem.New(filesystem.Config{
		MaxAge: int(time.Hour) * 2,
		Root:   pkger.Dir("/static"),
	}))
	app.Use(core.NotFound)

	log.Fatal(app.Listen(config.PORT))
}
