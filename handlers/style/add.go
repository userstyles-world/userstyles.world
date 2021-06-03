package style

import (
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/form3tech-oss/jwt-go"
	"github.com/gofiber/fiber/v2"
	"github.com/vednoc/go-usercss-parser"

	"userstyles.world/database"
	jwtware "userstyles.world/handlers/jwt"
	"userstyles.world/images"
	"userstyles.world/models"
	"userstyles.world/search"
	"userstyles.world/utils"
)

func CreateGet(c *fiber.Ctx) error {
	u, _ := jwtware.User(c)

	return c.Render("add", fiber.Map{
		"Title":  "Add userstyle",
		"User":   u,
		"Method": "add",
	})
}

func CreatePost(c *fiber.Ctx) error {
	u, _ := jwtware.User(c)
	secureToken, OAuthID := c.Query("token"), c.Query("oauthID")

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
			"Style":  s,
			"Method": "add",
			"Errors": errs,
		}
		if OAuthID != "" {
			arguments["Method"] = "add_api"
			arguments["OAuthID"] = OAuthID
			arguments["SecureToken"] = secureToken
		}
		return c.Render("add", arguments)
	}

	// Prevent adding multiples of the same style.
	err := models.CheckDuplicateStyle(database.DB, s)
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
				log.Println("Opening image , err:", err)
				return c.Render("err", fiber.Map{
					"Title": "Internal server error.",
					"User":  u,
				})
			}
		}
	}
	s, err = models.CreateStyle(database.DB, s)
	if err != nil {
		log.Println("Style creation failed, err:", err)
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
			log.Println("Style creation failed, err:", err)
			return c.Render("err", fiber.Map{
				"Title": "Internal server error.",
				"User":  u,
			})
		}
		if s.Preview == "" {
			s.Preview = "https://userstyles.world/api/preview/" + styleID + ".jpeg"
			database.DB.
				Model(new(models.Style)).
				Where("id", styleID).
				Updates(s)
		}
	}

	if err = search.IndexStyle(s.ID); err != nil {
		log.Printf("Re-indexing style %d failed, err: %s", s.ID, err.Error())
	}

	if OAuthID != "" {
		return handleAPIStyle(c, secureToken, OAuthID, styleID, s)
	}

	return c.Redirect(fmt.Sprintf("/style/%d", int(s.ID)), fiber.StatusSeeOther)
}

func handleAPIStyle(c *fiber.Ctx, secureToken, oauthID, styleID string, style *models.Style) error {
	u, _ := jwtware.User(c)

	OAuth, err := models.GetOAuthByID(database.DB, oauthID)
	if err != nil || OAuth.ID == 0 {
		return c.Status(400).
			JSON(fiber.Map{
				"data": "Incorrect oauthID specified",
			})
	}

	unsealedText, err := utils.DecodePreparedText(secureToken, utils.AEAD_OAUTHP)
	if err != nil {
		log.Println("Error: Couldn't unseal JWT Token:", err.Error())
		return c.Status(500).
			JSON(fiber.Map{
				"data": "JWT Token error, please notify the admins.",
			})
	}

	token, err := jwt.Parse(unsealedText, utils.OAuthPJwtKeyFunction)
	if err != nil || !token.Valid {
		log.Println("Error: Couldn't unseal JWT Token:", err.Error())
		return c.Status(500).
			JSON(fiber.Map{
				"data": "JWT Token error, please notify the admins.",
			})
	}
	claims := token.Claims.(jwt.MapClaims)

	userID, ok := claims["userID"].(float64)
	if !ok || userID != float64(u.ID) {
		log.Println("WARNING!: Invalid userID")
		return c.Status(500).
			JSON(fiber.Map{
				"data": "JWT Token error, please notify the admins.",
			})
	}

	state, ok := claims["state"].(string)
	if !ok {
		log.Println("WARNING!: Invalid state")
		return c.Status(500).
			JSON(fiber.Map{
				"data": "JWT Token error, please notify the admins.",
			})
	}

	if style.UserID != u.ID {
		log.Println("WARNING!: Invalid style's user ID")
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
		log.Println("Error: Couldn't create JWT Token:", err.Error())
		return c.Status(500).
			JSON(fiber.Map{
				"data": "JWT Token error, please notify the admins.",
			})
	}

	returnCode := "?code=" + utils.PrepareText(jwt, utils.AEAD_OAUTHP)
	returnCode += "&style_id=" + styleID
	if state != "" {
		returnCode += "&state=" + state
	}

	return c.Redirect(OAuth.RedirectURI + "/" + returnCode)
}
