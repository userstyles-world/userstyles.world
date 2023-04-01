package style

import (
	"fmt"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"

	"userstyles.world/handlers/jwt"
	"userstyles.world/modules/cache"
	"userstyles.world/modules/mirror"
	"userstyles.world/modules/storage"
)

func Mirror(c *fiber.Ctx) error {
	u, _ := jwt.User(c)
	m := fiber.Map{"User": u}

	key := "mirror-" + c.Params("id")
	_, t, found := cache.Store.GetWithExpiration(key)
	if found {
		// TODO: Explore adding future time to `RelTime`.
		dur := time.Since(t).Truncate(time.Second).String()
		m["Title"] = "You can use this again in " + strings.TrimPrefix(dur, "-")
		return c.Render("err", m)
	}

	id, err := c.ParamsInt("id")
	if err != nil {
		m["Title"] = "Invalid userstyle ID"
		return c.Render("err", m)
	}

	s, err := storage.FindStyleForMirror(id)
	if err != nil {
		m["Title"] = "Mirrored style not found"
		return c.Render("err", m)
	}

	if u.ID != s.UserID {
		m["Title"] = "User and style author don't match"
		return c.Render("err", m)
	}

	mirror.MirrorStyle(s)
	cache.Store.Add(key, "", time.Hour)

	return c.Redirect(fmt.Sprintf("/style/%d", s.ID))
}
