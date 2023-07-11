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
	r.Get("/edit", updatePage)
	r.Post("/edit", updateForm)
	r.Get("/delete", deletePage)
	r.Post("/delete", deleteForm)
}
