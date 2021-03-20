package jwt

import (
	"log"

	jwt "github.com/form3tech-oss/jwt-go"
	"github.com/gofiber/fiber/v2"
	jwtware "github.com/gofiber/jwt/v2"

	"userstyles.world/config"
	"userstyles.world/models"
)

var Protected = jwtware.New(jwtware.Config{
	SigningMethod: "HS512",
	SigningKey:    []byte(config.JWT_SIGNING_KEY),
	TokenLookup:   "cookie:" + fiber.HeaderAuthorization,
	ErrorHandler:  loginRequired,
})

var NoLoggedInUsers = jwtware.New(jwtware.Config{
	SigningMethod: "HS512",
	SigningKey:    []byte(config.JWT_SIGNING_KEY),
	TokenLookup:   "cookie:" + fiber.HeaderAuthorization,
	ErrorHandler:  loginRequired,
	SuccessHandler: func(c *fiber.Ctx) error {
		log.Printf("User %d has set valid JWT Token, redirecting.", User(c).ID)
		return c.Redirect("/account", fiber.StatusSeeOther)
	},
})

var Everyone = jwtware.New(jwtware.Config{
	SigningMethod: "HS512",
	SigningKey:    []byte(config.JWT_SIGNING_KEY),
	TokenLookup:   "cookie:" + fiber.HeaderAuthorization,
	ErrorHandler: func(c *fiber.Ctx, e error) error {
		return c.Next()
	},
})

func loginRequired(c *fiber.Ctx, e error) error {
	if c.Route().Path == "/login" {
		return c.Next()
	}
	c.Status(fiber.StatusUnauthorized)
	return c.Render("login", fiber.Map{
		"Title": "Login is required",
		"Error": "You must log in to do this action.",
	})
}

func MapClaim(c *fiber.Ctx) jwt.MapClaims {
	user, ok := c.Locals("user").(*jwt.Token)
	if !ok {
		return nil
	}
	claims := user.Claims.(jwt.MapClaims)

	return claims
}

func User(c *fiber.Ctx) *models.APIUser {
	s := MapClaim(c)
	u := &models.APIUser{}

	if s == nil {
		return u
	}

	// Type assertion will convert interface{} to other types.
	u.Username = s["name"].(string)
	u.Email = s["email"].(string)
	u.ID = uint(s["id"].(float64))

	return u
}
