package api

import (
	"io/fs"
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
	c.Status(fiber.StatusNotFound)
	return c.SendString("Screenshot not found")
}

func GetPreviewScreenshot(c *fiber.Ctx) error {
	styleID := c.Params("id")
	format := getFileExtension(styleID)
	styleID = strings.TrimSuffix(styleID, format)

	info, err := images.GetImageFromStyle(styleID)
	if err != nil {
		return notFound(c)
	}

	var stat fs.FileInfo
	var fileName string
	var mimeType string
	orignalFile := images.CacheFolder + styleID + ".original"

	switch format[1:] {
	case "jpeg":
		if info.Jpeg == nil {
			return notFound(c)
		}
		stat = info.Jpeg
		fileName = images.CacheFolder + styleID + ".jpeg"
		mimeType = "image/jpeg"
	case "webp":
		fileName = images.CacheFolder + styleID + ".webp"
		if info.WebP == nil {
			err = images.DecodeImage(orignalFile, fileName, images.ImageTypeWEBP)
			if err != nil {
				return notFound(c)
			}
			webpStat, err := os.Stat(fileName)
			if err != nil {
				return notFound(c)
			}
			stat = webpStat
			break
		}
		stat = info.WebP
		mimeType = "image/webp"
	}

	if stat == nil || fileName == "" {
		return notFound(c)
	}

	file, err := os.Open(fileName)
	if err != nil {
		return notFound(c)
	}

	// Set caching to two weeks.
	c.Response().Header.Set(fiber.HeaderCacheControl, "public, max-age=1209600")

	c.Response().Header.SetContentType(mimeType)
	c.Response().SetBodyStream(file, int(stat.Size()))

	return nil
}
