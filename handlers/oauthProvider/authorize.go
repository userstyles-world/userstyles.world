package oauthprovider

import (
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/form3tech-oss/jwt-go"
	"github.com/gofiber/fiber/v2"

	jwtware "userstyles.world/handlers/jwt"
	"userstyles.world/models"
	"userstyles.world/modules/config"
	"userstyles.world/utils"
)

func errorMessage(c *fiber.Ctx, status int, errorMessage string) error {
	return c.Status(status).
		JSON(fiber.Map{
			"data": errorMessage,
		})
}

func redirectFunction(c *fiber.Ctx, state, redirectURI string) error {
	u, _ := jwtware.User(c)

	jwt, err := utils.NewJWTToken().
		SetClaim("state", state).
		SetClaim("userID", u.ID).
		SetExpiration(time.Now().Add(time.Minute * 10)).
		GetSignedString(utils.OAuthPSigningKey)
	if err != nil {
		log.Println("Error: Couldn't create JWT Token:", err.Error())
		return errorMessage(c, 500, "JWT Token error, please notify the admins.")
	}

	returnCode := "?code=" + utils.EncryptText(jwt, utils.AEADOAuthp, config.ScrambleConfig)
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
	OAuth, err := models.GetOAuthByClientID(clientID)
	if err != nil || OAuth.ID == 0 {
		return errorMessage(c, 400, "Incorrect client_id specified")
	}

	user, err := models.FindUserByName(u.Username)
	if err != nil {
		return errorMessage(c, 500, "Notify the admins.")
	}

	// Check if the user has already authorized this OAuth application.
	if utils.Contains(user.AuthorizedOAuth, strconv.Itoa(int(OAuth.ID))) {
		return redirectFunction(c, state, OAuth.RedirectURI)
	}

	// Convert it to actual []string
	scopes := strings.Split(scope, " ")

	// Just check if the application has actually set if they will request these scopes.
	if !utils.Every(scopes, func(name interface{}) bool {
		return utils.Contains(OAuth.Scopes, name.(string))
	}) {
		return errorMessage(c, 400, "An scope was provided which isn't selected in the OAuth's settings selection.")
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

	arguments := fiber.Map{
		"User":        u,
		"OAuth":       OAuth,
		"SecureToken": utils.EncryptText(jwt, utils.AEADOAuthp, config.ScrambleConfig),
	}
	for _, v := range OAuth.Scopes {
		arguments["Scope_"+v] = true
	}

	return c.Render("authorize", arguments)
}

func AuthPost(c *fiber.Ctx) error {
	u, _ := jwtware.User(c)
	oauthID, secureToken := c.Params("id"), c.Params("token")

	OAuth, err := models.GetOAuthByID(oauthID)
	if err != nil || OAuth.ID == 0 {
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
		log.Println("Error: Couldn't parse JWT Token:", err.Error())
		return errorMessage(c, 500, "JWT Token error, please notify the admins.")
	}

	userID, ok := claims["userID"].(float64)
	if !ok || userID != float64(u.ID) {
		log.Println("WARNING!: Invalid userID")
		return errorMessage(c, 500, "JWT Token error, please notify the admins.")
	}

	user, err := models.FindUserByName(u.Username)
	if err != nil {
		log.Println("Error: Couldn't retrieve user:", err.Error())
		return errorMessage(c, 500, "JWT Token error, please notify the admins.")
	}

	user.AuthorizedOAuth = append(user.AuthorizedOAuth, oauthID)
	if err = models.UpdateUser(user); err != nil {
		log.Println("Error: couldn't update user:", err.Error())
		return errorMessage(c, 500, "JWT Token error, please notify the admins.")
	}

	return redirectFunction(c, claims["state"].(string), OAuth.RedirectURI)
}
