package jwt

import (
	jwt "github.com/form3tech-oss/jwt-go"
	"github.com/gofiber/fiber/v2"

	"userstyles.world/models"
)

var Protected = func(c *fiber.Ctx) error {
	if _, ok := User(c); !ok {
		c.Status(fiber.StatusUnauthorized)
		return c.Render("login", fiber.Map{
			"Title": "Login is required",
			"Error": "You must log in to do this action.",
		})
	}
	return c.Next()
}

func MapClaim(c *fiber.Ctx) jwt.MapClaims {
	user, ok := c.Locals("user").(*jwt.Token)
	if !ok {
		return nil
	}
	claims := user.Claims.(jwt.MapClaims)

	return claims
}

func User(c *fiber.Ctx) (*models.APIUser, bool) {
	s := MapClaim(c)
	u := &models.APIUser{}

	if s == nil {
		return u, false
	}

	// Type assertion will convert interface{} to other types.
	u.Username = s["name"].(string)
	if s["email"] != nil {
		u.Email = s["email"].(string)
	}
	u.ID = uint(s["id"].(float64))
	u.Role = models.Role(s["role"].(float64))

	return u, true
}
