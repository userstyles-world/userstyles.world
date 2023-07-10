package main

import (
	"net/http"
	"os"
	"os/signal"
	"runtime/debug"
	"syscall"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/filesystem"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/pprof"

	"userstyles.world/handlers/api"
	"userstyles.world/handlers/core"
	jwtware "userstyles.world/handlers/jwt"
	oauthprovider "userstyles.world/handlers/oauthProvider"
	"userstyles.world/handlers/review"
	"userstyles.world/handlers/style"
	"userstyles.world/handlers/user"
	"userstyles.world/modules/cache"
	"userstyles.world/modules/config"
	"userstyles.world/modules/cron"
	database "userstyles.world/modules/database/init"
	"userstyles.world/modules/email"
	"userstyles.world/modules/images"
	"userstyles.world/modules/log"
	"userstyles.world/modules/templates"
	"userstyles.world/modules/util"
	"userstyles.world/modules/validator"
	"userstyles.world/web"
)

func main() {
	log.Initialize()
	cache.Initialize()
	images.CheckVips()
	util.InitCrypto()
	validator.Init()
	database.Initialize()
	cron.Initialize()

	app := fiber.New(fiber.Config{
		Views:       templates.New(http.FS(web.ViewsDir)),
		ViewsLayout: "layouts/main",
		ProxyHeader: config.ProxyRealIP,
		JSONEncoder: util.JSONEncoder,
		IdleTimeout: 5 * time.Second,

		// TODO: Explore using this more.
		PassLocalsToViews: true,
	})

	email.SetRenderer(app)

	if !config.Production {
		app.Use(logger.New())
	}

	api.FastRoutes(app)

	app.Use(core.HSTSMiddleware)
	app.Use(core.CSPMiddleware)
	app.Use(core.FlagsMiddleware)
	app.Use(jwtware.New("user", jwtware.NormalJWTSigning))

	if config.PerformanceMonitor {
		perf := app.Group("/debug")
		perf.Use(jwtware.Admin)
		perf.Use(pprof.New())
		perf.Get("/free", func(c *fiber.Ctx) error {
			debug.FreeOSMemory()
			return c.Redirect("/debug/pprof")
		})
	}

	// Mount routes.
	core.Routes(app)
	user.Routes(app)
	style.Routes(app)
	review.Routes(app)
	api.Routes(app)
	oauthprovider.Routes(app)

	// Embed static files.
	app.Use(filesystem.New(filesystem.Config{
		MaxAge: 2 * int(time.Hour.Seconds()),
		Root:   http.FS(web.StaticDir),
	}))

	// TODO: Investigate how to "truly" inline sourcemaps in Sass.
	if !config.Production {
		app.Static("/scss", "web/scss")
	}

	// Fallback route.
	app.Use(core.NotFound)

	go func() {
		if err := app.Listen(config.Port); err != nil {
			log.Warn.Fatal(err)
		}
	}()

	// Block and listen.
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	<-c

	// Close everything and exit.
	log.Info.Println("Shutting down...")
	t := time.Now()
	_ = app.Shutdown()
	cache.Code.Close()
	cache.InstallStats.Close()
	cache.ViewStats.Close()
	cache.SaveStore()
	_ = database.Close()
	log.Info.Printf("Done in %s.\n", time.Since(t))
}
