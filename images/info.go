package images

import (
	"io"
	"io/fs"
	"log"
	"net/http"
	"os"

	"github.com/davidbyttow/govips/v2/vips"

	"userstyles.world/models"
)

type ImageInfo struct {
	Original fs.FileInfo
	Jpeg     fs.FileInfo
	WebP     fs.FileInfo
}

var CacheFolder = "./data/images/"

func Initialize() {
	if fileInfo := fileExist(CacheFolder); fileInfo == nil {
		if err := os.Mkdir(CacheFolder, 0o755); err != nil {
			log.Fatalln(err)
		}
	}

	vips.Startup(&vips.Config{
		ConcurrencyLevel: 0,
		MaxCacheFiles:    0,
		MaxCacheMem:      0,
		MaxCacheSize:     0,
	})
}

func fileExist(path string) fs.FileInfo {
	stat, err := os.Stat(path)
	if os.IsNotExist(err) {
		return nil
	}
	return stat
}

func GetImageFromStyle(id string) (ImageInfo, error) {
	template := CacheFolder + id
	original := template + ".original"
	jpeg := template + ".jpeg"
	webp := template + ".webp"
	if fileExist(original) == nil {
		style, err := models.GetStyleByID(id)
		if err != nil || style.Preview == "" {
			return ImageInfo{}, err
		}

		req, err := http.Get(style.Preview)
		if err != nil {
			log.Println("Error fetching URL:", err)
			return ImageInfo{}, err
		}
		defer req.Body.Close()

		// Make sure to received the full body before processing it.
		data, _ := io.ReadAll(req.Body)
		err = os.WriteFile(original, data, 0o600)
		if err != nil {
			log.Println("Error processing:", err)
			return ImageInfo{}, err
		}

		err = DecodeImage(original, jpeg, vips.ImageTypeJPEG)
		if err != nil {
			log.Println("Error processing:", err)
			return ImageInfo{}, err
		}

		return ImageInfo{
			Original: fileExist(original),
		}, nil
	}

	return ImageInfo{
		Original: fileExist(original),
		Jpeg:     fileExist(jpeg),
		WebP:     fileExist(webp),
	}, nil
}
