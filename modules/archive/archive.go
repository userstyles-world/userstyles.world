package archive

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"regexp"
	"strings"

	"github.com/vednoc/go-usercss-parser"

	"userstyles.world/models"
	"userstyles.world/modules/errors"
	"userstyles.world/modules/log"
)

const (
	ArchiveURL = orgURL
	DataURL    = ArchiveURL + "styles/"
	StyleURL   = ArchiveURL + "usercss/"
	PreviewURL = ArchiveURL + "screenshots/"

	cdnURL = "https://cdn.jsdelivr.net/gh/33kk/uso-archive@flomaster/data/"
	hubURL = "https://raw.githubusercontent.com/33kk/uso-archive/flomaster/data/"
	orgURL = "https://raw.githubusercontent.com/uso-archive/data/flomaster/data/"
	oldURL = "https://uso-archive.surge.sh/"
	newURL = "https://uso.kkx.one/style/"
)

var (
	// ErrFailedToExtractID returns an error if no ID has been extracted.
	ErrFailedToExtractID = fmt.Errorf("failed to extract ID")

	// idRe holds a regexp for extracting style IDs.
	idRe = regexp.MustCompile(`.*?/(\?(?:page=\d+\&)style[=/])?(\d+)(\.user\.css)?$`)
)

// IsFromArchive checks whether a userstyle comes from a USo-archive.
func IsFromArchive(url string) bool {
	for _, prefix := range [...]string{cdnURL, hubURL, orgURL, newURL, oldURL} {
		if strings.HasPrefix(url, prefix) {
			return true
		}
	}

	return false
}

// RewriteURL consolidates disparate URLs that point to USo-archive.
func RewriteURL(url string) (string, error) {
	if !strings.HasPrefix(url, orgURL) && IsFromArchive(url) {
		id, err := extractID(url)
		if err != nil {
			return "", err
		}

		return fmt.Sprintf("%s%s.user.css", StyleURL, id), nil
	}

	return url, nil
}

// data holds userstyle fields that we need.
type data struct {
	Info struct {
		Description    string
		AdditionalInfo string
		Category       string
	}
	Screenshots struct {
		Main struct {
			Name string
		}
	}
}

func ImportFromArchive(url string, u models.APIUser) (*models.Style, error) {
	id, err := extractID(url)
	if err != nil {
		log.Info.Printf("Failed to extract ID from %q: %s\n", url, err)
		return nil, ErrFailedToExtractID
	}

	res, err := fetchJSON(id)
	if err != nil {
		log.Info.Printf("Failed to fetch data for %s: %s\n", id, err)
		return nil, errors.ErrFailedFetch
	}

	data, err := unmarshalJSON(res)
	if err != nil {
		log.Info.Printf("Failed to unmarshal style %s: %s\n", id, err)
		return nil, errors.ErrFailedProcessData
	}

	// Fetch generated UserCSS format.
	uc := new(usercss.UserCSS)
	source := StyleURL + id + ".user.css"
	if err = uc.ParseURL(source); err != nil {
		log.Info.Printf("Failed to parse style for %s: %s\n", source, err)
		return nil, errors.ErrFailedFetch
	}

	s := &models.Style{
		UserID:      u.ID,
		Name:        uc.Name,
		Description: data.Info.Description,
		Notes:       data.Info.AdditionalInfo,
		Code:        uc.SourceCode,
		License:     uc.License,
		Preview:     PreviewURL + data.Screenshots.Main.Name,
		Homepage:    uc.HomepageURL,
		Category:    data.Info.Category,
		Original:    url,
	}

	// Disallow GIF format.
	if strings.HasSuffix(s.Preview, ".gif") {
		log.Info.Printf("Removed GIF image for style %s\n", id)
		s.Preview = ""
	}

	return s, nil
}

// extractID tries to extract ID from a provided URL.
func extractID(url string) (string, error) {
	s := idRe.ReplaceAllString(url, "$2")
	if s == url {
		return "", ErrFailedToExtractID
	}

	return s, nil
}

func fetchJSON(id string) ([]byte, error) {
	req, err := http.Get(DataURL + id + ".json")
	if err != nil {
		return nil, err
	}
	defer req.Body.Close()

	body, err := io.ReadAll(req.Body)
	if err != nil {
		return nil, err
	}

	// Return error if style doesn't exist.
	if string(body) == "404: Not Found" {
		return nil, errors.ErrStyleNotFound
	}

	return body, nil
}

func unmarshalJSON(raw []byte) (*data, error) {
	var data data
	if err := json.Unmarshal(raw, &data); err != nil {
		return nil, err
	}

	return &data, nil
}
