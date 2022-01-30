// Package oauthprovider provides OAuth Provider endpoints.
package oauthprovider

import (
	"github.com/gofiber/fiber/v2"

	jwtware "userstyles.world/handlers/jwt"
)

// Routes provides routes for Fiber's router.
func Routes(app *fiber.App) {
	r := app.Group("/api/oauth")
	r.Get("/auth", jwtware.Protected, AuthorizeGet)
	r.Get("/settings/:id?", jwtware.Protected, OAuthSettingsGet)
	r.Post("/settings/:id?", jwtware.Protected, OAuthSettingsPost)
	r.Get("/style/link", jwtware.Protected, OAuthStyleGet)
	r.Post("/style/link", jwtware.Protected, OAuthStylePost)
	r.Get("/style/new", jwtware.Protected, OAuthStyleNewPost)
	r.Post("/style/new", jwtware.Protected, OAuthStyleNewPost)
	r.Post("/auth/:id/:token", jwtware.Protected, AuthPost)
	r.Post("/token", TokenPost)
}
