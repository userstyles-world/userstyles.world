package images

import (
	"fmt"
	"io"
	"io/fs"
	"log"
	"net/http"
	"os"

	"github.com/davidbyttow/govips/v2/vips"
	"userstyles.world/database"
	"userstyles.world/models"
)

type ImageInfo struct {
	Originial fs.FileInfo
	Jpeg      fs.FileInfo
	WebP      fs.FileInfo
}

var (
	CacheFolder = "./image_cache/"
)

func Initialize() {
	if fileInfo := fileExist(CacheFolder); fileInfo == nil {
		if err := os.Mkdir(CacheFolder, 0755); err != nil {
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
	originial := template + ".originial"
	jpeg := template + ".jpeg"
	webp := template + ".webp"
	if fileExist(originial) == nil {
		style, err := models.GetStyleByID(database.DB, id)
		if err != nil || style.Preview == "" {
			return ImageInfo{}, err
		}

		req, err := http.Get(style.Preview)
		if err != nil {
			fmt.Println("Error fetching URL:", err)
			return ImageInfo{}, err
		}
		defer req.Body.Close()

		// Make sure to received the full body before processing it.
		data, _ := io.ReadAll(req.Body)
		err = os.WriteFile(originial, data, 0644)
		if err != nil {
			fmt.Println("Error processing:", err)
			return ImageInfo{}, err
		}

		err = DecodeImage(originial, jpeg, vips.ImageTypeJPEG)
		if err != nil {
			fmt.Println("Error processing:", err)
			return ImageInfo{}, err
		}

		return ImageInfo{
			Originial: fileExist(originial),
		}, nil
	}
	return ImageInfo{
		Originial: fileExist(originial),
		Jpeg:      fileExist(jpeg),
		WebP:      fileExist(webp),
	}, nil
}
