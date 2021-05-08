package api

import (
	"fmt"
	"strconv"
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

type JWTAPIUser struct {
	models.APIUser
	StyleID uint
}

func APIUser(c *fiber.Ctx) (*JWTAPIUser, bool) {
	s := MapClaim(c)
	u := &JWTAPIUser{}

	if s == nil {
		return u, false
	}

	// Just make sure it's the real deal.
	fUserID, ok := s["userID"].(float64)
	if !ok {
		return u, false
	}
	userID := strconv.Itoa(int(fUserID))

	user, err := models.FindUserByID(database.DB, userID)
	if err != nil || user.ID == 0 {
		return u, false
	}

	u.Username = user.Username
	u.Email = user.Email
	u.ID = uint(fUserID)
	u.Role = user.Role

	// As these are "optional" we need to check them first.
	if Scopes, ok := s["scopes"].(string); ok {
		u.Scopes = strings.Split(Scopes, ",")
	}
	if StyleID, ok := s["styleID"].(float64); ok {
		u.StyleID = uint(StyleID)
	}
	return u, true
}
