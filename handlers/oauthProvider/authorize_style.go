package oauthprovider

import (
	"fmt"
	"time"

	"github.com/form3tech-oss/jwt-go"
	"github.com/gofiber/fiber/v2"

	jwtware "userstyles.world/handlers/jwt"
	"userstyles.world/models"
	"userstyles.world/modules/config"
	"userstyles.world/modules/log"
	"userstyles.world/utils"
)

func OAuthStyleGet(c *fiber.Ctx) error {
	u, _ := jwtware.User(c)

	// Under no circumstance this page should be loaded in some third-party frame.
	// It should be fully the user's consent to choose to authorize.
	// https://developer.mozilla.org/en-US/docs/Web/HTTP/Headers/X-Frame-Options
	c.Response().Header.Set("X-Frame-Options", "DENY")

	clientID, state, vendorData := c.Query("client_id"), c.Query("state"), c.Query("vendor_data")
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
		log.Warn.Println("Failed to create a JWT Token:", err.Error())
		return errorMessage(c, 500, "Couldn't make JWT Token, please notify the admins.")
	}
	secureToken := utils.EncryptText(jwt, utils.AEADOAuthp, config.ScrambleConfig)

	styles, err := models.GetStylesByUser(u.Username)
	if err != nil {
		log.Info.Println("Failed to get styles from user:", err.Error())
		return errorMessage(c, 500, "Couldn't retrieve styles of user")
	}

	if len(styles) == 0 {
		return c.Redirect(fmt.Sprintf(
			"/api/oauth/style/new?token=%s&oauthID=%d&vendor_data=%s", secureToken, oauth.ID, vendorData),
			fiber.StatusSeeOther)
	}

	arguments := fiber.Map{
		"User":        u,
		"Styles":      styles,
		"OAuth":       oauth,
		"SecureToken": secureToken,
		"VendorData":  vendorData,
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
		log.Warn.Println("Failed to unseal JWT text:", err.Error())
		return errorMessage(c, 500, "JWT Token error, please notify the admins.")
	}

	token, err := jwt.Parse(unsealedText, utils.OAuthPJwtKeyFunction)
	if err != nil || !token.Valid {
		log.Warn.Println("Failed to unseal JWT token:", err.Error())
		return errorMessage(c, 500, "JWT Token error, please notify the admins.")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		log.Warn.Println("Failed to parse JWT Token:", err.Error())
		return errorMessage(c, 500, "JWT Token error, please notify the admins.")
	}

	userID, ok := claims["userID"].(float64)
	if !ok || userID != float64(u.ID) {
		log.Warn.Println("Failed to get userID from parsed token.")
		return errorMessage(c, 500, "JWT Token error, please notify the admins.")
	}

	state, ok := claims["state"].(string)
	if !ok {
		log.Warn.Println("Invalid JWT state.")
		return errorMessage(c, 500, "JWT Token error, please notify the admins.")
	}

	style, err := models.GetStyleByID(styleID)
	if err != nil {
		log.Warn.Printf("Failed to find style %v: %v\n", styleID, err)
		return errorMessage(c, 500, "Couldn't retrieve style of user")
	}

	if style.UserID != u.ID {
		log.Warn.Println("Failed to match style author and userID.")
		return errorMessage(c, 500, "JWT Token error, please notify the admins.")
	}

	jwt, err := utils.NewJWTToken().
		SetClaim("state", state).
		SetClaim("userID", u.ID).
		SetClaim("styleID", style.ID).
		SetExpiration(time.Now().Add(time.Minute * 10)).
		GetSignedString(utils.OAuthPSigningKey)
	if err != nil {
		log.Warn.Println("Failed to create JWT Token:", err.Error())
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
		log.Warn.Println("Failed to unseal JWT text:", err.Error())
		return errorMessage(c, 500, "JWT Token error, please notify the admins.")
	}

	token, err := jwt.Parse(unsealedText, utils.OAuthPJwtKeyFunction)
	if err != nil || !token.Valid {
		log.Warn.Println("Failed to unseal JWT token:", err.Error())
		return errorMessage(c, 500, "JWT Token error, please notify the admins.")
	}
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		log.Warn.Println("Failed to parse JWT Token:", err.Error())
		return errorMessage(c, 500, "JWT Token error, please notify the admins.")
	}

	userID, ok := claims["userID"].(float64)
	if !ok || userID != float64(u.ID) {
		log.Warn.Println("Failed to get userID from parsed token.")
		return errorMessage(c, 500, "JWT Token error, please notify the admins.")
	}

	_, ok = claims["state"].(string)
	if !ok {
		log.Warn.Println("Invalid JWT state.")
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
