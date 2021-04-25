package handlers

import (
	"fmt"
	"html/template"
	"log"
	"time"

	"github.com/form3tech-oss/jwt-go"
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
	jwtware "userstyles.world/handlers/jwt"
	"userstyles.world/handlers/oauth_provider"
	"userstyles.world/handlers/style"
	"userstyles.world/handlers/user"
	"userstyles.world/models"
	"userstyles.world/utils"
)

// TODO: Refactor this as a separate package.
func renderEngine() *html.Engine {
	engine := html.NewFileSystem(pkger.Dir("/views"), ".html")

	engine.AddFunc("Markdown", func(s string) template.HTML {
		// Generate Markdown then sanitize it before returning HTML.
		gen := blackfriday.Run([]byte(s), blackfriday.WithExtensions(blackfriday.HardLineBreak))
		out := bluemonday.UGCPolicy().SanitizeBytes(gen)

		return template.HTML(out)
	})

	engine.AddFunc("GitCommit", func() template.HTML {
		if !config.Production {
			return template.HTML("dev")
		}

		return template.HTML(config.GIT_COMMIT)
	})

	engine.AddFunc("Role", func(i models.Role) template.HTML {
		r := ""
		switch i {
		case 0:
			r = "Regular"
		case 1:
			r = "Moderator"
		case 2:
			r = "Admin"
		}
		return template.HTML(r)
	})

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
		JSONEncoder: utils.JsonEncoder,
	})

	if !config.Production {
		app.Use(logger.New())
	}

	app.Use(compress.New())
	if config.Production {
		app.Use(limiter.New(limiter.Config{Max: 300}))
	}
	app.Use(jwtware.New("user", func(t *jwt.Token) (interface{}, error) {
		if t.Method.Alg() != jwtware.SigningMethod {
			return nil, fmt.Errorf("unexpected jwt signing method=%v", t.Header["alg"])
		}
		return jwtware.JWTSigningKey, nil
	}))

	app.Get("/", core.Home)
	app.Get("/search", core.Search)
	app.Get("/login", user.LoginGet)
	app.Post("/login", user.LoginPost)
	app.Get("/register", user.RegisterGet)
	app.Post("/register", user.RegisterPost)
	app.Get("/oauth_login/:type", user.AuthLoginGet)
	app.Get("/verify/:key", user.VerifyGet)
	app.Get("/recover", user.RecoverGet)
	app.Post("/recover", user.RecoverPost)
	app.Get("/reset/:key", user.ResetGet)
	app.Post("/reset/:key", user.ResetPost)
	app.Get("/explore", style.GetExplore)
	app.Get("/style/:id/:name?", style.GetStyleSlug)
	app.Get("/user/:name", user.Profile)
	app.Get("~:name", user.Profile)

	app.Get("/logout", jwtware.Protected, user.Logout)
	app.Get("/account", jwtware.Protected, user.Account)
	app.Post("/account", jwtware.Protected, user.EditAccount)
	app.Post("/style/:id", jwtware.Protected, style.DeleteByID)
	app.Get("/add", jwtware.Protected, style.StyleCreateGet)
	app.Post("/add", jwtware.Protected, style.StyleCreatePost)
	app.Get("/import", jwtware.Protected, style.StyleImportGet)
	app.Post("/import", jwtware.Protected, style.StyleImportPost)
	app.Get("/edit/:id", jwtware.Protected, style.StyleEditGet)
	app.Post("/edit/:id", jwtware.Protected, style.StyleEditPost)
	app.Post("/style/:id/promote", jwtware.Protected, style.StylePromote)
	app.Get("/oauth_settings/:id?", jwtware.Protected, oauth_provider.OAuthSettingsGet)
	app.Post("/oauth_settings/:id?", jwtware.Protected, oauth_provider.OAuthSettingsPost)
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

	oauthV1 := app.Group("/oauth")
	oauthV1.Get("/authorize", oauth_provider.AuthorizeGet)
	oauthV1.Post("/authorize/:id/:token", jwtware.Protected, oauth_provider.AuthorizePost)
	oauthV1.Post("/access_token", oauth_provider.AccessTokenPost)

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
