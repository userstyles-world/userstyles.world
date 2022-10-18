package jwt

import (
	"net/url"

	"github.com/gofiber/fiber/v2"
	lib "github.com/golang-jwt/jwt"

	"userstyles.world/models"
	"userstyles.world/modules/errors"
	"userstyles.world/utils"
)

var NormalJWTSigning = func(t *lib.Token) (any, error) {
	if t.Method.Alg() != SigningMethod {
		return nil, errors.UnexpectedSigningMethod(t.Method.Alg())
	}
	return JWTSigningKey, nil
}

var Protected = func(c *fiber.Ctx) error {
	if _, ok := User(c); !ok {
		redirectURI := utils.UnsafeString(c.Request().URI().Path())
		if c.Context().QueryArgs().Len() != 0 {
			redirectURI += "?" + c.Context().QueryArgs().String()
		}

		return c.Redirect("/login?r=" + url.QueryEscape(redirectURI))
	}
	return c.Next()
}

var Admin = func(c *fiber.Ctx) error {
	u, ok := User(c)
	if !ok || !u.IsAdmin() {
		return c.Redirect("/login?r=" + url.QueryEscape(c.Path()))
	}

	return c.Next()
}

func MapClaim(c *fiber.Ctx) lib.MapClaims {
	user, ok := c.Locals("user").(*lib.Token)
	if !ok {
		return nil
	}
	claims, ok := user.Claims.(lib.MapClaims)
	if !ok {
		return nil
	}
	return claims
}

func User(c *fiber.Ctx) (*models.APIUser, bool) {
	s := MapClaim(c)
	u := &models.APIUser{}

	if s == nil {
		return u, false
	}

	// Type assertion will convert interface{} to other types.
	if name, ok := s["name"].(string); ok {
		u.Username = name
	}
	if email, ok := s["email"].(string); ok {
		u.Email = email
	}
	if id, ok := s["id"].(float64); ok {
		u.ID = uint(id)
	}
	if role, ok := s["role"].(float64); ok {
		u.Role = models.Role(role)
	}

	return u, true
}
