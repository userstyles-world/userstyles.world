package review

import (
	"github.com/gofiber/fiber/v2"

	"userstyles.world/handlers/jwt"
)

// Routes provides routes for Fiber's router.
func Routes(app *fiber.App) {
	r := app.Group("/styles/review/:id", jwt.Protected)
	r.Get("/", createPage)
	r.Post("/", createForm)
	r.Get("/edit", updatePage)
	r.Post("/edit", updateForm)
	r.Get("/delete", deletePage)
	r.Post("/delete", deleteForm)
}
