package oauthprovider

import (
	"errors"
	"strconv"
	"strings"

	val "github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"

	"userstyles.world/handlers/jwt"
	"userstyles.world/models"
	"userstyles.world/modules/log"
	"userstyles.world/modules/util"
	"userstyles.world/modules/validator"
	"userstyles.world/utils"
)

type OAuthSettingMethod uint8

const (
	methodAdd OAuthSettingMethod = iota
	methodEdit
)

func OAuthSettingsGet(c *fiber.Ctx) error {
	u, _ := jwt.User(c)

	var method OAuthSettingMethod
	var oauth *models.APIOAuth
	var err error
	if id := c.Params("id"); id != "" {
		method = methodEdit
		oauth, err = models.GetOAuthByID(id)
	} else {
		method = methodAdd
	}

	if err != nil {
		log.Warn.Printf("Failed to OAuth %d: %s\n", oauth.ID, err.Error())
		return c.Render("err", fiber.Map{
			"Title": "Internal server error",
			"User":  u,
		})
	}

	if method == methodEdit {
		if u.ID != oauth.UserID {
			return c.Render("err", fiber.Map{
				"Title": "Users don't match",
				"User":  u,
			})
		}
	}
	oauths, err := models.ListOAuthsOfUser(u.Username)
	if err != nil {
		if method == methodEdit {
			arguments := fiber.Map{
				"Title":  "OAuth Settings",
				"User":   u,
				"OAuth":  oauth,
				"OAuths": nil,
				"Method": method,
			}
			for _, v := range oauth.Scopes {
				arguments["Scope_"+v] = true
			}
			return c.Render("oauth/settings", arguments)
		}
		return c.Render("oauth/settings", fiber.Map{
			"Title":  "OAuth Settings",
			"User":   u,
			"OAuths": nil,
			"Method": method,
		})
	}

	if method == methodEdit {
		arguments := fiber.Map{
			"Title":  "OAuth Settings",
			"User":   u,
			"OAuth":  oauth,
			"OAuths": oauths,
			"Method": method,
		}
		for _, v := range oauth.Scopes {
			arguments["Scope_"+v] = true
		}
		return c.Render("oauth/settings", arguments)
	}
	return c.Render("oauth/settings", fiber.Map{
		"Title":  "OAuth Settings",
		"User":   u,
		"OAuths": oauths,
		"Method": method,
	})
}

func OAuthSettingsPost(c *fiber.Ctx) error {
	u, _ := jwt.User(c)
	id := c.Params("id")

	q := models.OAuth{
		Name:        c.FormValue("name"),
		Description: c.FormValue("description"),
		RedirectURI: strings.TrimSuffix(c.FormValue("redirect_uri"), "/"),
		Scopes: utils.Filter([]string{"style", "user"}, func(name any) bool {
			return c.FormValue(name.(string)) == "on"
		}).([]string),
		UserID: u.ID,
	}

	if err := validator.V.StructPartial(q, "Name", "Description"); err != nil {
		var validationError val.ValidationErrors
		if ok := errors.As(err, &validationError); ok {
			log.Info.Println("Validation errors:", validationError)
		}

		arguments := fiber.Map{
			"Title":  "OAuth Settings",
			"Error":  "Failed to validate inputs.",
			"User":   u,
			"OAuth":  q,
			"Method": "add",
		}
		for _, v := range q.Scopes {
			arguments["Scope_"+v] = true
		}
		return c.Status(500).
			Render("oauth/settings", arguments)
	}

	var err error
	var dbOAuth *models.OAuth
	if id != "" {
		err = models.UpdateOAuth(&q, id)
	} else {
		q.ClientID = util.RandomString(32)
		q.ClientSecret = util.RandomString(128)
		dbOAuth, err = models.CreateOAuth(&q)
	}

	if err != nil {
		log.Warn.Printf("Updating OAuth settings for %v failed: %s\n", id, err.Error())
		return c.Render("err", fiber.Map{
			"Title": "Internal server error.",
			"User":  u,
		})
	}

	oauthID := strconv.FormatUint(uint64(dbOAuth.ID), 10)
	return c.Redirect("/oauth/settings/"+oauthID, fiber.StatusSeeOther)
}
