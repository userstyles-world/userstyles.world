package utils

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"strings"

	"github.com/vednoc/go-usercss-parser"

	"userstyles.world/models"
	"userstyles.world/modules/errors"
)

const (
	ArchiveURL = "https://raw.githubusercontent.com/33kk/uso-archive/flomaster/data/"
	DataURL    = ArchiveURL + "styles/"
	StyleURL   = ArchiveURL + "usercss/"
	PreviewURL = ArchiveURL + "screenshots/"
)

// Using only the data that we need.
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
	s := &models.Style{}
	if err != nil {
		log.Printf("failed to extract id, err: %v\n", err)
		return s, errors.ErrFailedProcessData
	}

	data, err := fetchJSON(id)
	if err != nil {
		log.Printf("failed to fetch json, err: %v\n", err)
		return s, errors.ErrFailedFetch
	}

	res, err := unmarshalJSON(data)
	if err != nil {
		log.Printf("failed to unmarshal json, err: %v\n", err)
		return s, errors.ErrFailedProcessData
	}

	// Fetch generated UserCSS format.
	source := StyleURL + id + ".user.css"
	uc, err := usercss.ParseFromURL(source)
	if err != nil {
		log.Printf("failed to parse style from URL, err: %v\n", err)
		return s, errors.ErrFailedFetch
	}

	s = &models.Style{
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
		log.Println("Error fetching URL:", err)
		return nil, err
	}
	defer req.Body.Close()

	body, err := ioutil.ReadAll(req.Body)
	if err != nil {
		log.Println("Error reading body:", err)
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
		log.Printf("failed to unmarshal json, err: %v\n", err)
		return data, err
	}

	return data, nil
}
