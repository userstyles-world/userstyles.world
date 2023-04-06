// Package style provides style endpoints.
package style

import (
	"github.com/gofiber/fiber/v2"

	jwtware "userstyles.world/handlers/jwt"
	"userstyles.world/modules/config"
)

// Routes provides routes for Fiber's router.
func Routes(app *fiber.App) {
	r := app.Group("/")
	r.Get("/explore", GetExplore)
	r.Get("/style/:id/:name?", GetStylePage)
	r.Get("/add", jwtware.Protected, CreateGet)
	r.Post("/add", jwtware.Protected, CreatePost)
	r.Get("/delete/:id", jwtware.Protected, DeleteGet)
	r.Post("/delete/:id", jwtware.Protected, DeletePost)
	r.Get("/import", jwtware.Protected, ImportGet)
	r.Post("/import", jwtware.Protected, ImportPost)
	r.Get("/edit/:id", jwtware.Protected, EditGet)
	r.Post("/edit/:id", jwtware.Protected, EditPost)
	r.Get("/mirror/:id", jwtware.Protected, Mirror)
	r.Get("/styles/promote/:id", jwtware.Protected, Promote)
	r.Get("/styles/ban/:id", jwtware.Protected, BanGet)
	r.Post("/styles/ban/:id", jwtware.Protected, BanPost)
	r.Get("/styles/review/:id", jwtware.Protected, ReviewGet)
	r.Post("/styles/review/:id", jwtware.Protected, ReviewPost)
	r.Static("/preview", config.PublicDir, fiber.Static{
		MaxAge: 2678400, // 1 month
	})
}
