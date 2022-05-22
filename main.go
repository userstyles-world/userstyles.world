package main

import (
	"embed"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/compress"
	"github.com/gofiber/fiber/v2/middleware/filesystem"
	"github.com/gofiber/fiber/v2/middleware/limiter"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/pprof"

	"userstyles.world/handlers/api"
	"userstyles.world/handlers/core"
	jwtware "userstyles.world/handlers/jwt"
	oauthprovider "userstyles.world/handlers/oauthProvider"
	"userstyles.world/handlers/style"
	"userstyles.world/handlers/user"
	"userstyles.world/modules/cache"
	"userstyles.world/modules/config"
	database "userstyles.world/modules/database/init"
	"userstyles.world/modules/log"
	"userstyles.world/modules/templates"
	"userstyles.world/modules/util/httputil"
	"userstyles.world/search"
	"userstyles.world/services/cron"
	"userstyles.world/utils"
)

var (
	//go:embed views/*
	views embed.FS

	//go:embed static/*
	static embed.FS
)

func main() {
	log.Initialize()
	cache.Initialize()
	utils.InitalizeCrypto()
	utils.InitializeValidator()
	database.Initialize()
	cron.Initialize()
	search.Initialize()

	app := fiber.New(fiber.Config{
		Views:       templates.New(views, "views"),
		ViewsLayout: "layouts/main",
		ProxyHeader: httputil.ProxyHeader(config.Production),
		JSONEncoder: utils.JSONEncoder,
	})

	if !config.Production {
		app.Use(logger.New())
	}

	if config.Production {
		app.Use(core.HSTSMiddleware)
		app.Use(core.CSPMiddleware)
		app.Use(limiter.New(limiter.Config{
			Max:               350,
			Expiration:        time.Second * 60,
			LimiterMiddleware: limiter.SlidingWindow{},
		}))
	}
	app.Use(compress.New())
	app.Use(jwtware.New("user", jwtware.NormalJWTSigning))

	if config.PerformanceMonitor {
		app.Use(pprof.New())
	}

	// Mount routes.
	core.Routes(app)
	user.Routes(app)
	style.Routes(app)
	api.Routes(app)
	oauthprovider.Routes(app)

	// Embed static files.
	fsys, err := httputil.EmbedFS(static, "static", config.Production)
	if err != nil {
		log.Warn.Fatal(err)
	}
	app.Use(filesystem.New(filesystem.Config{
		MaxAge: 2 * int(time.Hour.Seconds()),
		Root:   fsys,
	}))

	// Fallback route.
	app.Use(core.NotFound)

	log.Warn.Fatal(app.Listen(config.Port))
}
