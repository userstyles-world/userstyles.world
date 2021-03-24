package jwt

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/form3tech-oss/jwt-go"
	"github.com/gofiber/fiber/v2"
	"userstyles.world/config"
)

type Config struct {
	SigningKey interface{}

	SigningMethod string

	Claims jwt.Claims

	keyFunc jwt.Keyfunc
}

func New() fiber.Handler {
	// Init config
	var cfg Config
	cfg.SigningKey = []byte(config.JWT_SIGNING_KEY)
	cfg.SigningMethod = "HS512"
	cfg.Claims = jwt.MapClaims{}
	cfg.keyFunc = func(t *jwt.Token) (interface{}, error) {
		if t.Method.Alg() != cfg.SigningMethod {
			return nil, fmt.Errorf("unexpected jwt signing method=%v", t.Header["alg"])
		}
		return cfg.SigningKey, nil
	}

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
		token := new(jwt.Token)
		var err error

		if _, ok := cfg.Claims.(jwt.MapClaims); ok {
			token, err = jwt.Parse(auth, cfg.keyFunc)
		} else {
			t := reflect.ValueOf(cfg.Claims).Type().Elem()
			claims := reflect.New(t).Interface().(jwt.Claims)
			token, err = jwt.ParseWithClaims(auth, claims, cfg.keyFunc)
		}
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
