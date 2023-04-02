// Package user provides user endpoints.
package user

import (
	"github.com/gofiber/fiber/v2"

	jwtware "userstyles.world/handlers/jwt"
)

// Routes provides routes for Fiber's router.
func Routes(app *fiber.App) {
	r := app.Group("/")
	r.Get("/login", LoginGet)
	r.Post("/login", LoginPost)
	r.Get("/register", RegisterGet)
	r.Post("/register", RegisterPost)
	r.Get("/oauth/:type", AuthLoginGet)
	r.Get("/verify/:key", VerifyGet)
	r.Get("/recover", RecoverGet)
	r.Post("/recover", RecoverPost)
	r.Get("/reset/:key", ResetGet)
	r.Post("/reset/:key", ResetPost)
	r.Get("/user/:name", Profile)
	r.Get("~:name", Profile)
	r.Get("/logout", jwtware.Protected, Logout)
	r.Get("/account", jwtware.Protected, Account)
	r.Post("/account/:form", jwtware.Protected, EditAccount)
	r.Get("/user/ban/:id", jwtware.Protected, Ban)
	r.Post("/user/ban/:id", jwtware.Protected, ConfirmBan)
	r.Get("/user/delete/:id", jwtware.Protected, DeleteGet)
	r.Post("/user/delete/:id", jwtware.Protected, DeletePost)
}
