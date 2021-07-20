package oauthprovider

import (
	"fmt"
	"log"
	"time"

	"github.com/form3tech-oss/jwt-go"
	"github.com/gofiber/fiber/v2"

	jwtware "userstyles.world/handlers/jwt"
	"userstyles.world/models"
	"userstyles.world/modules/config"
	"userstyles.world/utils"
)

func OAuthStyleGet(c *fiber.Ctx) error {
	u, _ := jwtware.User(c)

	// Under no circumstance this page should be loaded in some third-party frame.
	// It should be fully the user's consent to choose to authorize.
	// https://developer.mozilla.org/en-US/docs/Web/HTTP/Headers/X-Frame-Options
	c.Response().Header.Set("X-Frame-Options", "DENY")

	clientID, state := c.Query("client_id"), c.Query("state")
	if clientID == "" {
		return errorMessage(c, 400, "No client_id specified")
	}
	oauth, err := models.GetOAuthByClientID(clientID)
	if err != nil || oauth.ID == 0 {
		return errorMessage(c, 400, "Incorrect client_id specified")
	}

	// User has to authorize within 2 hours.
	// To migate any weird attack we include the ID of the user that wishes to authorize.
	// Such that this key cannot be replaced by some other user.
	// And to follow our weird state-less design we include the state.
	// Thus not storing the state.
	jwt, err := utils.NewJWTToken().
		SetClaim("state", state).
		SetClaim("userID", u.ID).
		SetExpiration(time.Now().Add(time.Hour * 2)).
		GetSignedString(utils.OAuthPSigningKey)
	if err != nil {
		log.Println("Error: Couldn't make a JWT Token due to:", err.Error())
		return errorMessage(c, 500, "Couldn't make JWT Token, please notify the admins.")
	}
	secureToken := utils.EncryptText(jwt, utils.AEADOAuthp, config.ScrambleConfig)

	styles, err := models.GetStylesByUser(u.Username)
	if err != nil {
		log.Println("Error: Mo styles find for user", err.Error())
		return errorMessage(c, 500, "Couldn't retrieve styles of user")
	}

	if len(styles) == 0 {
		return c.Redirect(fmt.Sprintf(
			"/api/oauth/style/new?token=%s&oauthID=%d", secureToken, oauth.ID),
			fiber.StatusSeeOther)
	}

	arguments := fiber.Map{
		"User":        u,
		"Styles":      styles,
		"OAuth":       oauth,
		"SecureToken": secureToken,
	}
	for _, v := range oauth.Scopes {
		arguments["Scope_"+v] = true
	}

	return c.Render("authorize_style", arguments)
}

func OAuthStylePost(c *fiber.Ctx) error {
	u, _ := jwtware.User(c)
	styleID, oauthID, secureToken := c.Query("styleID"), c.Query("oauthID"), c.Query("token")

	oauth, err := models.GetOAuthByID(oauthID)
	if err != nil || oauth.ID == 0 {
		return errorMessage(c, 400, "Incorrect oauthID specified")
	}

	unsealedText, err := utils.DecryptText(secureToken, utils.AEADOAuthp, config.ScrambleConfig)
	if err != nil {
		log.Println("Error: Couldn't unseal JWT Token:", err.Error())
		return errorMessage(c, 500, "JWT Token error, please notify the admins.")
	}

	token, err := jwt.Parse(unsealedText, utils.OAuthPJwtKeyFunction)
	if err != nil || !token.Valid {
		log.Println("Error: Couldn't unseal JWT Token:", err.Error())
		return errorMessage(c, 500, "JWT Token error, please notify the admins.")
	}
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		log.Println("Error: Couldn't get claims from JWT Token:", err.Error())
		return errorMessage(c, 500, "JWT Token error, please notify the admins.")
	}

	userID, ok := claims["userID"].(float64)
	if !ok || userID != float64(u.ID) {
		log.Println("WARNING!: Invalid userID")
		return errorMessage(c, 500, "JWT Token error, please notify the admins.")
	}

	state, ok := claims["state"].(string)
	if !ok {
		log.Println("WARNING!: Invalid state")
		return errorMessage(c, 500, "JWT Token error, please notify the admins.")
	}

	style, err := models.GetStyleByID(styleID)
	if err != nil {
		log.Println("Error: Style wasn't found, due to: ", err.Error())
		return errorMessage(c, 500, "Couldn't retrieve style of user")
	}

	if style.UserID != u.ID {
		log.Println("WARNING!: Invalid style's user ID")
		return errorMessage(c, 500, "JWT Token error, please notify the admins.")
	}

	jwt, err := utils.NewJWTToken().
		SetClaim("state", state).
		SetClaim("userID", u.ID).
		SetClaim("styleID", style.ID).
		SetExpiration(time.Now().Add(time.Minute * 10)).
		GetSignedString(utils.OAuthPSigningKey)
	if err != nil {
		log.Println("Error: Couldn't create JWT Token:", err.Error())
		return errorMessage(c, 500, "JWT Token error, please notify the admins.")
	}

	returnCode := "?code=" + utils.EncryptText(jwt, utils.AEADOAuthp, config.ScrambleConfig)
	returnCode += "&style_id=" + styleID
	if state != "" {
		returnCode += "&state=" + state
	}

	return c.Redirect(oauth.RedirectURI + "/" + returnCode)
}

func OAuthStyleNewPost(c *fiber.Ctx) error {
	u, _ := jwtware.User(c)
	oauthID, secureToken := c.Query("oauthID"), c.Query("token")

	oauth, err := models.GetOAuthByID(oauthID)
	if err != nil || oauth.ID == 0 {
		return errorMessage(c, 400, "Incorrect oauthID specified")
	}

	unsealedText, err := utils.DecryptText(secureToken, utils.AEADOAuthp, config.ScrambleConfig)
	if err != nil {
		log.Println("Error: Couldn't unseal JWT Token:", err.Error())
		return errorMessage(c, 500, "JWT Token error, please notify the admins.")
	}

	token, err := jwt.Parse(unsealedText, utils.OAuthPJwtKeyFunction)
	if err != nil || !token.Valid {
		log.Println("Error: Couldn't unseal JWT Token:", err.Error())
		return errorMessage(c, 500, "JWT Token error, please notify the admins.")
	}
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		log.Println("WARNING!: Invalid conversion")
		return errorMessage(c, 500, "JWT Token error, please notify the admins.")
	}

	userID, ok := claims["userID"].(float64)
	if !ok || userID != float64(u.ID) {
		log.Println("WARNING!: Invalid userID")
		return errorMessage(c, 500, "JWT Token error, please notify the admins.")
	}

	_, ok = claims["state"].(string)
	if !ok {
		log.Println("WARNING!: Invalid state")
		return errorMessage(c, 500, "JWT Token error, please notify the admins.")
	}

	return c.Render("style/create", fiber.Map{
		"Title":       "Add userstyle",
		"User":        u,
		"Method":      "add_api",
		"OAuthID":     oauthID,
		"SecureToken": secureToken,
	})
}
