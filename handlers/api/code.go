package api

import (
	"fmt"
	"hash/crc32"
	"os"
	"path/filepath"
	"strconv"

	"github.com/gofiber/fiber/v2"

	"userstyles.world/modules/cache"
	"userstyles.world/modules/config"
)

func GetStyleCode(c *fiber.Ctx) error {
	kind := c.Params("ext")
	if kind != "css" && kind != "styl" && kind != "less" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "invalid userstyle extension",
		})
	}

	id := c.Params("id")
	i, err := strconv.Atoi(id)
	if err != nil || i < 1 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "invalid userstyle ID",
		})
	}

	code := cache.Code.Get(i)
	if code == nil {
		code, err = os.ReadFile(filepath.Join(config.StyleDir, id))
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"message": "userstyle not found",
			})
		}
		cache.Code.Add(i, code)
	}

	if c.Method() == fiber.MethodGet {
		c.Type("css", "utf-8") // #107
	}

	cl := strconv.Itoa(len(code))
	cs := crc32.ChecksumIEEE(code)
	c.Set("ETag", fmt.Sprintf("%s-%d", cl, cs))

	cache.InstallStats.Add(c.IP() + " " + id)

	return c.Send(code)
}
