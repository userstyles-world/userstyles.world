package api

import (
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"

	"userstyles.world/models"
	"userstyles.world/modules/config"
	"userstyles.world/modules/database"
	"userstyles.world/modules/errors"
	"userstyles.world/modules/log"
	"userstyles.world/modules/oauthlogin"
	"userstyles.world/utils"
)

var allowedErrosList = []error{
	errors.ErrPrimaryEmailNotVerified,
	errors.ErrNoServiceDetected,
}

func CallbackGet(c *fiber.Ctx) error {
	// Get the necessary information.
	redirectCode, tempCode, state := c.Params("rcode"), c.Query("code"), c.Query("state")
	if redirectCode == "" || tempCode == "" {
		log.Info.Println("No redirectCode or tempCode was detected.")
		// Give them the bad endpoint error.
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

	user, err := findOrMigrateUser(response)
	if err != nil {
		log.Warn.Printf("Failed to find or migrate %q: %s\n", response.Username, err)
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

	if err := user.UpdateLastLogin(); err != nil {
		log.Database.Printf("Failed to update last_login for %d\n", user.ID)
	}

	c.Cookie(&fiber.Cookie{
		Name:     fiber.HeaderAuthorization,
		Value:    t,
		Path:     "/",
		Expires:  expiration,
		Secure:   config.Production,
		HTTPOnly: true,
		SameSite: fiber.CookieSameSiteLaxMode,
	})

	return c.Redirect("/account", fiber.StatusSeeOther)
}

func findOrMigrateUser(res oauthlogin.OAuthResponse) (models.User, error) {
	var eu models.ExternalUser
	err := database.Conn.
		Preload("User").Model(eu).
		Where("provider = ?", string(res.Provider)).
		Where("external_id = ?", res.ExternalID).
		First(&eu).Error

	if err == gorm.ErrRecordNotFound {
		err = database.Conn.First(&eu.User, "username = ?", res.Username).Error
		if err != nil {
			return models.User{}, err
		}

		eu.ExternalID = strconv.Itoa(res.ExternalID)
		eu.Provider = string(res.Provider)
		eu.Email = res.Email
		eu.Username = res.Username
		eu.ExternalURL = setExternalURL(res.Provider, res.Username)
		eu.AccessToken = res.AccessToken
		eu.RawData = res.RawData

		if err = database.Conn.Create(&eu).Error; err != nil {
			return models.User{}, err
		}
	} else if err != nil {
		return models.User{}, err
	}

	return eu.User, nil
}

func setExternalURL(service oauthlogin.Service, username string) string {
	switch service {
	case oauthlogin.GithubService:
		return "https://github.com/" + username
	case oauthlogin.GitlabService:
		return "https://gitlab.com/" + username
	case oauthlogin.CodebergService:
		return "https://codeberg.org/" + username
	default:
		return ""
	}
}
