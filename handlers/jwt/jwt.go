package jwt

import (
	"net/url"

	lib "github.com/form3tech-oss/jwt-go"
	"github.com/gofiber/fiber/v2"

	"userstyles.world/models"
	"userstyles.world/modules/errors"
	"userstyles.world/utils"
)

var NormalJWTSigning = func(t *lib.Token) (interface{}, error) {
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
