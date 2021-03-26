package utils

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"log"
	"net/http"
	"strings"

	"github.com/vednoc/go-usercss-parser"

	"userstyles.world/models"
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
	// TODO: Implement better error handling.
	id, err := extractID(url)
	s := &models.Style{}
	if err != nil {
		log.Printf("failed to extract id, err: %v\n", err)
		return s, errors.New("Failed to extract Style ID")
	}

	data, err := fetchJSON(id)
	if err != nil {
		log.Printf("failed to fetch json, err: %v\n", err)
		return s, errors.New("Failed to fetch Style data")
	}

	res, err := unmarshalJSON(data)
	if err != nil {
		log.Printf("failed to unmarshal json, err: %v\n", err)
		return s, errors.New("Failed to process Style data")
	}

	// Fetch generated UserCSS format.
	source := StyleURL + id + ".user.css"
	uc, err := usercss.ParseFromURL(source)
	if err != nil {
		log.Printf("failed to parse style from URL, err: %v\n", err)
		return s, errors.New("Ffailed to fetch Style")
	}

	s = &models.Style{
		UserID:      u.ID,
		Name:        uc.Name,
		Description: uc.Description,
		Notes:       res.Info.AdditionalInfo,
		Code:        uc.SourceCode,
		License:     uc.License,
		Preview:     PreviewURL + res.Screenshots.Main.Name,
		Homepage:    uc.HomepageURL,
		Category:    res.Info.Category,
		Original:    url,
	}

	return s, nil
}

func extractID(url string) (string, error) {
	if !strings.HasPrefix(url, ArchiveURL) {
		return "", errors.New("style isn't from uso-archive")
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
		return nil, errors.New("style not found")
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
