// Package images implements functions for parsing and converting images.
package images

import (
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path"
	"regexp"

	"userstyles.world/modules/config"
	"userstyles.world/modules/log"
)

// Generate will generate all needed images used for styles.
func Generate(fh *multipart.FileHeader, id, version, oldURL, newURL string) error {
	switch {
	case fh != nil:
		return GenerateFromUpload(fh, id, version)
	case oldURL != newURL:
		return GenerateFromURL(id, version, newURL)
	default:
		return nil
	}
}

// GenerateFromUpload will generate images from an uploaded file.
func GenerateFromUpload(fh *multipart.FileHeader, id, version string) error {
	img, err := fh.Open()
	if err != nil {
		return err
	}

	data, _ := io.ReadAll(img)
	return processImages(id, version, "file upload", data)
}

// GenerateFromURL will generate images from an external URL.
func GenerateFromURL(id, version, url string) error {
	url = fixRawURL(url)
	req, err := http.Get(url)
	if err != nil {
		log.Warn.Printf("Failed to fetch %q for %s: %s\n", url, id, err)
		return err
	}
	defer req.Body.Close()

	data, _ := io.ReadAll(req.Body)
	return processImages(id, version, url, data)
}

func fixRawURL(url string) string {
	re := regexp.MustCompile(`(?mi)^(http.*)/(raw|src|blob)/(.*.(png|jpe?g|avif|webp))(\?.*)*$`)
	return re.ReplaceAllString(url, "${1}/raw/${3}")
}

func processImages(id, version, from string, data []byte) error {
	original := path.Join(config.ImageDir, id+".original")
	if err := os.WriteFile(original, data, 0o600); err != nil {
		return fmt.Errorf("failed to save image for %s from %q: %s",
			id, from, err)
	}

	dir := path.Join(config.PublicDir, id)
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		if err := os.MkdirAll(dir, 0o755); err != nil {
			log.Warn.Fatalf("Failed to create %q: %s\n", dir, err)
		}
	}

	webp := path.Join(dir, version+".webp")
	if err := decodeImage(original, webp, imageFullWebP); err != nil {
		return fmt.Errorf("Failed to decode WebP image for %s from %q: %s",
			id, from, err)
	}
	if err := decodeImage(original, webp, imageThumbWebP); err != nil {
		return fmt.Errorf("Failed to decode WebP thumb for %s from %q: %s",
			id, from, err)
	}

	jpeg := path.Join(dir, version+".jpeg")
	if err := decodeImage(original, jpeg, imageFullJPEG); err != nil {
		return fmt.Errorf("Failed to decode JPEG image for %s from %q: %s",
			id, from, err)
	}
	if err := decodeImage(original, jpeg, imageThumbJPEG); err != nil {
		return fmt.Errorf("Failed to decode JPEG thumb for %s from %q: %s",
			id, from, err)
	}

	return nil
}
