package oauthprovider

import (
	"fmt"
	"log"
	"time"

	"github.com/form3tech-oss/jwt-go"
	"github.com/gofiber/fiber/v2"

	"userstyles.world/database"
	jwtware "userstyles.world/handlers/jwt"
	"userstyles.world/models"
	"userstyles.world/utils"
)

func AuthorizeStyleGet(c *fiber.Ctx) error {
	u, _ := jwtware.User(c)

	// Under no circumstance this page should be loaded in some third-party frame.
	// It should be fully the user's consent to choose to authorize.
	// https://developer.mozilla.org/en-US/docs/Web/HTTP/Headers/X-Frame-Options
	c.Response().Header.Set("X-Frame-Options", "DENY")

	clientID, state := c.Query("client_id"), c.Query("state")
	if clientID == "" {
		return errorMessage(c, 400, "No client_id specified")
	}
	OAuth, err := models.GetOAuthByClientID(database.DB, clientID)
	if err != nil || OAuth.ID == 0 {
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

	styles, err := models.GetStylesByUser(database.DB, u.Username)
	if err != nil {
		log.Println("Error: Mo styles find for user", err.Error())
		return errorMessage(c, 500, "Couldn't retrieve styles of user")
	}

	arguments := fiber.Map{
		"User":        u,
		"Styles":      styles,
		"OAuth":       OAuth,
		"SecureToken": utils.PrepareText(jwt, utils.AEAD_OAUTHP),
	}
	for _, v := range OAuth.Scopes {
		arguments["Scope_"+v] = true
	}

	return c.Render("authorize_style", arguments)
}

func AuthorizeStylePost(c *fiber.Ctx) error {
	u, _ := jwtware.User(c)
	styleID, oauthID, secureToken := c.Query("styleID"), c.Query("oauthID"), c.Query("token")

	OAuth, err := models.GetOAuthByID(database.DB, oauthID)
	if err != nil || OAuth.ID == 0 {
		return errorMessage(c, 400, "Incorrect oauthID specified")
	}

	unsealedText, err := utils.DecodePreparedText(secureToken, utils.AEAD_OAUTHP)
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

	userID, ok := claims["userID"].(float64)
	if !ok || userID != float64(u.ID) {
		fmt.Println("WARNING!: Invalid userID")
		return errorMessage(c, 500, "JWT Token error, please notify the admins.")
	}

	state, ok := claims["state"].(string)
	if !ok {
		fmt.Println("WARNING!: Invalid state")
		return errorMessage(c, 500, "JWT Token error, please notify the admins.")
	}

	style, err := models.GetStyleByID(database.DB, styleID)
	if err != nil {
		fmt.Println("Error: Style wasn't found, due to: ", err.Error())
		return errorMessage(c, 500, "Couldn't retrieve style of user")
	}

	if style.UserID != u.ID {
		fmt.Println("WARNING!: Invalid style's user ID")
		return errorMessage(c, 500, "JWT Token error, please notify the admins.")
	}

	jwt, err := utils.NewJWTToken().
		SetClaim("state", state).
		SetClaim("userID", u.ID).
		SetClaim("styleID", style.ID).
		SetExpiration(time.Now().Add(time.Minute * 10)).
		GetSignedString(utils.OAuthPSigningKey)
	if err != nil {
		fmt.Println("Error: Couldn't create JWT Token:", err.Error())
		return errorMessage(c, 500, "JWT Token error, please notify the admins.")
	}

	returnCode := "?code=" + utils.PrepareText(jwt, utils.AEAD_OAUTHP)
	returnCode += "&style_id=" + styleID
	if state != "" {
		returnCode += "&state=" + state
	}

	return c.Redirect(OAuth.RedirectURI + "/" + returnCode)
}
