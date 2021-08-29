package style

import (
	"io"
	"mime/multipart"
	"os"
	"strings"

	"github.com/userstyles-world/fiber/v2"

	"userstyles.world/handlers/jwt"
	"userstyles.world/models"
	"userstyles.world/modules/database"
	"userstyles.world/modules/images"
	"userstyles.world/modules/log"
)

func EditGet(c *fiber.Ctx) error {
	u, _ := jwt.User(c)
	p := c.Params("id")

	s, err := models.GetStyleByID(p)
	if err != nil {
		return c.Render("err", fiber.Map{
			"Title": "Style not found",
			"User":  u,
		})
	}

	// Check if logged-in user matches style author.
	if u.ID != s.UserID {
		return c.Render("err", fiber.Map{
			"Title": "User and style author don't match",
			"User":  u,
		})
	}

	return c.Render("style/create", fiber.Map{
		"Title":  "Edit userstyle",
		"Method": "edit",
		"User":   u,
		"Styles": s,
	})
}

func EditPost(c *fiber.Ctx) error {
	u, _ := jwt.User(c)
	styleID, t := c.Params("id"), new(models.Style)

	s, err := models.GetStyleByID(styleID)
	if err != nil {
		return c.Render("err", fiber.Map{
			"Title": "Style not found",
			"User":  u,
		})
	}

	// Check if logged-in user matches style author.
	if u.ID != s.UserID {
		return c.Render("err", fiber.Map{
			"Title": "Users don't match",
			"User":  u,
		})
	}

	// Check if userstyle name is empty.
	if strings.TrimSpace(c.FormValue("name")) == "" {
		return c.Render("err", fiber.Map{
			"Title": "Style name can't be empty",
			"User":  u,
		})
	}

	q := models.Style{
		Name:        c.FormValue("name"),
		Description: c.FormValue("description"),
		Notes:       c.FormValue("notes"),
		Code:        c.FormValue("code"),
		Homepage:    c.FormValue("homepage"),
		Preview:     c.FormValue("previewUrl"),
		License:     strings.TrimSpace(c.FormValue("license", "No License")),
		Category:    strings.TrimSpace(c.FormValue("category", "unset")),
		MirrorURL:   strings.TrimSpace(c.FormValue("mirrorURL")),
		UserID:      u.ID,
	}

	var image multipart.File
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

	err = database.Conn.
		Model(t).
		Where("id", styleID).
		Updates(q).
		// GORM doesn't update non-zero values in structs.
		Update("mirror_code", c.FormValue("mirrorCode") == "on").
		Update("mirror_meta", c.FormValue("mirrorMeta") == "on").
		Error

	if err != nil {
		log.Warn.Println("Failed to update style:", err.Error())
		return c.Render("err", fiber.Map{
			"Title": "Internal server error.",
			"User":  u,
		})
	}

	go func(image multipart.File, style *models.Style, styleID, preview string) {
		isLocal := false
		style.Preview = "https://userstyles.world/api/style/preview/" + styleID + ".jpeg"

		var err error
		if image != nil {
			isLocal = true
			data, _ := io.ReadAll(image)
			err = os.WriteFile(images.CacheFolder+styleID+".original", data, 0o600)
			if err != nil {
				log.Warn.Println("Failed to write image:", err.Error())
				return
			}
		}
		err = images.GenerateImagesForStyle(styleID, preview, isLocal)
		if err != nil {
			s.Preview = ""
			log.Warn.Println("Failed to generate images:", err.Error())
			return
		}

		err = database.Conn.
			Model(new(models.Style)).
			Where("id", styleID).
			Updates(style).
			Error
		if err != nil {
			log.Warn.Println("Failed to update style:", err.Error())
		}
	}(image, &q, styleID, q.Preview)

	return c.Redirect("/style/"+c.Params("id"), fiber.StatusSeeOther)
}
