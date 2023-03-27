package api

import (
	"fmt"
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

	user, err := flow(response)
	if err != nil {
		log.Warn.Printf("User %q failed to sign in: %s\n", response.Username, err)
		msg := "Please contact us and provide this timestamp: " + time.Now().Format(time.RFC3339)
		return c.Render("err", fiber.Map{"Title": msg})
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

func flow(o oauthlogin.OAuthResponse) (*models.User, error) {
	var eu models.ExternalUser

	// Check if external user exists.
	err := database.Conn.Debug().
		Model(eu).Preload("User").
		First(&eu, "provider = ? AND external_id = ?", o.Provider, o.ExternalID).Error

	switch {
	case err == gorm.ErrRecordNotFound:
		// Check if user exists.
		err := database.Conn.Debug().
			First(&eu.User, "username = ? AND o_auth_provider = ?", o.Username, o.Provider).Error
		if err != nil {
			eu = models.ExternalUser{
				ExternalID:  o.ExternalID,
				Provider:    string(o.Provider),
				Email:       o.Email,
				Username:    o.Username,
				ExternalURL: o.ProfileURL(),
				AccessToken: o.AccessToken,
				RawData:     o.RawData,
				User: models.User{
					Email:         o.Email,
					Username:      o.Username,
					OAuthProvider: string(o.Provider),
				},
			}

			var count int64
			err := database.Conn.Debug().
				Model(eu.User).
				Where("username = ? OR email = ?", eu.Username, eu.Email).
				Count(&count).
				Error
			if err != nil {
				return nil, err
			}
			if count > 0 {
				return nil, fmt.Errorf("user %s already exists", eu.Username)
			}

			if err := database.Conn.Debug().Create(&eu).Error; err != nil {
				return nil, err
			}

			log.Info.Printf("Migrated user id=%d username=%s", eu.User.ID, eu.User.Username)

			return &eu.User, nil
		} else {
			// NOTE: This overrides existing values.
			eu = models.ExternalUser{
				User:        eu.User,
				UserID:      eu.User.ID,
				ExternalID:  o.ExternalID,
				Provider:    string(o.Provider),
				Email:       o.Email,
				Username:    o.Username,
				ExternalURL: o.ProfileURL(),
				AccessToken: o.AccessToken,
				RawData:     o.RawData,
			}
			if err := database.Conn.Debug().Create(&eu).Error; err != nil {
				return nil, err
			}

			log.Info.Printf("Created user id=%d username=%s", eu.User.ID, eu.User.Username)

			return &eu.User, nil
		}

	case err != nil:
		return nil, err
	}

	log.Info.Printf("Logged in user id=%d username=%s", eu.User.ID, eu.User.Username)

	return &eu.User, nil
}
