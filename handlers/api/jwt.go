package api

import (
	"fmt"
	"strings"

	"github.com/form3tech-oss/jwt-go"
	"github.com/gofiber/fiber/v2"

	"userstyles.world/database"
	jwtware "userstyles.world/handlers/jwt"
	"userstyles.world/models"
	"userstyles.world/utils"
)

var parseJWT = jwtware.New("apiUser", func(t *jwt.Token) (interface{}, error) {
	if t.Method.Alg() != jwtware.SigningMethod {
		return nil, fmt.Errorf("unexpected jwt signing method=%v", t.Header["alg"])
	}
	return utils.OAuthPSigningKey, nil
})

func ProtectedAPI(c *fiber.Ctx) error {
	parseJWT(c)
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
	user, ok := c.Locals("apiUser").(*jwt.Token)
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

	user, err := models.FindUserByName(database.DB, s["username"].(string))

	if err != nil || user.ID == 0 {
		return u, false
	}

	// Type assertion will convert interface{} to other types.
	u.Username = user.Username
	u.Email = user.Email
	u.ID = user.ID
	u.Role = user.Role
	u.Scopes = strings.Split(s["scopes"].(string), ",")

	return u, true
}
