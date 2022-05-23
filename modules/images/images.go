// Package images implements functions for parsing and converting images.
package images

import (
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path"

	"userstyles.world/modules/config"
	"userstyles.world/modules/log"
)

// Generate will generate all needed images used for styles.
func Generate(fh *multipart.FileHeader, id, old, new string) error {
	if fh != nil {
		GenerateFromUpload(fh, id)
	} else if fh == nil && old != new {
		return GenerateFromURL(id, new)
	}

	return nil
}

// GenerateFromUpload will generate images from an uploaded file.
func GenerateFromUpload(fh *multipart.FileHeader, id string) error {
	img, err := fh.Open()
	if err != nil {
		return err
	}

	data, _ := io.ReadAll(img)
	return processImages(id, "file upload", data)
}

// GenerateFromURL will generate images from an external URL.
func GenerateFromURL(id, url string) error {
	url = fixRawURL(url)
	req, err := http.Get(url)
	if err != nil {
		log.Warn.Printf("Failed to fetch %q for %s: %s\n", url, id, err)
		return err
	}
	defer req.Body.Close()

	data, _ := io.ReadAll(req.Body)
	return processImages(id, url, data)
}

func processImages(id, from string, data []byte) error {
	template := path.Join(config.ImageDir, id)
	original := template + ".original"
	jpeg := template + ".jpeg"
	webp := template + ".webp"

	if err := os.WriteFile(original, data, 0o600); err != nil {
		return fmt.Errorf("failed to save image for %s from %q: %s",
			id, from, err)
	}

	if err := decodeImage(original, jpeg, ImageTypeJPEG); err != nil {
		return fmt.Errorf("failed to generate JPEG image for %s from %q: %s",
			id, from, err)
	}

	if err := decodeImage(original, webp, ImageTypeWEBP); err != nil {
		return fmt.Errorf("failed to generate WebP image for %s from %q: %s",
			id, from, err)
	}

	return nil
}
