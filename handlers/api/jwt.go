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

var ParseAPIJWT = jwtware.New("apiUser", func(t *jwt.Token) (interface{}, error) {
	if t.Method.Alg() != jwtware.SigningMethod {
		return nil, fmt.Errorf("unexpected jwt signing method=%v", t.Header["alg"])
	}
	return utils.OAuthPSigningKey, nil
})

func ProtectedAPI(c *fiber.Ctx) error {
	if _, ok := APIUser(c); !ok {
		return c.Status(401).
			JSON(fiber.Map{
				"data": "You need to provide an access_token within the Authorization header.",
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

func APIUser(c *fiber.Ctx) (*models.APIUser, bool) {
	s := MapClaim(c)
	u := &models.APIUser{}

	if s == nil {
		return u, false
	}

	user, err := models.FindUserByID(database.DB, fmt.Sprintf("%d", uint(s["userID"].(float64))))

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
