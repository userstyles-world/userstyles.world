package utils

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"log"
	"net/http"
	"strings"

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

	res, err := fetchJSON(id)
	if err != nil {
		log.Printf("failed to fetch json, err: %v\n", err)
		panic(err)
	}
	// log.Printf("json data: %v\n", string(res))

	final, err := unmarshalJSON(res)
	if err != nil {
		log.Printf("failed to unmarshal json, err: %v\n", err)
		panic(err)
	}

	log.Printf("final ->> %#+v\n", final)

	s := new(models.Style)
	s.UserID = u.ID
	s.Name = "test"
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
