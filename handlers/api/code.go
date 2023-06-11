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
	"userstyles.world/modules/log"
)

func GetStyleCode(c *fiber.Ctx) error {
	kind := c.Params("ext")
	if kind != "css" && kind != "styl" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "invalid userstyle extension",
		})
	}

	id := c.Params("id")
	if i, err := strconv.Atoi(id); err != nil || i < 1 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "invalid userstyle ID",
		})
	}

	code, err := os.ReadFile(filepath.Join(config.StyleDir, id))
	if err != nil {
		log.Info.Println(err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "userstyle not found",
		})
	}

	if c.Method() == fiber.MethodGet {
		c.Type("css", "utf-8") // #107
	}

	cl := strconv.Itoa((len(code)))
	cs := crc32.ChecksumIEEE([]byte(code))
	c.Set("ETag", fmt.Sprintf("%s-%d", cl, cs))

	cache.InstallStats.Add(c.IP() + " " + id)

	return c.Send(code)
}
