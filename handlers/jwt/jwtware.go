package jwt

import (
	"fmt"
	"strings"

	"github.com/form3tech-oss/jwt-go"
	"github.com/gofiber/fiber/v2"
	"userstyles.world/config"
)

var (
	signingKey    = []byte(config.JWT_SIGNING_KEY)
	signingMethod = "HS512"
)

func KeyFuncion(t *jwt.Token) (interface{}, error) {
	if t.Method.Alg() != signingMethod {
		return nil, fmt.Errorf("unexpected jwt signing method=%v", t.Header["alg"])
	}
	return signingKey, nil
}

func New() fiber.Handler {
	extractors := []func(c *fiber.Ctx) (string, bool){jwtFromCookie(fiber.HeaderAuthorization), jwtFromHeader(fiber.HeaderAuthorization)}

	return func(c *fiber.Ctx) error {
		var auth string
		var ok bool

		for _, extractor := range extractors {
			auth, ok = extractor(c)
			if auth != "" && ok {
				break
			}
		}

		if !ok {
			return c.Next()
		}

		token, err := jwt.Parse(auth, KeyFuncion)

		if err == nil && token.Valid {
			// Store user information from token into context.
			c.Locals("user", token)
			return c.Next()
		}
		return c.Next()
	}
}

func jwtFromHeader(header string) func(c *fiber.Ctx) (string, bool) {
	return func(c *fiber.Ctx) (string, bool) {
		auth := c.Get(header)
		l := len("Bearer")
		if len(auth) > l+1 && strings.EqualFold(auth[:l], "Bearer") {
			return auth[l+1:], true
		}
		return "", false
	}
}

func jwtFromCookie(name string) func(c *fiber.Ctx) (string, bool) {
	return func(c *fiber.Ctx) (string, bool) {
		token := c.Cookies(name)
		if token == "" {
			return "", false
		}
		return token, true
	}
}
