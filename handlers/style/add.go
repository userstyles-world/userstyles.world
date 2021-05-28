package style

import (
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"os"
	"strconv"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/vednoc/go-usercss-parser"

	"userstyles.world/database"
	"userstyles.world/handlers/jwt"
	"userstyles.world/images"
	"userstyles.world/models"
)

func CreateGet(c *fiber.Ctx) error {
	u, _ := jwt.User(c)

	return c.Render("add", fiber.Map{
		"Title":  "Add userstyle",
		"User":   u,
		"Method": "add",
	})
}

func CreatePost(c *fiber.Ctx) error {
	u, _ := jwt.User(c)

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
		return c.Render("add", fiber.Map{
			"Title":  "Add userstyle",
			"User":   u,
			"Style":  s,
			"Method": "add",
			"Errors": errs,
		})
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

	if image != nil {
		styleID := strconv.FormatUint(uint64(s.ID), 10)
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

	return c.Redirect(fmt.Sprintf("/style/%d", int(s.ID)), fiber.StatusSeeOther)
}
