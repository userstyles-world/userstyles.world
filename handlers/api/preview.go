package api

import (
	"io/fs"
	"os"
	"strings"

	"github.com/gofiber/fiber/v2"
	"userstyles.world/images"
)

func getFileExtension(path string) string {
	n := strings.LastIndexByte(path, '.')
	if n < 0 {
		return ""
	}
	return path[n:]
}

func GetPreviewScreenshot(c *fiber.Ctx) error {
	styleId := c.Params("id")
	format := getFileExtension(styleId)
	styleId = strings.TrimSuffix(styleId, format)

	notFound := func(c *fiber.Ctx) error {
		c.Status(fiber.StatusNotFound)
		return c.SendString("Screenshot not found")
	}
	info, err := images.GetImageFromStyle(styleId)
	if err != nil {
		return notFound(c)
	}

	var stat fs.FileInfo
	var fileName string
	orignialFile := images.CacheFolder + styleId + ".originial"

	switch format[1:] {
	case "jpeg":
		if info.Jpeg == nil {
			return notFound(c)
		}
		stat = info.Originial
		fileName = images.CacheFolder + styleId + ".jpeg"
	case "webp":
		webpName := images.CacheFolder + styleId + ".webp"
		if info.WebP == nil {
			file, err := os.Open(orignialFile)
			if err != nil {
				return notFound(c)
			}
			if err = images.ProcessToWebp(file, webpName); err != nil {
				return notFound(c)
			}
			webpStat, err := os.Stat(webpName)
			if err != nil {
				return notFound(c)
			}
			fileName = webpName
			stat = webpStat
			break
		}
		stat = info.WebP
		fileName = webpName
	case "avif":
		avifName := images.CacheFolder + styleId + ".avif"
		if info.Avif == nil {
			file, err := os.Open(orignialFile)
			if err != nil {
				return notFound(c)
			}
			if err = images.ProcessToAvif(file, avifName); err != nil {
				return notFound(c)
			}
			avifStat, err := os.Stat(avifName)
			if err != nil {
				return notFound(c)
			}
			fileName = avifName
			stat = avifStat
			break
		}
		stat = info.Avif
		fileName = avifName
	}

	if stat == nil || fileName == "" {
		return notFound(c)
	}

	file, err := os.Open(fileName)
	if err != nil {
		return notFound(c)
	}
	contentLength := int(stat.Size())

	// fasthttp currently doesn't recognize avif format.
	// Nor does nginx?
	if format == ".avif" {
		c.Context().SetContentType("image/avif")
	} else {
		c.Type(getFileExtension(stat.Name()))
	}

	c.Response().SetBodyStream(file, contentLength)
	return nil
}
