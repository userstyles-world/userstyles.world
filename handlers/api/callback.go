package api

import (
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"

	"userstyles.world/models"
	"userstyles.world/modules/config"
	"userstyles.world/modules/database"
	"userstyles.world/modules/errors"
	"userstyles.world/modules/log"
	"userstyles.world/modules/oauthlogin"
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

var allowedErrosList []error = []error{
	errors.ErrPrimaryEmailNotVerified,
	errors.ErrNoServiceDetected,
}

func CallbackGet(c *fiber.Ctx) error {
	// Get the necessary information.
	redirectCode, tempCode, state := c.Params("rcode"), c.Query("code"), c.Query("state")
	if redirectCode == "" || tempCode == "" {
		log.Info.Println("No redirectCode or tempCode was detected.")
		// Give them the bad enpoint error.
		return c.Next()
	}
	var service string
	var rState string
	if redirectCode != "codeberg" && redirectCode != "gitlab" {
		service = "github"
		// Decode the string so we get our actual information back.
		code, err := utils.DecryptText(redirectCode, utils.AEADOAuth, config.ScrambleConfig)
		if err != nil {
			log.Warn.Println("Failed to decode prepared text.")
			return c.Next()
		}
		rState = code

		if rState != state {
			log.Warn.Println("Failed to match states.")
			return c.Next()
		}
	} else {
		service = redirectCode
	}

	response, err := oauthlogin.CallbackOAuth(tempCode, rState, service)
	if err != nil {
		log.Warn.Println("Ouch, the response failed:", err.Error())
		// We only allow a certain amount of errors to be displayed to the
		// user. So we will now check if the error is in the "allowed" list
		// and if it is, we will display it to the user.
		if utils.ContainsError(allowedErrosList, err) {
			return c.Render("err", fiber.Map{
				"Title": err.Error(),
			})
		}
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
			log.Warn.Printf("Failed to register %s: %s", response.Username, regErr.Error)
			return c.Status(fiber.StatusInternalServerError).
				JSON(fiber.Map{
					"data": "Internal Error.",
				})
		}
	}

	// TODO: Simplify this logic.
	if (user.OAuthProvider == "none" || user.OAuthProvider != service) &&
		!strings.EqualFold(getSocialMediaValue(user, service), response.Username) {
		log.Warn.Println("User detected but the social media value wasn't set of this user.")
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
		log.Warn.Println("Failed to create JWT Token:", err.Error())
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
