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
		return errorMessage(c, 400, "No client_id specified")
	}
	if clientSecret == "" {
		return errorMessage(c, 400, "No client_secret specified")
	}
	if tCode == "" {
		return errorMessage(c, 400, "No code specified")

	}

	OAuth, err := models.GetOAuthByClientID(database.DB, clientID)
	if err != nil || OAuth.ID == 0 {
		return errorMessage(c, 400, "Incorrect client_id specified")
	}
	if OAuth.ClientSecret != clientSecret {
		return errorMessage(c, 400, "Incorrect client_secret specified")
	}

	unsealedText, err := utils.DecodePreparedText(tCode, utils.AEAD_OAUTHP)
	if err != nil {
		fmt.Println("Error: Couldn't unseal JWT Token:", err.Error())
		return errorMessage(c, 500, "JWT Token error, please notify the admins.")
	}

	token, err := jwt.Parse(unsealedText, utils.OAuthPJwtKeyFunction)
	if err != nil || !token.Valid {
		fmt.Println("Error: Couldn't unseal JWT Token:", err.Error())
		return errorMessage(c, 500, "JWT Token error, please notify the admins.")
	}

	claims := token.Claims.(jwt.MapClaims)

	state, userID := claims["state"].(string), uint(claims["userID"].(float64))

	styleID, ok := claims["styleID"].(string)
	if !ok {
		styleID = ""
	}

	if stateQuery != state {
		return errorMessage(c, 500, "State doesn't match.")
	}

	user, err := models.FindUserByID(database.DB, fmt.Sprintf("%d", userID))
	if err != nil || user.ID == 0 {
		return errorMessage(c, 500, "Couldn't find the user that was specified, please notify the admins.")
	}

	var jwt string

	if styleID != "" {
		jwt, err = utils.NewJWTToken().
			SetClaim("styleID", styleID).
			SetClaim("userID", user.ID).
			GetSignedString(utils.OAuthPSigningKey)
	} else {
		jwt, err = utils.NewJWTToken().
			SetClaim("scopes", strings.Join(OAuth.Scopes, ",")).
			SetClaim("userID", user.ID).
			GetSignedString(utils.OAuthPSigningKey)
	}

	if err != nil {
		return errorMessage(c, 500, "Couldn't create access_token please notify the admins.")
	}

	if c.Accepts("application/json", "text/plain ") == "application/json" {
		return c.JSON(fiber.Map{
			"access_token": jwt,
			"token_type":   "Bearer",
		})
	}

	return c.SendString(jwt + "&token_type=Bearer")

}
