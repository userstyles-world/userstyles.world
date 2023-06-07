package style

import (
	"strconv"
	"strings"

	"github.com/gofiber/fiber/v2"

	"userstyles.world/handlers/jwt"
	"userstyles.world/models"
	"userstyles.world/modules/cache"
	"userstyles.world/modules/images"
	"userstyles.world/modules/log"
	"userstyles.world/modules/search"
	"userstyles.world/modules/util"
	"userstyles.world/utils"
)

func EditGet(c *fiber.Ctx) error {
	u, _ := jwt.User(c)
	c.Locals("User", u)
	c.Locals("Title", "Edit userstyle")

	i, err := c.ParamsInt("id")
	if err != nil || i < 1 {
		c.Locals("Title", "Invalid style ID")
		return c.Status(fiber.StatusBadRequest).Render("err", fiber.Map{})
	}

	s, err := models.GetStyleFromAuthor(i, int(u.ID))
	if err != nil {
		c.Locals("Title", "Style not found")
		return c.Status(fiber.StatusNotFound).Render("err", fiber.Map{})
	}
	c.Locals("Style", s)

	return c.Render("style/edit", fiber.Map{})
}

func EditPost(c *fiber.Ctx) error {
	u, _ := jwt.User(c)
	c.Locals("User", u)
	c.Locals("Title", "Edit userstyle")

	i, err := c.ParamsInt("id")
	if err != nil || i < 1 {
		c.Locals("Title", "Invalid style ID")
		return c.Status(fiber.StatusBadRequest).Render("err", fiber.Map{})
	}
	id := c.Params("id")

	s, err := models.GetStyleFromAuthor(i, int(u.ID))
	if err != nil {
		c.Locals("Title", "Style not found")
		return c.Status(fiber.StatusNotFound).Render("err", fiber.Map{})
	}

	s.Name = strings.TrimSpace(c.FormValue("name"))
	s.Description = strings.TrimSpace(c.FormValue("description"))
	s.Notes = strings.TrimSpace(c.FormValue("notes"))
	s.Code = util.RemoveUpdateURL(c.FormValue("code"))
	s.Category = strings.TrimSpace(c.FormValue("category"))
	s.Homepage = strings.TrimSpace(c.FormValue("homepage"))
	s.License = strings.TrimSpace(c.FormValue("license", "No License"))
	s.MirrorURL = strings.TrimSpace(c.FormValue("mirrorURL"))
	s.ImportPrivate = c.FormValue("importPrivate") == "on"
	s.MirrorPrivate = c.FormValue("mirrorPrivate") == "on"
	s.MirrorCode = c.FormValue("mirrorCode") == "on"
	s.MirrorMeta = c.FormValue("mirrorMeta") == "on"
	c.Locals("Style", s)

	m, err := s.Validate(utils.Validate(), false)
	if err != nil {
		c.Locals("err", m)
		c.Locals("Error", "Incorrect userstyle data was entered. Please review the fields bellow.")
		return c.Status(fiber.StatusBadRequest).Render("style/edit", fiber.Map{})
	}

	// Prevent adding multiples of the same style.
	err = models.CheckDuplicateStyle(&s)
	if err != nil {
		c.Locals("dupName", "Duplicate userstyle names aren't allowed.")
		c.Locals("Error", "Incorrect userstyle data was entered. Please review the fields bellow.")
		return c.Status(fiber.StatusBadRequest).Render("style/edit", fiber.Map{})
	}

	// Check for new preview image.
	file, _ := c.FormFile("preview")
	version := strconv.Itoa(s.PreviewVersion + 1)
	preview := strings.TrimSpace(c.FormValue("previewURL"))
	if file != nil || (preview != s.Preview && preview != "") {
		if err := images.Generate(file, id, version, s.Preview, preview); err != nil {
			log.Warn.Println("Error:", err)
		} else {
			s.PreviewVersion++
			s.SetPreview()
		}
	} else if preview == "" {
		// TODO: Figure out a better UI/UX for this functionality.  ATM, one has
		// to set "Preview image URL" field to be empty in order to remove it.
		s.Preview = ""
	}

	// TODO: Split updates into sections.
	if err = models.SelectUpdateStyle(s); err != nil {
		log.Database.Printf("Failed to update style %d: %s\n", s.ID, err)
		c.Locals("Error", "Failed to update userstyle in database. Please try again.")
		return c.Status(fiber.StatusBadRequest).Render("style/edit", fiber.Map{})
	}

	// TODO: Move to code section once we refactor this messy logic.
	// sftodo: is it needed?
	cache.LRU.Remove(id)

	if err = search.IndexStyle(s.ID); err != nil {
		log.Warn.Printf("Failed to re-index style %d: %s\n", s.ID, err)
	}

	return c.Redirect("/style/"+id, fiber.StatusSeeOther)
}
