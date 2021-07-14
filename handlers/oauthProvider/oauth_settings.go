package oauthprovider

import (
	"errors"
	"log"
	"strconv"
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"

	"userstyles.world/handlers/jwt"
	"userstyles.world/models"
	"userstyles.world/utils"
)

func OAuthSettingsGet(c *fiber.Ctx) error {
	u, _ := jwt.User(c)
	id := c.Params("id")

	var method string
	var oauth *models.APIOAuth
	var err error
	if isEdit := id != ""; isEdit {
		method = "edit"
		oauth, err = models.GetOAuthByID(id)
	} else {
		method = "add"
	}

	if err != nil {
		log.Printf("Failed to oauth, err: %#+v\n", err)
		return c.Render("err", fiber.Map{
			"Title": "Internal server error",
			"User":  u,
		})
	}

	if isEdit {
		if u.ID != oauth.UserID {
			return c.Render("err", fiber.Map{
				"Title": "Users don't match",
				"User":  u,
			})
		}
	}
	oauths, err := models.ListOAuthsOfUser(u.Username)
	if err != nil {
		if isEdit {
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
			return c.Render("oauth_settings", arguments)
		}
		return c.Render("oauth_settings", fiber.Map{
			"Title":  "OAuth Settings",
			"User":   u,
			"OAuths": nil,
			"Method": method,
		})
	}

	if isEdit {
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
		return c.Render("oauth_settings", arguments)
	}
	return c.Render("oauth_settings", fiber.Map{
		"Title":  "OAuth Settings",
		"User":   u,
		"OAuths": oauths,
		"Method": method,
	})
}

func OAuthSettingsPost(c *fiber.Ctx) error {
	u, _ := jwt.User(c)

	q := models.OAuth{
		Name:        c.FormValue("name"),
		Description: c.FormValue("description"),
		RedirectURI: strings.TrimSuffix(c.FormValue("redirect_uri"), "/"),
		Scopes: utils.Filter([]string{"style", "user"}, func(name interface{}) bool {
			return c.FormValue(name.(string)) == "on"
		}).([]string),
		UserID: u.ID,
	}

	if err := utils.Validate().StructPartial(q, "Name", "Description"); err != nil {
		var validationError validator.ValidationErrors
		if ok := errors.As(err, &validationError); ok {
			log.Println("Validation errors:", validationError)
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
			Render("oauth_settings", arguments)
	}

	var err error
	var dbOAuth *models.OAuth
	if id := c.Params("id"); id != "" {
		err = models.UpdateOAuth(&q, id)
	} else {
		q.ClientID = utils.UnsafeString((utils.RandStringBytesMaskImprSrcUnsafe(32)))
		q.ClientSecret = utils.UnsafeString((utils.RandStringBytesMaskImprSrcUnsafe(128)))
		dbOAuth, err = models.CreateOAuth(&q)
	}

	if err != nil {
		log.Println("Updating style failed, err:", err)
		return c.Render("err", fiber.Map{
			"Title": "Internal server error.",
			"User":  u,
		})
	}

	OAuthID := strconv.FormatUint(uint64(dbOAuth.ID), 10)
	return c.Redirect("/oauth_settings/"+OAuthID, fiber.StatusSeeOther)
}
