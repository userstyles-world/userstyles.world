package oauthprovider

import (
	"strconv"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"

	jwtware "userstyles.world/handlers/jwt"
	"userstyles.world/models"
	"userstyles.world/modules/config"
	"userstyles.world/modules/log"
	"userstyles.world/modules/util"
)

func errorMessage(c *fiber.Ctx, status int, errorMessage string) error {
	return c.Status(status).
		JSON(fiber.Map{
			"data": errorMessage,
		})
}

func redirectFunction(c *fiber.Ctx, state, redirectURI string) error {
	u, _ := jwtware.User(c)

	jwtToken, err := util.NewJWT().
		SetClaim("state", state).
		SetClaim("userID", u.ID).
		SetExpiration(time.Now().Add(time.Minute * 10)).
		GetSignedString(util.OAuthPSigningKey)
	if err != nil {
		log.Warn.Println("Failed to create a JWT Token:", err.Error())
		return errorMessage(c, 500, "Error: Please notify the UserStyles.world admins.")
	}

	returnCode := "?code=" + util.EncryptText(jwtToken, util.AEADOAuthp, config.Secrets)
	if state != "" {
		returnCode += "&state=" + state
	}

	return c.Redirect(redirectURI + "/" + returnCode)
}

func AuthorizeGet(c *fiber.Ctx) error {
	u, _ := jwtware.User(c)

	// Under no circumstance this page should be loaded in some third-party frame.
	// It should be fully the user's consent to choose to authorize.
	// https://developer.mozilla.org/en-US/docs/Web/HTTP/Headers/X-Frame-Options
	c.Response().Header.Set("X-Frame-Options", "DENY")

	clientID, state, scope := c.Query("client_id"), c.Query("state"), c.Query("scope")
	if clientID == "" {
		return errorMessage(c, 400, "No client_id specified")
	}
	oauth, err := models.GetOAuthByClientID(clientID)
	if err != nil || oauth.ID == 0 {
		return errorMessage(c, 400, "Incorrect client_id specified")
	}

	user, err := models.FindUserByName(u.Username)
	if err != nil {
		return errorMessage(c, 500, "Notify the admins.")
	}

	// Check if the user has already authorized this OAuth application.
	if util.ContainsString(user.AuthorizedOAuth, strconv.Itoa(int(oauth.ID))) {
		return redirectFunction(c, state, oauth.RedirectURI)
	}

	// Convert it to actual []string
	scopes := strings.Split(scope, " ")

	// Just check if the application has actually set if they will request these scopes.
	if !util.EveryString(scopes, func(name string) bool {
		return util.ContainsString(oauth.Scopes, name)
	}) {
		return errorMessage(c, 400, "An scope was provided which isn't selected in the OAuth's settings selection.")
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

	arguments := fiber.Map{
		"User":        u,
		"OAuth":       oauth,
		"SecureToken": util.EncryptText(jwtToken, util.AEADOAuthp, config.Secrets),
	}
	for _, v := range oauth.Scopes {
		arguments["Scope_"+v] = true
	}

	return c.Render("oauth/authorize", arguments)
}

func AuthPost(c *fiber.Ctx) error {
	u, _ := jwtware.User(c)
	oauthID, secureToken := c.Params("id"), c.Params("token")

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
		log.Warn.Println("Failed to parse JWT token:", err.Error())
		return errorMessage(c, 500, "Error: Please notify the UserStyles.world admins.")
	}

	userID, ok := claims["userID"].(float64)
	if !ok || userID != float64(u.ID) {
		log.Warn.Println("Failed to get userID from parsed token")
		return errorMessage(c, 500, "Error: Please notify the UserStyles.world admins.")
	}

	user, err := models.FindUserByName(u.Username)
	if err != nil {
		log.Warn.Printf("Failed to find user %v: %v\n", userID, err)
		return errorMessage(c, 500, "Error: Please notify the UserStyles.world admins.")
	}

	user.AuthorizedOAuth = append(user.AuthorizedOAuth, oauthID)
	if err = models.UpdateUser(user); err != nil {
		log.Warn.Printf("Failed to update user %d: %v\n", user.ID, err)
		return errorMessage(c, 500, "Error: Please notify the UserStyles.world admins.")
	}

	return redirectFunction(c, claims["state"].(string), oauth.RedirectURI)
}
