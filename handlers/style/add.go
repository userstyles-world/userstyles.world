package style

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt"

	jwtware "userstyles.world/handlers/jwt"
	"userstyles.world/models"
	"userstyles.world/modules/config"
	"userstyles.world/modules/images"
	"userstyles.world/modules/log"
	"userstyles.world/modules/util"
	"userstyles.world/modules/validator"
)

func CreateGet(c *fiber.Ctx) error {
	u, _ := jwtware.User(c)
	c.Locals("User", u)
	c.Locals("Title", "Add userstyle")

	return c.Render("style/add", fiber.Map{})
}

func CreatePost(c *fiber.Ctx) error {
	u, _ := jwtware.User(c)
	c.Locals("User", u)
	c.Locals("Title", "Add userstyle")

	secureToken := c.Query("token")
	if secureToken != "" {
		c.Locals("SecureToken", secureToken)
	}
	oauthID := c.Query("oauthID")
	if oauthID != "" {
		c.Locals("OAuthID", oauthID)
		c.Locals("Method", "api")
	}

	s := &models.Style{
		Name:        strings.TrimSpace(c.FormValue("name")),
		Description: strings.TrimSpace(c.FormValue("description")),
		Notes:       strings.TrimSpace(c.FormValue("notes")),
		Homepage:    strings.TrimSpace(c.FormValue("homepage")),
		License:     strings.TrimSpace(c.FormValue("license", "No License")),
		Code:        strings.TrimSpace(util.RemoveUpdateURL(c.FormValue("code"))),
		Category:    strings.TrimSpace(c.FormValue("category")),
		UserID:      u.ID,
	}
	c.Locals("Style", s)

	// Get previewURL
	preview := c.FormValue("previewURL")
	c.Locals("PreviewURL", preview)

	m, err := s.Validate(validator.V, true)
	if err != nil {
		c.Locals("err", m)
		c.Locals("Error", "Incorrect userstyle data was entered. Please review the fields bellow.")
		return c.Status(fiber.StatusBadRequest).Render("style/add", fiber.Map{})
	}

	// Prevent adding multiples of the same style.
	err = models.CheckDuplicateStyle(s)
	if err != nil {
		c.Locals("dupName", "Duplicate userstyle names aren't allowed.")
		c.Locals("Error", "Incorrect userstyle data was entered. Please review the fields bellow.")
		return c.Status(fiber.StatusBadRequest).Render("style/add", fiber.Map{})
	}

	s, err = models.CreateStyle(s)
	if err != nil {
		log.Warn.Println("Failed to create style:", err)
		c.Locals("Error", "Failed to add userstyle to database. Please try again.")
		return c.Status(fiber.StatusBadRequest).Render("style/add", fiber.Map{})
	}

	err = models.SaveStyleCode(strconv.Itoa(int(s.ID)), s.Code)
	if err != nil {
		log.Warn.Printf("kind=code id=%v err=%q\n", s.ID, err)
	}

	// Check preview image.
	file, _ := c.FormFile("preview")
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

	unsealedText, err := util.DecryptText(secureToken, util.AEADOAuthp, config.Secrets)
	if err != nil {
		log.Warn.Println("Failed to unseal JWT text:", err.Error())
		return c.Status(500).
			JSON(fiber.Map{
				"data": "Error: Please notify the UserStyles.world admins.",
			})
	}

	token, err := jwt.Parse(unsealedText, util.OAuthPJwtKeyFunction)
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

	jwtToken, err := util.NewJWT().
		SetClaim("state", state).
		SetClaim("userID", u.ID).
		SetClaim("styleID", style.ID).
		SetExpiration(time.Now().Add(time.Minute * 10)).
		GetSignedString(util.OAuthPSigningKey)
	if err != nil {
		log.Warn.Println("Failed to create a JWT Token:", err.Error())
		return c.Status(500).
			JSON(fiber.Map{
				"data": "Error: Please notify the UserStyles.world admins.",
			})
	}

	returnCode := "?code=" + util.EncryptText(jwtToken, util.AEADOAuthp, config.Secrets)
	returnCode += "&style_id=" + styleID
	if state != "" {
		returnCode += "&state=" + state
	}

	return c.Redirect(oauth.RedirectURI + "/" + returnCode)
}
