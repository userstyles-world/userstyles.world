package oauthprovider

import (
	"fmt"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"

	"userstyles.world/models"
	"userstyles.world/modules/config"
	"userstyles.world/modules/log"
	"userstyles.world/modules/util"
)

func TokenPost(c *fiber.Ctx) error {
	clientID, clientSecret, stateQuery, tCode :=
		c.FormValue("client_id"), c.FormValue("client_secret"), c.FormValue("state"), c.FormValue("code")

	if clientID == "" {
		return errorMessage(c, 400, "No client_id specified")
	}
	if clientSecret == "" {
		return errorMessage(c, 400, "No client_secret specified")
	}
	if tCode == "" {
		return errorMessage(c, 400, "No code specified")
	}

	oauth, err := models.GetOAuthByClientID(clientID)
	if err != nil || oauth.ID == 0 {
		return errorMessage(c, 400, "Incorrect client_id specified")
	}
	if oauth.ClientSecret != clientSecret {
		return errorMessage(c, 400, "Incorrect client_secret specified")
	}

	unsealedText, err := util.DecryptText(tCode, util.AEADOAuthp, config.Secrets)
	if err != nil {
		log.Warn.Println("Failed to unseal JWT text:", err.Error())
		return errorMessage(c, 500, "Error: Please notify the UserStyles.world admins.")
	}

	token, err := jwt.Parse(unsealedText, util.OAuthPJwtKeyFunction)
	if err != nil || !token.Valid {
		log.Warn.Println("Failed to unseal JWT token:", err.Error())
		return errorMessage(c, 500, "Error: Please notify the UserStyles.world admins.")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		log.Warn.Println("Failed to parse JWT Token:", err.Error())
		return errorMessage(c, 500, "Error: Please notify the UserStyles.world admins.")
	}

	state, ok := claims["state"].(string)
	if !ok {
		log.Warn.Println("Invalid JWT state.")
		return errorMessage(c, 500, "Error: Please notify the UserStyles.world admins.")
	}

	floatUserID, ok := claims["userID"].(float64)
	if !ok {
		log.Warn.Println("Failed to get userID from parsed token.")
		return errorMessage(c, 500, "Error: Please notify the UserStyles.world admins.")
	}
	userID := uint(floatUserID)

	fStyleID, ok := claims["styleID"].(float64)
	if !ok {
		fStyleID = 0
	}

	if stateQuery != state {
		return errorMessage(c, 500, "State doesn't match.")
	}

	user, err := models.FindUserByID(fmt.Sprintf("%d", userID))
	if err != nil || user.ID == 0 {
		return errorMessage(c, 500, "Couldn't find the user that was specified, Error: Please notify the UserStyles.world admins.")
	}

	var jwtToken string

	if styleID := uint(fStyleID); styleID != 0 {
		jwtToken, err = util.NewJWT().
			SetClaim("styleID", styleID).
			SetClaim("userID", user.ID).
			GetSignedString(util.OAuthPSigningKey)
	} else {
		jwtToken, err = util.NewJWT().
			SetClaim("scopes", strings.Join(oauth.Scopes, ",")).
			SetClaim("userID", user.ID).
			GetSignedString(util.OAuthPSigningKey)
	}

	if err != nil {
		return errorMessage(c, 500, "Couldn't create access_token Error: Please notify the UserStyles.world admins.")
	}

	if c.Accepts("application/json", "text/plain ") == "application/json" {
		return c.JSON(fiber.Map{
			"access_token": jwtToken,
			"token_type":   "Bearer",
		})
	}

	return c.SendString(jwtToken + "&token_type=Bearer")
}
