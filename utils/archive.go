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

const ArchiveURL = "https://raw.githubusercontent.com/33kk/uso-archive/flomaster/data/"

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

func ImportFromArchive(url string, u models.APIUser) *models.Style {
	// TODO: Implement better error handling.
	id, err := extractID(url)
	if err != nil {
		log.Printf("failed to extract id, err: %v\n", err)
		panic(err)
	}
	log.Printf("extract id: %v\n", id)

	data, err := fetchJSON(id)
	if err != nil {
		log.Printf("failed to fetch json, err: %v\n", err)
		panic(err)
	}

	res, err := unmarshalJSON(data)
	if err != nil {
		log.Printf("failed to unmarshal json, err: %v\n", err)
		panic(err)
	}

	log.Printf("final ->> %#+v\n", res)

	// Fetch generated UserCSS format.
	source := ArchiveURL + "usercss/" + id + ".user.css"
	uc, err := usercss.ParseFromURL(source)
	if err != nil {
		log.Printf("failed to parse style from URL, err: %v\n", err)
		panic(err)
	}

	s := &models.Style{
		UserID:      u.ID,
		Name:        uc.Name,
		Description: uc.Description,
		Notes:       res.Info.AdditionalInfo,
		Code:        uc.SourceCode,
		License:     uc.License,
		Preview:     ArchiveURL + "screenshots/" + res.Screenshots.Main.Name,
		Homepage:    uc.HomepageURL,
		Category:    res.Info.Category,
		Original:    url,
	}

	return s
}

func extractID(url string) (string, error) {
	if !strings.HasPrefix(url, ArchiveURL) {
		return "", errors.New("style isn't from uso-archive")
	}

	// Trim everything except style id.
	url = strings.TrimPrefix(url, ArchiveURL+"usercss/")
	url = strings.TrimSuffix(url, ".user.css")

	return url, nil
}

func fetchJSON(id string) ([]byte, error) {
	url := ArchiveURL + "styles/" + id + ".json"

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
