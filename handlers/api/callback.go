package api

import (
	"log"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"

	"userstyles.world/models"
	"userstyles.world/modules/config"
	"userstyles.world/modules/database"
	"userstyles.world/utils"
)

func getSocialMediaValue(user *models.User, social string) string {
	switch social {
	case "github":
		return user.Socials.Github
	case "gitlab":
		return user.Socials.Gitlab
	case "codeberg":
		return user.Socials.Codeberg
	default:
		return ""
	}
}

func CallbackGet(c *fiber.Ctx) error {
	// Get the necessary information.
	redirectCode, tempCode, state := c.Params("rcode"), c.Query("code"), c.Query("state")
	if redirectCode == "" || tempCode == "" {
		log.Println("No redirectcode or tempCode was detected")
		// Give them the bad enpoint error.
		return c.Next()
	}
	var service string
	var rState string
	if redirectCode != "codeberg" && redirectCode != "gitlab" {
		service = "github"
		// Decode the string so we get our actual information back.
		code, err := utils.DecodePreparedText(redirectCode, utils.AEAD_OAUTH)
		if err != nil {
			log.Println("Error: Couldn't decode our prepared text.")
			return c.Next()
		}
		rState = code

		if rState != state {
			log.Println("Error: The state doesn't match!")
			return c.Next()
		}
	} else {
		service = redirectCode
	}

	response, err := utils.CallbackOAuth(tempCode, rState, service)
	if err != nil {
		log.Println("Ouch, the response failed, due to: " + err.Error())
		return c.Next()
	}

	user, err := models.FindUserByNameOrEmail(response.Username, response.Email)
	if err != nil {
		if err.Error() != "User not found." && err.Error() != "record not found" {
			return c.Next()
		}
		user = &models.User{
			Username:      response.Username,
			Email:         response.Email,
			Role:          models.Regular,
			OAuthProvider: service,
		}
		regErr := database.Conn.Create(user)

		if regErr.Error != nil {
			log.Printf("Failed to register %s, error: %s", response.Username, regErr.Error)
			return c.Status(fiber.StatusInternalServerError).
				JSON(fiber.Map{
					"data": "Internal Error.",
				})
		}
	}

	// TODO: Simplify this logic.
	if (user.OAuthProvider == "none" || user.OAuthProvider != service) &&
		!strings.EqualFold(getSocialMediaValue(user, service), response.Username) {
		log.Println("User detected but the social media value wasn't set of this user.")
		return c.Next()
	}

	expiration := time.Now().Add(time.Hour * 24 * 14)
	t, err := utils.NewJWTToken().
		SetClaim("id", user.ID).
		SetClaim("name", user.Username).
		SetClaim("role", user.Role).
		SetExpiration(expiration).
		GetSignedString(nil)
	if err != nil {
		log.Println("Couldn't create JWT Token, due to " + err.Error())
		return c.Status(fiber.StatusInternalServerError).
			JSON(fiber.Map{
				"data": "Internal Error.",
			})
	}

	c.Cookie(&fiber.Cookie{
		Name:     fiber.HeaderAuthorization,
		Value:    t,
		Path:     "/",
		Expires:  expiration,
		Secure:   config.Production,
		HTTPOnly: true,
		SameSite: "lax",
	})

	return c.Redirect("/account", fiber.StatusSeeOther)
}
