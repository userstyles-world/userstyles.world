package style

import (
	"strconv"

	"github.com/gofiber/fiber/v2"

	"userstyles.world/handlers/jwt"
	"userstyles.world/models"
	"userstyles.world/modules/cache"
	"userstyles.world/modules/database"
	"userstyles.world/modules/log"
	"userstyles.world/modules/search"
)

func DeleteGet(c *fiber.Ctx) error {
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
			"Title": "Users don't match",
			"User":  u,
		})
	}

	return c.Render("style/delete", fiber.Map{
		"Title": "Confirm deletion",
		"User":  u,
		"Style": s,
	})
}

func DeletePost(c *fiber.Ctx) error {
	u, _ := jwt.User(c)

	i, err := c.ParamsInt("id")
	if err != nil || i < 1 {
		return c.Status(fiber.StatusBadRequest).Render("err", fiber.Map{
			"User":  u,
			"Title": "Invalid style ID",
		})
	}
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

	// Delete style from database.
	q := new(models.Style)
	if err = database.Conn.Delete(q, "styles.id = ?", id).Error; err != nil {
		log.Warn.Printf("Failed to delete style %d: %s\n", s.ID, err.Error())
		return c.Render("err", fiber.Map{
			"Title": "Internal server error",
			"User":  u,
		})
	}

	// Delete stats from database.
	if err = new(models.Stats).Delete(s.ID); err != nil {
		log.Warn.Printf("Failed to delete stats for style %d: %s\n", s.ID, err.Error())
		c.Status(fiber.StatusInternalServerError)
		return c.Render("err", fiber.Map{
			"Title": "Internal server error",
			"User":  u,
		})
	}

	if err = models.RemoveStyleCode(strconv.Itoa(int(s.ID))); err != nil {
		log.Warn.Printf("kind=removecode id=%v err=%q\n", s.ID, err)
	}

	if err = search.DeleteStyle(s.ID); err != nil {
		log.Warn.Printf("Failed to delete style %d from index: %s", s.ID, err)
	}

	cache.Code.Remove(i)

	return c.Redirect("/user/"+u.Username, fiber.StatusSeeOther)
}
