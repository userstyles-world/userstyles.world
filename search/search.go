package search

import (
	"fmt"
	"time"

	"github.com/blevesearch/bleve/v2"

	"userstyles.world/utils/strings"
)

const (
	timeFormat = "2006-01-02T15:04:05Z"
)

type MinimalStyle struct {
	ID          int       `json:"id"`
	UpdatedAt   time.Time `json:"updated_at"`
	Username    string    `json:"username"`
	DisplayName string    `json:"display_name"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Preview     string    `json:"preview"`
	Notes       string    `json:"notes"`
	Views       int       `json:"views"`
	Installs    int       `json:"installs"`
}

func (s MinimalStyle) Slug() string {
	return strings.SlugifyURL(s.Name)
}

func (s MinimalStyle) StyleURL() string {
	return fmt.Sprintf("/style/%d/%s", s.ID, s.Slug())
}

func (s MinimalStyle) Author() string {
	if s.DisplayName != "" {
		return s.DisplayName
	}

	return s.Username
}

func FindStylesByText(text string) ([]MinimalStyle, error) {
	// See https://github.com/blevesearch/bleve/issues/1290
	// FuzzySearch won't work the way I'd like the search to behave.
	// This way it will be more "loslly" and actually uses the tokenizers.
	// That we provide within the mappings.go and provide better results.
	sanitzedQuery := bleve.NewMatchQuery(text)

	searchRequest := bleve.NewSearchRequestOptions(sanitzedQuery, 99, 0, false)
	searchRequest.Fields = []string{"*"}

	sr, err := StyleIndex.Search(searchRequest)
	if err != nil {
		return nil, err
	}

	returnResult := make([]MinimalStyle, 0, len(sr.Hits))
	for _, hit := range sr.Hits {
		if err != nil {
			return nil, err
		}

		updated, err := time.Parse(timeFormat, hit.Fields["updated_at"].(string))
		if err != nil {
			return nil, err
		}

		styleInfo := MinimalStyle{
			UpdatedAt:   updated,
			ID:          int(hit.Fields["id"].(float64)),
			Username:    hit.Fields["username"].(string),
			DisplayName: hit.Fields["display_name"].(string),
			Name:        hit.Fields["name"].(string),
			Description: hit.Fields["description"].(string),
			Preview:     hit.Fields["preview"].(string),
			Notes:       hit.Fields["notes"].(string),
			Views:       int(hit.Fields["views"].(float64)),
			Installs:    int(hit.Fields["installs"].(float64)),
		}

		returnResult = append(returnResult, styleInfo)
	}
	return returnResult, nil
}
