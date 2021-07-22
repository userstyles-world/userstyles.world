package style

import (
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/form3tech-oss/jwt-go"
	"github.com/gofiber/fiber/v2"
	"github.com/vednoc/go-usercss-parser"

	jwtware "userstyles.world/handlers/jwt"
	"userstyles.world/models"
	"userstyles.world/modules/config"
	"userstyles.world/modules/database"
	"userstyles.world/modules/images"
	"userstyles.world/modules/log"
	"userstyles.world/search"
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

	// Check if userstyle name is empty.
	if strings.TrimSpace(c.FormValue("name")) == "" {
		return c.Render("err", fiber.Map{
			"Title": "Style name can't be empty",
			"User":  u,
		})
	}

	s := &models.Style{
		Name:        c.FormValue("name"),
		Description: c.FormValue("description"),
		Notes:       c.FormValue("notes"),
		Homepage:    c.FormValue("homepage"),
		Preview:     c.FormValue("previewUrl"),
		Code:        c.FormValue("code"),
		License:     strings.TrimSpace(c.FormValue("license", "No License")),
		Category:    strings.TrimSpace(c.FormValue("category", "unset")),
		UserID:      u.ID,
	}

	code := usercss.ParseFromString(c.FormValue("code"))
	if errs := usercss.BasicMetadataValidation(code); errs != nil {
		arguments := fiber.Map{
			"Title":  "Add userstyle",
			"User":   u,
			"Styles": s,
			"Method": "add",
			"Errors": errs,
		}
		if oauthID != "" {
			arguments["Method"] = "add_api"
			arguments["OAuthID"] = oauthID
			arguments["SecureToken"] = secureToken
		}
		return c.Render("style/create", arguments)
	}

	// Prevent adding multiples of the same style.
	err := models.CheckDuplicateStyle(s)
	if err != nil {
		return c.Render("err", fiber.Map{
			"Title": err,
			"User":  u,
		})
	}

	var image multipart.File
	if s.Preview == "" {
		if ff, _ := c.FormFile("preview"); ff != nil {
			image, err = ff.Open()
			if err != nil {
				log.Warn.Println("Failed to open image:", err.Error())
				return c.Render("err", fiber.Map{
					"Title": "Internal server error.",
					"User":  u,
				})
			}
		}
	}
	s, err = models.CreateStyle(s)
	if err != nil {
		log.Warn.Println("Failed to create style:", err.Error())
		return c.Render("err", fiber.Map{
			"Title": "Internal server error.",
			"User":  u,
		})
	}

	styleID := strconv.FormatUint(uint64(s.ID), 10)
	if image != nil {
		data, _ := io.ReadAll(image)
		err = os.WriteFile(images.CacheFolder+styleID+".original", data, 0o600)
		if err != nil {
			log.Warn.Println("Failed to write image:", err.Error())
			return c.Render("err", fiber.Map{
				"Title": "Internal server error.",
				"User":  u,
			})
		}
		if s.Preview == "" {
			s.Preview = "https://userstyles.world/api/style/preview/" + styleID + ".jpeg"
			database.Conn.
				Model(new(models.Style)).
				Where("id", styleID).
				Updates(s)
		}
	}

	go func(style *models.Style) {
		if err = search.IndexStyle(style.ID); err != nil {
			log.Warn.Printf("Failed to re-index style %d: %\ns", style.ID, err.Error())
		}
	}(s)

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
				"data": "JWT Token error, please notify the admins.",
			})
	}

	token, err := jwt.Parse(unsealedText, utils.OAuthPJwtKeyFunction)
	if err != nil || !token.Valid {
		log.Warn.Println("Failed to unseal JWT token:", err.Error())
		return c.Status(500).
			JSON(fiber.Map{
				"data": "JWT Token error, please notify the admins.",
			})
	}
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		log.Warn.Println("Failed to parse JWT claims:", err.Error())
		return c.Status(500).
			JSON(fiber.Map{
				"data": "JWT Token error, please notify the admins.",
			})
	}

	userID, ok := claims["userID"].(float64)
	if !ok || userID != float64(u.ID) {
		log.Warn.Println("Failed to get userID from parsed token.")
		return c.Status(500).
			JSON(fiber.Map{
				"data": "JWT Token error, please notify the admins.",
			})
	}

	state, ok := claims["state"].(string)
	if !ok {
		log.Warn.Println("Invalid JWT state.")
		return c.Status(500).
			JSON(fiber.Map{
				"data": "JWT Token error, please notify the admins.",
			})
	}

	if style.UserID != u.ID {
		log.Warn.Println("Failed to match style author and userID.")
		return c.Status(500).
			JSON(fiber.Map{
				"data": "JWT Token error, please notify the admins.",
			})
	}

	jwt, err := utils.NewJWTToken().
		SetClaim("state", state).
		SetClaim("userID", u.ID).
		SetClaim("styleID", style.ID).
		SetExpiration(time.Now().Add(time.Minute * 10)).
		GetSignedString(utils.OAuthPSigningKey)
	if err != nil {
		log.Warn.Println("Failed to create a JWT Token:", err.Error())
		return c.Status(500).
			JSON(fiber.Map{
				"data": "JWT Token error, please notify the admins.",
			})
	}

	returnCode := "?code=" + utils.EncryptText(jwt, utils.AEADOAuthp, config.ScrambleConfig)
	returnCode += "&style_id=" + styleID
	if state != "" {
		returnCode += "&state=" + state
	}

	return c.Redirect(oauth.RedirectURI + "/" + returnCode)
}
