package oauth_provider

import (
	"fmt"
	"reflect"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"userstyles.world/database"
	"userstyles.world/handlers/jwt"
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

func AuthorizeGet(c *fiber.Ctx) error {
	u, _ := jwt.User(c)
	// TODO: Chekc if user is not logged in and ask if they want to login/register.

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

	// Convert it to actual []string
	scopes := strings.Split(scope, " ")

	// Just check if the application has actually set if they will request these scopes.
	if !every(scopes, func(name interface{}) bool {
		return contains(OAuth.Scopes, name.(string))
	}) {
		return c.Status(400).
			JSON(fiber.Map{
				"error": "An scope was provided which isn't selected in your OAuth selection.",
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
