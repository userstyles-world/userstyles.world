package style

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt"
	"github.com/vednoc/go-usercss-parser"

	jwtware "userstyles.world/handlers/jwt"
	"userstyles.world/models"
	"userstyles.world/modules/config"
	"userstyles.world/modules/images"
	"userstyles.world/modules/log"
	"userstyles.world/modules/search"
	"userstyles.world/modules/util"
	"userstyles.world/utils"
)

func CreateGet(c *fiber.Ctx) error {
	u, _ := jwtware.User(c)

	return c.Render("style/create", fiber.Map{
		"Title":  "Add userstyle",
		"User":   u,
		"Method": "add",
	})
}

func CreatePost(c *fiber.Ctx) error {
	u, _ := jwtware.User(c)
	secureToken, oauthID := c.Query("token"), c.Query("oauthID")

	s := &models.Style{
		Name:        strings.TrimSpace(c.FormValue("name")),
		Description: strings.TrimSpace(c.FormValue("description")),
		Notes:       strings.TrimSpace(c.FormValue("notes")),
		Homepage:    strings.TrimSpace(c.FormValue("homepage")),
		License:     strings.TrimSpace(c.FormValue("license", "No License")),
		Category:    strings.TrimSpace(c.FormValue("category")),
		UserID:      u.ID,
	}

	m, msg, err := s.Validate(utils.Validate())
	if err != nil {
		return c.Render("style/create", fiber.Map{
			"Title":  "Add userstyle",
			"User":   u,
			"Styles": s,
			"Method": "add",
			"Error":  msg,
			"err":    m,
		})
	}

	s.Code = strings.TrimSpace(util.RemoveUpdateURL(c.FormValue("code")))

	uc := new(usercss.UserCSS)
	if err := uc.Parse(s.Code); err != nil {
		// TODO: Fix this in UserCSS parser.
		e := err.Error()
		msg := strings.ToUpper(string(e[0])) + e[1:] + "."

		arguments := fiber.Map{
			"Title":   "Add userstyle",
			"User":    u,
			"Styles":  s,
			"Method":  "add",
			"Error":   "Invalid source code.",
			"errCode": msg,
		}
		if oauthID != "" {
			arguments["Method"] = "add_api"
			arguments["OAuthID"] = oauthID
			arguments["SecureToken"] = secureToken
		}
		return c.Render("style/create", arguments)
	}
	if errs := uc.Validate(); errs != nil {
		arguments := fiber.Map{
			"Title":  "Add userstyle",
			"User":   u,
			"Styles": s,
			"Method": "add",
			"Error":  "Missing mandatory fields in source code.",
			"errors": errs,
		}
		if oauthID != "" {
			arguments["Method"] = "add_api"
			arguments["OAuthID"] = oauthID
			arguments["SecureToken"] = secureToken
		}
		return c.Render("style/create", arguments)
	}

	// Prevent broken traditional userstyles.
	// TODO: Remove a week or two after Stylus v1.5.20 is released.
	if len(uc.MozDocument) == 0 {
		return c.Render("err", fiber.Map{
			"Title":  "Bad style format",
			"Stylus": "Your style is affected by a bug in Stylus integration.",
			"User":   u,
		})
	}

	// Prevent adding multiples of the same style.
	err = models.CheckDuplicateStyle(s)
	if err != nil {
		return c.Render("err", fiber.Map{
			"Title": err,
			"User":  u,
		})
	}

	s, err = models.CreateStyle(s)
	if err != nil {
		log.Warn.Println("Failed to create style:", err.Error())
		return c.Render("err", fiber.Map{
			"Title": "Internal server error",
			"User":  u,
		})
	}

	// Check preview image.
	file, _ := c.FormFile("preview")
	preview := c.FormValue("previewURL")
	styleID := strconv.FormatUint(uint64(s.ID), 10)
	if file != nil || preview != "" {
		if err = images.Generate(file, styleID, "0", "", preview); err != nil {
			log.Warn.Println("Error:", err)
			s.Preview = ""
		} else {
			s.SetPreview()
			if err = s.UpdateColumn("preview", s.Preview); err != nil {
				log.Warn.Printf("Failed to update preview for %d: %s\n", s.ID, err)
			}
		}
	}

	if err = search.IndexStyle(s.ID); err != nil {
		log.Warn.Printf("Failed to index style %d: %s\n", s.ID, err)
	}

	if oauthID != "" {
		return handleAPIStyle(c, secureToken, oauthID, styleID, s)
	}

	return c.Redirect(fmt.Sprintf("/style/%d", int(s.ID)), fiber.StatusSeeOther)
}

func handleAPIStyle(c *fiber.Ctx, secureToken, oauthID, styleID string, style *models.Style) error {
	u, _ := jwtware.User(c)

	oauth, err := models.GetOAuthByID(oauthID)
	if err != nil || oauth.ID == 0 {
		return c.Status(400).
			JSON(fiber.Map{
				"data": "Incorrect oauthID specified",
			})
	}

	unsealedText, err := utils.DecryptText(secureToken, utils.AEADOAuthp, config.ScrambleConfig)
	if err != nil {
		log.Warn.Println("Failed to unseal JWT text:", err.Error())
		return c.Status(500).
			JSON(fiber.Map{
				"data": "Error: Please notify the UserStyles.world admins.",
			})
	}

	token, err := jwt.Parse(unsealedText, utils.OAuthPJwtKeyFunction)
	if err != nil || !token.Valid {
		log.Warn.Println("Failed to unseal JWT token:", err.Error())
		return c.Status(500).
			JSON(fiber.Map{
				"data": "Error: Please notify the UserStyles.world admins.",
			})
	}
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		log.Warn.Println("Failed to parse JWT claims:", err.Error())
		return c.Status(500).
			JSON(fiber.Map{
				"data": "Error: Please notify the UserStyles.world admins.",
			})
	}

	userID, ok := claims["userID"].(float64)
	if !ok || userID != float64(u.ID) {
		log.Warn.Println("Failed to get userID from parsed token.")
		return c.Status(500).
			JSON(fiber.Map{
				"data": "Error: Please notify the UserStyles.world admins.",
			})
	}

	state, ok := claims["state"].(string)
	if !ok {
		log.Warn.Println("Invalid JWT state.")
		return c.Status(500).
			JSON(fiber.Map{
				"data": "Error: Please notify the UserStyles.world admins.",
			})
	}

	if style.UserID != u.ID {
		log.Warn.Println("Failed to match style author and userID.")
		return c.Status(500).
			JSON(fiber.Map{
				"data": "Error: Please notify the UserStyles.world admins.",
			})
	}

	jwtToken, err := utils.NewJWTToken().
		SetClaim("state", state).
		SetClaim("userID", u.ID).
		SetClaim("styleID", style.ID).
		SetExpiration(time.Now().Add(time.Minute * 10)).
		GetSignedString(utils.OAuthPSigningKey)
	if err != nil {
		log.Warn.Println("Failed to create a JWT Token:", err.Error())
		return c.Status(500).
			JSON(fiber.Map{
				"data": "Error: Please notify the UserStyles.world admins.",
			})
	}

	returnCode := "?code=" + utils.EncryptText(jwtToken, utils.AEADOAuthp, config.ScrambleConfig)
	returnCode += "&style_id=" + styleID
	if state != "" {
		returnCode += "&state=" + state
	}

	return c.Redirect(oauth.RedirectURI + "/" + returnCode)
}
