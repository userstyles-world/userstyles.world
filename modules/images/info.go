package images

import (
	"io"
	"io/fs"
	"net/http"
	"os"
	"strings"

	"userstyles.world/modules/log"
)

var CacheFolder = "./data/images/"

func Initialize() {
	if fileInfo := fileExist(CacheFolder); fileInfo == nil {
		if err := os.Mkdir(CacheFolder, 0o755); err != nil {
			log.Warn.Fatalln(err)
		}
	}
}

func fileExist(path string) fs.FileInfo {
	stat, err := os.Stat(path)
	if os.IsNotExist(err) {
		return nil
	}
	return stat
}

func fixGitHubURL(url string) string {
	if strings.Contains(url, "/blob/") {
		parts := strings.Split(url, "/blob/")
		url = parts[0] + "/raw/" + parts[1]
	}

	return url
}

func GenerateImagesForStyle(id, preview string, isOriginalLocal bool) error {
	template := CacheFolder + id
	original := template + ".original"
	jpeg := template + ".jpeg"
	webp := template + ".webp"

	// Is the preview image not a local file?
	// Let's download it.
	if !isOriginalLocal {
		if strings.HasPrefix(preview, "https://github.com/") {
			preview = fixGitHubURL(preview)
		}

		req, err := http.Get(preview)
		if err != nil {
			log.Warn.Println("Error fetching image URL:", err)
			return err
		}
		defer req.Body.Close()

		// Make sure to received the full body before processing it.
		data, _ := io.ReadAll(req.Body)
		err = os.WriteFile(original, data, 0o600)
		if err != nil {
			log.Warn.Println("Error processing image:", err)
			return err
		}
	}

	err := decodeImage(original, jpeg, ImageTypeJPEG)
	if err != nil {
		log.Warn.Println("Error processing image:", err)
		return err
	}

	err = decodeImage(original, webp, ImageTypeWEBP)
	if err != nil {
		log.Warn.Println("Error processing image:", err)
		return err
	}

	return nil
}
