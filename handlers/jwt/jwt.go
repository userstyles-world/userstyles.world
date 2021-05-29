package jwt

import (
	"fmt"
	"net/url"

	"github.com/form3tech-oss/jwt-go"
	lib "github.com/form3tech-oss/jwt-go"
	"github.com/gofiber/fiber/v2"

	"userstyles.world/models"
	"userstyles.world/utils"
)

var NormalJWTSigning = func(t *jwt.Token) (interface{}, error) {
	if t.Method.Alg() != SigningMethod {
		return nil, fmt.Errorf("unexpected jwt signing method=%v", t.Header["alg"])
	}
	return JWTSigningKey, nil
}

var Protected = func(c *fiber.Ctx) error {
	if _, ok := User(c); !ok {
		redirectUri := utils.UnsafeString(c.Request().URI().Path())
		if c.Context().QueryArgs().Len() != 0 {
			redirectUri += "?" + c.Context().QueryArgs().String()
		}

		return c.Redirect("/login?after_login=" + url.QueryEscape(redirectUri))
	}
	return c.Next()
}

func MapClaim(c *fiber.Ctx) lib.MapClaims {
	user, ok := c.Locals("user").(*lib.Token)
	if !ok {
		return nil
	}
	claims := user.Claims.(lib.MapClaims)

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
