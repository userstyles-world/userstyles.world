package review

import (
	"github.com/gofiber/fiber/v2"

	"userstyles.world/handlers/jwt"
)

// Routes provides routes for Fiber's router.
func Routes(app *fiber.App) {
	r := app.Group("/styles/:s/reviews", jwt.Protected)
	r.Get("/create", createPage)
	r.Post("/create", createForm)
	r.Get("/:r", viewPage)
	r.Get("/:r/edit", editPage)
	r.Post("/:r/edit", editForm)
	r.Get("/:r/delete", deletePage)
	r.Post("/:r/delete", deleteForm)
}
