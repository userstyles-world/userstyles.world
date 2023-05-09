package style

import (
	"strconv"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/vednoc/go-usercss-parser"

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

	i, err := c.ParamsInt("id")
	if err != nil || i < 1 {
		return c.Render("err", fiber.Map{
			"User":  u,
			"Title": "Invalid style ID",
		})
	}
	id := c.Params("id")

	args := fiber.Map{
		"User":   u,
		"Title":  "Edit userstyle",
		"Method": "edit", // TODO: Remove later.
	}

	s, err := models.GetStyle(id)
	if err != nil {
		args["Title"] = "Style not found"
		return c.Render("err", args)
	}

	// Check if logged-in user matches style author.
	if u.ID != s.UserID {
		args["Title"] = "User and style author don't match"
		return c.Render("err", args)
	}

	s.Name = strings.TrimSpace(c.FormValue("name"))
	s.Description = strings.TrimSpace(c.FormValue("description"))
	s.Notes = strings.TrimSpace(c.FormValue("notes"))
	s.Code = util.RemoveUpdateURL(c.FormValue("code"))
	s.Category = strings.TrimSpace(c.FormValue("category"))
	args["Styles"] = s

	m, msg, err := s.Validate(utils.Validate())
	if err != nil {
		args["Error"] = msg
		args["err"] = m

		return c.Render("style/create", args)
	}

	var uc usercss.UserCSS
	if err := uc.Parse(s.Code); err != nil {
		// TODO: Fix this in UserCSS parser.
		e := err.Error()
		msg := strings.ToUpper(string(e[0])) + e[1:] + "."

		args["Error"] = "Invalid source code."
		args["errCode"] = msg
		return c.Render("style/create", args)
	}
	if errs := uc.Validate(); errs != nil {
		args["Error"] = "Missing mandatory fields in source code."
		args["errors"] = errs
		return c.Render("style/create", args)
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

	// Update the other fields with new data.
	s.Homepage = strings.TrimSpace(c.FormValue("homepage"))
	s.License = strings.TrimSpace(c.FormValue("license", "No License"))
	s.MirrorURL = strings.TrimSpace(c.FormValue("mirrorURL"))
	s.MirrorCode = c.FormValue("mirrorCode") == "on"
	s.MirrorMeta = c.FormValue("mirrorMeta") == "on"
	args["Styles"] = s // NOTE: Add new data.

	// TODO: Split updates into sections.
	if err = models.SelectUpdateStyle(s); err != nil {
		log.Database.Printf("Failed to update style %d: %s\n", s.ID, err)
		args["Title"] = "Failed to update userstyle"
		args["Error"] = "Failed to update style in database."
		return c.Render("style/create", args)
	}

	// TODO: Move to code section once we refactor this messy logic.
	cache.LRU.Remove(id)

	if err = search.IndexStyle(s.ID); err != nil {
		log.Warn.Printf("Failed to re-index style %d: %s\n", s.ID, err)
	}

	return c.Redirect("/style/"+id, fiber.StatusSeeOther)
}
