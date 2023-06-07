// Package api provides API endpoints.
package api

import (
	"fmt"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/gofiber/adaptor/v2"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/filesystem"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"userstyles.world/modules/cache"
	"userstyles.world/modules/config"
)

// FastRoutes provides auth-less routes for Fiber's router.
func FastRoutes(app *fiber.App) {
	r := app.Group("/api")
	r.Get("/health", GetHealth)
	// sftodo: sort out etags
	//r.Head("/style/:id.user.css", GetStyleEtag)
	//r.Get("/style/:id.user.css", GetStyleSource)
	r.Get("/metrics", adaptor.HTTPHandler(promhttp.Handler()))

	// sftodo: should and how we specify that its Get requtest?
	r1 := r.Group("/style")
	r1.Use(func(c *fiber.Ctx) error {
		if !strings.HasSuffix(c.Path(), ".user.css") {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"message": "Invalid requtest",
			})
		}

		id := ""
		// sftodo: refactor
		if strings.HasSuffix(c.Path(), ".user.css") {
			id = strings.TrimSuffix(c.Path(), ".user.css")
		} else if strings.HasSuffix(c.Path(), ".user.styl") {
			id = strings.TrimSuffix(c.Path(), ".user.styl")
		} else {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"message": "Invalid requtest",
			})
		}
		id = strings.TrimPrefix(id, "/api/style/")

		fmt.Printf("path: %v"+" id: "+id+"\n", c.Path())
		if _, err := strconv.Atoi(id); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"message": "Invalid requtest",
			})
		}
		cache.InstallStats.Add(c.IP() + " " + id)

		// sftodo: add check for existence of the style with such ID

		return c.Next()
	})
	r1.Use(filesystem.New(filesystem.Config{
		Root:   http.FS(os.DirFS(config.StyleDir)),
		Browse: true,
	}))
}

// Routes provides routes for Fiber's router.
func Routes(app *fiber.App) {
	r := app.Group("/api", ParseAPIJWT)
	//r.Get("/style/:id", GetStyleDetails) // sftodo: figure difference with FastRoutes
	r.Get("/style/stats/:id/:type?", GetStyleStats)
	r.Get("/index/:format?", GetStyleIndex)
	r.Get("/search/:query", GetSearchResult)
	r.Get("/callback/:rcode", CallbackGet)
	r.Get("/user", ProtectedAPI, UserGet)
	r.Get("/user/:identifier", SpecificUserGet)
	r.Get("/styles", ProtectedAPI, StylesGet)
	r.Post("/style/new", ProtectedAPI, NewStyle)
	r.Post("/style/:id", ProtectedAPI, StylePost)
	r.Delete("/style/:id", ProtectedAPI, DeleteStyle)
	r.Get("/style", ProtectedAPI, StyleGet)
}
