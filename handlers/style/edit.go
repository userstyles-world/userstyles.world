package style

import (
	"mime/multipart"
	"strings"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"

	"userstyles.world/handlers/jwt"
	"userstyles.world/models"
	"userstyles.world/modules/cache"
	"userstyles.world/modules/config"
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
		Model:       gorm.Model{ID: s.ID},
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

	// TODO: Move to code section once we refactor this messy logic.
	cache.LRU.Remove(styleID)

	// Check for new image.
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

	// Check for new preview image.
	if image != nil || s.Preview != q.Preview {
		err = images.Generate(image, styleID, q.Preview)
		if err != nil {
			log.Warn.Printf("Failed to generate images for %d: %s\n", s.ID, err.Error())
			q.Preview = ""
		} else {
			q.Preview = config.BaseURL + "/api/style/preview/" + styleID + ".jpeg"
		}
	}

	if err = q.UpdateColumn("preview", q.Preview); err != nil {
		log.Warn.Printf("Failed to update preview image for %s: %s\n", styleID, err.Error())
		return c.Render("err", fiber.Map{
			"Title": "Failed to update preview image",
			"User":  u,
		})
	}

	return c.Redirect("/style/"+c.Params("id"), fiber.StatusSeeOther)
}
