package style

import (
	"fmt"
	"mime/multipart"
	"strings"

	"github.com/gofiber/fiber/v2"

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
	id := c.Params("id")

	s, err := models.GetStyleByID(id)
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

	// Prepare initial data.
	m := map[string]interface{}{
		"id":      id,
		"name":    strings.TrimSpace(c.FormValue("name")),
		"preview": strings.TrimSpace(c.FormValue("previewURL")),
	}

	// Check if userstyle name is empty.
	// TODO: Implement proper validation.
	if m["name"] == "" {
		return c.Render("err", fiber.Map{
			"Title": "Style name can't be empty",
			"User":  u,
		})
	}

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
	if image != nil || s.Preview != m["preview"] {
		err = images.Generate(image, id, m["preview"].(string))
		if err != nil {
			log.Warn.Printf("Failed to generate images for %d: %s\n", s.ID, err.Error())
			m["preview"] = ""
		} else {
			m["preview"] = fmt.Sprintf("%s/api/style/preview/%s.jpeg", config.BaseURL, id)
		}
	}

	// Add the rest of the data.
	m["description"] = strings.TrimSpace(c.FormValue("description"))
	m["notes"] = strings.TrimSpace(c.FormValue("notes"))
	m["code"] = strings.TrimSpace(c.FormValue("code"))
	m["homepage"] = strings.TrimSpace(c.FormValue("homepage"))
	m["license"] = strings.TrimSpace(c.FormValue("license", "No License"))
	m["category"] = strings.TrimSpace(c.FormValue("category", "unset"))
	m["mirror_url"] = strings.TrimSpace(c.FormValue("mirrorURL"))
	m["mirror_code"] = c.FormValue("mirrorCode") == "on"
	m["mirror_meta"] = c.FormValue("mirrorMeta") == "on"

	// TODO: Split updates into sections.
	err = database.Conn.Debug().Model(models.Style{}).Where("id", id).Updates(m).Error
	if err != nil {
		log.Warn.Printf("Failed to update style %d: %v\n", s.ID, err)
		return c.Render("err", fiber.Map{
			"Title": "Failed to update userstyle",
			"User":  u,
		})
	}

	// TODO: Move to code section once we refactor this messy logic.
	cache.LRU.Remove(id)

	return c.Redirect("/style/"+id, fiber.StatusSeeOther)
}
