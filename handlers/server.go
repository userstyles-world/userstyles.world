package handlers

import (
	"log"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/compress"
	"github.com/gofiber/fiber/v2/middleware/filesystem"
	"github.com/gofiber/fiber/v2/middleware/limiter"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/pprof"
	"github.com/markbates/pkger"

	"userstyles.world/handlers/api"
	"userstyles.world/handlers/core"
	jwtware "userstyles.world/handlers/jwt"
	oauthprovider "userstyles.world/handlers/oauthProvider"
	"userstyles.world/handlers/style"
	"userstyles.world/handlers/user"
	"userstyles.world/modules/config"
	"userstyles.world/modules/templates"
	"userstyles.world/utils"
)

// Get proper IP depending on the environment.
func proxyHeader() (s string) {
	if config.Production {
		s = "X-Real-IP"
	}

	return s
}

func Initialize() {
	app := fiber.New(fiber.Config{
		Views:       templates.New(),
		ViewsLayout: "layouts/main",
		ProxyHeader: proxyHeader(),
		JSONEncoder: utils.JSONEncoder,
	})

	if !config.Production {
		app.Use(logger.New())
	}

	app.Use(core.HSTSMiddleware)
	app.Use(core.CSPMiddleware)
	app.Use(compress.New())
	if config.Production {
		app.Use(limiter.New(limiter.Config{Max: 600}))
	}
	app.Use(jwtware.New("user", jwtware.NormalJWTSigning))

	if config.PerformanceMonitor {
		app.Use(pprof.New())
		app.Get("/monitor", jwtware.Protected, core.Monitor)
	}

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
	app.Get("/docs/:document?", core.GetDocs)
	app.Get("/modlog", core.GetModLog)
	app.Get("/link/:site", core.GetLinkedSite)

	app.Get("/logout", jwtware.Protected, user.Logout)
	app.Get("/account", jwtware.Protected, user.Account)
	app.Post("/account", jwtware.Protected, user.EditAccount)
	app.Get("/add", jwtware.Protected, style.CreateGet)
	app.Post("/add", jwtware.Protected, style.CreatePost)
	app.Get("/delete/:id", jwtware.Protected, style.DeleteGet)
	app.Post("/delete/:id", jwtware.Protected, style.DeletePost)
	app.Get("/import", jwtware.Protected, style.ImportGet)
	app.Post("/import", jwtware.Protected, style.ImportPost)
	app.Get("/edit/:id", jwtware.Protected, style.EditGet)
	app.Post("/edit/:id", jwtware.Protected, style.EditPost)
	app.Post("/style/:id/promote", jwtware.Protected, style.Promote)
	app.Get("/styles/ban/:id", jwtware.Protected, style.BanGet)
	app.Post("/styles/ban/:id", jwtware.Protected, style.BanPost)
	app.Get("/styles/review/:id", jwtware.Protected, style.ReviewGet)
	app.Post("/styles/review/:id", jwtware.Protected, style.ReviewPost)
	app.Get("/oauth_settings/:id?", jwtware.Protected, oauthprovider.OAuthSettingsGet)
	app.Post("/oauth_settings/:id?", jwtware.Protected, oauthprovider.OAuthSettingsPost)
	app.Get("/user/ban/:id", jwtware.Protected, user.Ban)
	app.Post("/user/ban/:id", jwtware.Protected, user.ConfirmBan)
	app.Get("/dashboard", jwtware.Protected, core.Dashboard)

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
	oauthV1.Get("/auth", jwtware.Protected, oauthprovider.AuthorizeGet)
	oauthV1.Get("/style/link", jwtware.Protected, oauthprovider.OAuthStyleGet)
	oauthV1.Post("/style/link", jwtware.Protected, oauthprovider.OAuthStylePost)
	oauthV1.Get("/style/new", jwtware.Protected, oauthprovider.OAuthStyleNewPost)
	oauthV1.Post("/style/new", jwtware.Protected, oauthprovider.OAuthStyleNewPost)
	oauthV1.Post("/auth/:id/:token", jwtware.Protected, oauthprovider.AuthPost)
	oauthV1.Post("/token", oauthprovider.TokenPost)

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

	log.Fatal(app.Listen(config.Port))
}
