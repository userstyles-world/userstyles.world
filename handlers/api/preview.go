package api

import (
	"os"
	"strings"

	"github.com/gofiber/fiber/v2"

	"userstyles.world/modules/images"
)

func getFileExtension(path string) string {
	n := strings.LastIndexByte(path, '.')
	if n < 0 {
		return ""
	}
	return path[n:]
}

var notFound = func(c *fiber.Ctx) error {
	return c.Status(fiber.StatusNotFound).
		JSON(fiber.Map{"data": "Error: screenshot not found"})
}

func GetPreviewScreenshot(c *fiber.Ctx) error {
	styleID := c.Params("id")
	format := getFileExtension(styleID)
	styleID = strings.TrimSuffix(styleID, format)

	var fileName string

	// Redirect to prevent panic.
	if format == "" {
		return c.Redirect("/api/style/preview/"+styleID+".jpeg", fiber.StatusSeeOther)
	}

	// Only allow jpeg and webp as formats.
	switch format[1:] {
	case "jpeg":
		fileName = images.CacheFolder + styleID + ".jpeg"
	case "webp":
		fileName = images.CacheFolder + styleID + ".webp"
	default:
		return notFound(c)
	}

	var (
		file *os.File
		stat os.FileInfo
	)

	file, err := os.Open(fileName)
	if err != nil {
		return notFound(c)
	}

	if stat, err = file.Stat(); err != nil {
		return notFound(c)
	}

	contentLength := int(stat.Size())

	// Set Content Type header
	c.Type(getFileExtension(stat.Name()))

	// Set caching to 3 month.
	// Images are very likely not changing that often.
	// 60 * 60 * 24 * 31 * 3
	c.Response().Header.Set(fiber.HeaderCacheControl, "public, max-age=8035200")

	c.Response().SetBodyStream(file, contentLength)

	return nil
}
