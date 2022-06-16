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
	return processImages(id, version, fh.Filename, data)
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

var rawURLExtractor = regexp.MustCompile(`(?mi)^(http.*)/(raw|src|blob)/(.*.(png|jpe?g|avif|webp))(\?.*)*$`)

func fixRawURL(url string) string {
	return rawURLExtractor.ReplaceAllString(url, "${1}/raw/${3}")
}

func processImages(id, version, url string, data []byte) error {
	original := path.Join(config.ImageDir, id+".original")
	if err := os.WriteFile(original, data, 0o600); err != nil {
		return fmt.Errorf("failed to save image for %s from %q: %s",
			id, url, err)
	}

	dir := path.Join(config.PublicDir, id)
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		if err := os.MkdirAll(dir, 0o755); err != nil {
			log.Warn.Printf("Failed to create %q for %s: %s\n", dir, id, err)
			return err
		}
	}

	webp := path.Join(dir, version+".webp")
	jpeg := path.Join(dir, version+".jpeg")
	images := []struct {
		out  string
		kind imageKind
	}{
		{out: webp, kind: imageFullWebP},
		{out: webp, kind: imageThumbWebP},
		{out: jpeg, kind: imageFullJPEG},
		{out: jpeg, kind: imageThumbJPEG},
	}

	for _, image := range images {
		if err := decodeImage(original, image.out, image.kind); err != nil {
			return fmt.Errorf("failed to decode %s for style %s from %q: %s",
				image.kind, id, url, err)
		}
	}

	return nil
}
