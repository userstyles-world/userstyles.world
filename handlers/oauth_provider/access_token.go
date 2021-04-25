package oauth_provider

import (
	"fmt"
	"strings"

	"github.com/form3tech-oss/jwt-go"
	"github.com/gofiber/fiber/v2"
	"userstyles.world/database"
	"userstyles.world/models"
	"userstyles.world/utils"
)

func AccessTokenPost(c *fiber.Ctx) error {
	clientID, clientSecret, stateQuery, tCode := c.FormValue("client_id"), c.FormValue("client_secret"), c.FormValue("state"), c.FormValue("code")

	if clientID == "" {
		return c.Status(400).
			JSON(fiber.Map{
				"error": "No client_id specified",
			})
	}
	if clientSecret == "" {
		return c.Status(400).
			JSON(fiber.Map{
				"error": "No client_secret specified",
			})
	}
	if tCode == "" {
		return c.Status(400).
			JSON(fiber.Map{
				"error": "No code specified",
			})
	}

	OAuth, err := models.GetOAuthByClientID(database.DB, clientID)
	if err != nil || OAuth.ID == 0 {
		return c.Status(400).
			JSON(fiber.Map{
				"error": "Incorrect client_id specified",
			})
	}
	if OAuth.ClientSecret != clientSecret {
		return c.Status(400).
			JSON(fiber.Map{
				"error": "Incorrect client_secret specified",
			})
	}

	unsealedText, err := utils.DecodePreparedText(tCode, utils.AEAD_OAUTHP)
	if err != nil {
		fmt.Println("Error: Couldn't unseal JWT Token:", err.Error())
		return c.Status(500).
			JSON(fiber.Map{
				"error": "JWT Token error, please notify the admins.",
			})
	}

	token, err := jwt.Parse(unsealedText, utils.OAuthPJwtKeyFunction)
	if err != nil || !token.Valid {
		fmt.Println("Error: Couldn't unseal JWT Token:", err.Error())
		return c.Status(500).
			JSON(fiber.Map{
				"error": "JWT Token error, please notify the admins.",
			})
	}

	claims := token.Claims.(jwt.MapClaims)

	state, userName := claims["state"].(string), claims["userID"].(string)

	if stateQuery != state {
		return c.Status(400).
			JSON(fiber.Map{
				"error": "State doesn't match.",
			})
	}

	user, err := models.FindUserByName(database.DB, userName)
	if err != nil || user.ID == 0 {
		return c.Status(500).
			JSON(fiber.Map{
				"error": "Couldn't find the user that was specified, please notify the admins.",
			})
	}

	jwt, err := utils.NewJWTToken().
		SetClaim("scopes", strings.Join(OAuth.Scopes, ", ")).
		SetClaim("username", user.Username).
		GetSignedString(utils.OAuthPSigningKey)

	if err != nil {
		return c.Status(500).
			JSON(fiber.Map{
				"error": "Couldn't create access_token please notify the admins.",
			})
	}

	switch c.Accepts("application/json", "plain/text") {
	case "application/json":
		return c.JSON(fiber.Map{
			"access_token": jwt,
		})
	case "plain/text":
		return c.SendString(jwt)

	}

	return c.SendString(jwt)

}
