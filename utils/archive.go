package utils

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/vednoc/go-usercss-parser"

	"userstyles.world/models"
	"userstyles.world/modules/errors"
	"userstyles.world/modules/log"
)

const (
	ArchiveURL = "https://cdn.jsdelivr.net/gh/33kk/uso-archive@flomaster/data/"
	DataURL    = ArchiveURL + "styles/"
	StyleURL   = ArchiveURL + "usercss/"
	PreviewURL = ArchiveURL + "screenshots/"
)

// Data struct contains only the data that we need.
type Data struct {
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
		log.Info.Println("Failed to extract style id:", err.Error())
		return nil, errors.ErrFailedProcessData
	}

	data, err := fetchJSON(id)
	if err != nil {
		log.Info.Println("Failed to fetch style JSON:", err.Error())
		return nil, errors.ErrFailedFetch
	}

	res, err := unmarshalJSON(data)
	if err != nil {
		log.Info.Println("Failed to unmarshal style JSON:", err.Error())
		return nil, errors.ErrFailedProcessData
	}

	// Fetch generated UserCSS format.
	uc := new(usercss.UserCSS)
	source := StyleURL + id + ".user.css"
	if err = uc.ParseURL(source); err != nil {
		log.Info.Printf("Failed to parse style from URL %v: %v\n", source, err)
		return nil, errors.ErrFailedFetch
	}

	s := &models.Style{
		UserID:      u.ID,
		Name:        uc.Name,
		Description: res.Info.Description,
		Notes:       res.Info.AdditionalInfo,
		Code:        uc.SourceCode,
		License:     uc.License,
		Preview:     PreviewURL + res.Screenshots.Main.Name,
		Homepage:    uc.HomepageURL,
		Category:    res.Info.Category,
		Original:    url,
	}

	// Disallow GIF format.
	if strings.HasSuffix(s.Preview, ".gif") {
		s.Preview = ""
	}

	return s, nil
}

func extractID(url string) (string, error) {
	if !strings.HasPrefix(url, ArchiveURL) {
		return "", errors.ErrStyleNotFromUSO
	}

	// Trim everything except style id.
	url = strings.TrimPrefix(url, StyleURL)
	url = strings.TrimSuffix(url, ".user.css")

	return url, nil
}

func fetchJSON(id string) ([]byte, error) {
	url := DataURL + id + ".json"

	req, err := http.Get(url)
	if err != nil {
		log.Warn.Println("Error fetching style URL:", err.Error())
		return nil, err
	}
	defer req.Body.Close()

	body, err := ioutil.ReadAll(req.Body)
	if err != nil {
		log.Info.Println("Error reading body:", err.Error())
		return nil, err
	}

	// Return error if style doesn't exist.
	if string(body) == "404: Not Found" {
		return nil, errors.ErrStyleNotFound
	}

	return body, nil
}

func unmarshalJSON(raw []byte) (Data, error) {
	data := Data{}
	err := json.Unmarshal(raw, &data)
	if err != nil {
		log.Info.Println("Failed to unmarshal style JSON:", err.Error())
		return data, err
	}

	return data, nil
}
