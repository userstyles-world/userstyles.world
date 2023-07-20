package review

import (
	"github.com/gofiber/fiber/v2"

	"userstyles.world/handlers/jwt"
	"userstyles.world/handlers/middleware"
)

// Routes provides routes for Fiber's router.
func Routes(app *fiber.App) {
	r := app.Group("/styles/:s-:slug/reviews")
	r.Get("/create", jwt.Protected, createPage)
	r.Post("/create", jwt.Protected, createForm)

	r = app.Group("/styles/:s-:slug/reviews/:r", middleware.Alert)
	r.Get("/", viewPage)
	r.Use(jwt.Protected)
	r.Get("/edit", editPage)
	r.Post("/edit", editForm)
	r.Get("/delete", deletePage)
	r.Post("/delete", deleteForm)
	r.Get("/remove", removePage)
	r.Post("/remove", removeForm)
}
