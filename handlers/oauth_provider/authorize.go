package oauth_provider

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"
	"time"

	"github.com/form3tech-oss/jwt-go"
	"github.com/gofiber/fiber/v2"
	"userstyles.world/database"
	jwtware "userstyles.world/handlers/jwt"
	"userstyles.world/models"
	"userstyles.world/utils"
)

// Checks if every entry of slice fulfills condition.
func every(arr interface{}, cond func(interface{}) bool) bool {
	contentValue := reflect.ValueOf(arr)

	for i := 0; i < contentValue.Len(); i++ {
		if content := contentValue.Index(i); !cond(content.Interface()) {
			return false
		}
	}
	return true
}

// Check if array contains certain entry.
func contains(arr []string, entry string) bool {
	for _, possibleEntry := range arr {
		if possibleEntry == entry {
			return true
		}
	}
	return false
}

func redirectFunction(c *fiber.Ctx, state, redirect_uri string) error {
	u, _ := jwtware.User(c)

	jwt, err := utils.NewJWTToken().
		SetClaim("state", state).
		SetClaim("username", u.Username).
		SetExpiration(time.Now().Add(time.Minute * 10)).
		GetSignedString(utils.OAuthPSigningKey)

	if err != nil {
		fmt.Println("Error: Couldn't create JWT Token:", err.Error())
		return c.Status(500).
			JSON(fiber.Map{
				"error": "JWT Token error, please notify the admins.",
			})
	}

	returnCode := "?code=" + utils.PrepareText(jwt, utils.AEAD_OAUTHP)
	if state != "" {
		returnCode += "&state=" + state
	}

	return c.Redirect(redirect_uri + "/" + returnCode)
}

func AuthorizeGet(c *fiber.Ctx) error {
	u, ok := jwtware.User(c)
	if !ok {
		// TODO: Make this template.
		return c.Status(401).
			Render("ask_login", fiber.Map{})
	}

	// Under no circumstance this page should be loaded in some third-party frame.
	// It should be fully the user's consent to choose to authorize.
	// https://developer.mozilla.org/en-US/docs/Web/HTTP/Headers/X-Frame-Options
	c.Response().Header.Set("X-Frame-Options", "DENY")

	clientID, state, scope := c.Query("client_id"), c.Query("state"), c.Query("scope")
	if clientID == "" {
		return c.Status(400).
			JSON(fiber.Map{
				"error": "No client_id specified",
			})
	}
	OAuth, err := models.GetOAuthByClientID(database.DB, clientID)
	if err != nil || OAuth.ID == 0 {
		return c.Status(400).
			JSON(fiber.Map{
				"error": "Incorrect client_id specified",
			})
	}

	user, err := models.FindUserByName(database.DB, u.Username)
	if err != nil {
		return c.Status(500).
			JSON(fiber.Map{
				"error": "Notify the admins.",
			})
	}

	// Check if the user has already authorized this OAuth application.
	if contains(user.AuthorizedOAuth, strconv.Itoa(int(OAuth.ID))) {
		return redirectFunction(c, state, OAuth.RedirectURI)
	}

	// Convert it to actual []string
	scopes := strings.Split(scope, " ")

	// Just check if the application has actually set if they will request these scopes.
	if !every(scopes, func(name interface{}) bool {
		return contains(OAuth.Scopes, name.(string))
	}) {
		return c.Status(400).
			JSON(fiber.Map{
				"error": "An scope was provided which isn't selected in the OAuth's settings selection.",
			})
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
		fmt.Println("Error: Couldn't make a JWT Token due to:", err.Error())
		return c.Status(500).
			JSON(fiber.Map{
				"error": "Couldn't make JWT Token, please notify the admins.",
			})
	}

	return c.Render("authorize", fiber.Map{
		"User":        u,
		"OAuth":       OAuth,
		"SecureToken": utils.PrepareText(jwt, utils.AEAD_OAUTHP),
	})
}

func AuthorizePost(c *fiber.Ctx) error {
	u, _ := jwtware.User(c)

	oauthID, secureToken := c.Params("id"), c.Params("token")

	OAuth, err := models.GetOAuthByID(database.DB, oauthID)
	if err != nil || OAuth.ID == 0 {
		return c.Status(400).
			JSON(fiber.Map{
				"error": "Incorrect oauthID specified",
			})
	}

	unsealedText, err := utils.DecodePreparedText(secureToken, utils.AEAD_OAUTHP)
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

	if _, ok := claims["userID"].(float64); !ok {
		fmt.Println("WARNING!: Invalid userID")
		return c.Status(500).
			JSON(fiber.Map{
				"error": "JWT Token error, please notify the admins.",
			})
	}

	if uint(claims["userID"].(float64)) != u.ID {
		fmt.Println("WARNING!: User got valid encrypted state, but userID differ!!!")
		return c.Status(500).
			JSON(fiber.Map{
				"error": "JWT Token error, please notify the admins.",
			})
	}

	user, err := models.FindUserByName(database.DB, u.Username)
	if err != nil {
		fmt.Println("Error: Couldn't retrieve user:", err.Error())
		return c.Status(500).
			JSON(fiber.Map{
				"error": "JWT Token error, please notify the admins.",
			})
	}

	user.AuthorizedOAuth = append(user.AuthorizedOAuth, oauthID)
	if err = models.UpdateUser(database.DB, user); err != nil {
		fmt.Println("Error: couldn't update user:", err.Error())
		return c.Status(500).
			JSON(fiber.Map{
				"error": "JWT Token error, please notify the admins.",
			})
	}

	return redirectFunction(c, claims["state"].(string), OAuth.RedirectURI)
}
