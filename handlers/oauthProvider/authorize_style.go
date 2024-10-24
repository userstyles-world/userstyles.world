package oauthprovider

import (
	"fmt"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt"

	jwtware "userstyles.world/handlers/jwt"
	"userstyles.world/models"
	"userstyles.world/modules/config"
	"userstyles.world/modules/log"
	"userstyles.world/modules/storage"
	"userstyles.world/modules/util"
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
	jwtToken, err := util.NewJWT().
		SetClaim("state", state).
		SetClaim("userID", u.ID).
		SetExpiration(time.Now().Add(time.Hour * 2)).
		GetSignedString(util.OAuthPSigningKey)
	if err != nil {
		log.Warn.Println("Failed to create a JWT Token:", err.Error())
		return errorMessage(c, 500, "Couldn't make JWT Token, Error: Please notify the UserStyles.world admins.")
	}
	secureToken := util.EncryptText(jwtToken, util.AEADOAuthp, config.Secrets)

	styles, err := storage.FindStyleCardsForUsername(u.Username)
	if err != nil {
		log.Warn.Printf("Failed to get styles for %q: %s\n", u.Username, err)
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

	return c.Render("oauth/link", arguments)
}

func OAuthStylePost(c *fiber.Ctx) error {
	u, _ := jwtware.User(c)

	i := c.QueryInt("styleID")
	if i < 1 {
		return errorMessage(c, 400, "Invalid style ID")
	}

	oauth, err := models.GetOAuthByID(c.Query("oauthID"))
	if err != nil || oauth.ID == 0 {
		return errorMessage(c, 400, "Incorrect oauthID specified")
	}

	unsealedText, err := util.DecryptText(c.Query("token"), util.AEADOAuthp, config.Secrets)
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

	userID, ok := claims["userID"].(float64)
	if !ok || userID != float64(u.ID) {
		log.Warn.Println("Failed to get userID from parsed token.")
		return errorMessage(c, 500, "Error: Please notify the UserStyles.world admins.")
	}

	state, ok := claims["state"].(string)
	if !ok {
		log.Warn.Println("Invalid JWT state.")
		return errorMessage(c, 500, "Error: Please notify the UserStyles.world admins.")
	}

	style, err := models.GetStyleByID(i)
	if err != nil {
		log.Warn.Printf("Failed to find style %d: %s\n", i, err)
		return errorMessage(c, 500, "Couldn't retrieve style of user")
	}

	if style.UserID != u.ID {
		log.Warn.Println("Failed to match style author and userID.")
		return errorMessage(c, 500, "Error: Please notify the UserStyles.world admins.")
	}

	jwtToken, err := util.NewJWT().
		SetClaim("state", state).
		SetClaim("userID", u.ID).
		SetClaim("styleID", style.ID).
		SetExpiration(time.Now().Add(time.Minute * 10)).
		GetSignedString(util.OAuthPSigningKey)
	if err != nil {
		log.Warn.Println("Failed to create JWT Token:", err.Error())
		return errorMessage(c, 500, "Error: Please notify the UserStyles.world admins.")
	}

	returnCode := "?code=" + util.EncryptText(jwtToken, util.AEADOAuthp, config.Secrets)
	returnCode += "&style_id=" + strconv.Itoa(i)
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

	unsealedText, err := util.DecryptText(secureToken, util.AEADOAuthp, config.Secrets)
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

	userID, ok := claims["userID"].(float64)
	if !ok || userID != float64(u.ID) {
		log.Warn.Println("Failed to get userID from parsed token.")
		return errorMessage(c, 500, "Error: Please notify the UserStyles.world admins.")
	}

	_, ok = claims["state"].(string)
	if !ok {
		log.Warn.Println("Invalid JWT state.")
		return errorMessage(c, 500, "Error: Please notify the UserStyles.world admins.")
	}

	return c.Render("style/add", fiber.Map{
		"Title":       "Add userstyle",
		"User":        u,
		"Method":      "api",
		"OAuthID":     oauthID,
		"SecureToken": secureToken,
	})
}
