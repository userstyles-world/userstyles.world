package style

import (
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"

	"userstyles.world/handlers/jwt"
	"userstyles.world/models"
	"userstyles.world/modules/cache"
	"userstyles.world/modules/database"
	"userstyles.world/modules/log"
	"userstyles.world/modules/storage"
)

func DeleteGet(c *fiber.Ctx) error {
	u, _ := jwt.User(c)

	i, err := c.ParamsInt("id")
	if err != nil || i < 1 {
		return c.Status(fiber.StatusBadRequest).Render("err", fiber.Map{
			"Title": "Invalid style ID",
			"User":  u,
		})
	}

	s, err := models.GetStyleByID(i)
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

	s, err := models.GetStyleByID(i)
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

	err = database.Conn.Transaction(func(tx *gorm.DB) error {
		if err = storage.DeleteUserstyle(tx, i); err != nil {
			return err
		}
		if err = models.DeleteStats(tx, i); err != nil {
			return err
		}
		if err = storage.DeleteSearchData(tx, i); err != nil {
			return err
		}
		return models.RemoveStyleCode(c.Params("id"))
	})
	if err != nil {
		log.Database.Printf("Failed to delete %d: %s\n", i, err)
		return c.Render("err", fiber.Map{
			"Title": "Failed to remove userstyle",
			"User":  u,
		})
	}

	cache.Code.Remove(i)

	return c.Redirect("/user/"+u.Username, fiber.StatusSeeOther)
}
