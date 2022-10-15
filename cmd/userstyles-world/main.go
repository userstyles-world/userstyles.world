package main

import (
	"os"
	"os/signal"
	"runtime/debug"
	"syscall"
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
	"userstyles.world/modules/cron"
	database "userstyles.world/modules/database/init"
	"userstyles.world/modules/images"
	"userstyles.world/modules/log"
	"userstyles.world/modules/search"
	"userstyles.world/modules/templates"
	"userstyles.world/modules/util/httputil"
	"userstyles.world/utils"
	"userstyles.world/web"
)

func main() {
	log.Initialize()
	cache.Initialize()
	images.CheckVips()
	utils.InitalizeCrypto()
	utils.InitializeValidator()
	database.Initialize()
	cron.Initialize()
	search.Initialize()

	app := fiber.New(fiber.Config{
		Views:       templates.New(web.ViewsDir),
		ViewsLayout: "layouts/main",
		ProxyHeader: httputil.ProxyHeader(config.Production),
		JSONEncoder: utils.JSONEncoder,
		IdleTimeout: 5 * time.Second,
	})

	if !config.Production {
		app.Use(logger.New())
	}

	if config.Production {
		app.Use(core.HSTSMiddleware)
		app.Use(core.CSPMiddleware)
		app.Use(limiter.New(limiter.Config{
			Max:               400,
			Expiration:        time.Minute,
			LimiterMiddleware: limiter.FixedWindow{},
		}))
	}
	app.Use(compress.New())
	app.Use(jwtware.New("user", jwtware.NormalJWTSigning))

	if config.PerformanceMonitor {
		app.Use(pprof.New(pprof.Config{
			Next: func(c *fiber.Ctx) bool {
				u, _ := jwtware.User(c)
				return !u.IsAdmin()
			},
		}))

		app.Get("/debug/free", func(c *fiber.Ctx) error {
			u, _ := jwtware.User(c)
			if !u.IsAdmin() {
				return c.Next()
			}

			debug.FreeOSMemory()
			return c.Redirect("/debug/pprof")
		})
	}

	// Mount routes.
	core.Routes(app)
	user.Routes(app)
	style.Routes(app)
	api.Routes(app)
	oauthprovider.Routes(app)

	// Embed static files.
	app.Use(filesystem.New(filesystem.Config{
		MaxAge: 2 * int(time.Hour.Seconds()),
		Root:   web.StaticDir,
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
	_ = database.Close()
	_ = search.StyleIndex.Close()
	log.Info.Printf("Done in %s.\n", time.Since(t))
}
