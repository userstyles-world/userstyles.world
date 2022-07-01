package search

import (
	"errors"
	"fmt"
	"time"

	"github.com/blevesearch/bleve/v2"

	"userstyles.world/utils/strutils"
)

const (
	timeFormat = "2006-01-02T15:04:05Z"
)

var (
	// ErrSearchNoResults errors that it couldn't match anything in search index.
	ErrSearchNoResults = errors.New("no search results found")
)

type MinimalStyle struct {
	ID          int       `json:"id"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	Username    string    `json:"username"`
	DisplayName string    `json:"display_name"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Preview     string    `json:"preview"`
	Notes       string    `json:"notes"`
	Views       int64     `json:"views"`
	Installs    int64     `json:"installs"`
	Rating      float64   `json:"rating"`
}

type PerformanceMetrics struct {
	Hits      int
	TimeSpent time.Duration
}

func (s MinimalStyle) Slug() string {
	return strutils.SlugifyURL(s.Name)
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

func FindStylesByText(text string) ([]MinimalStyle, PerformanceMetrics, error) {
	performanceMetrics := PerformanceMetrics{}
	// See https://github.com/blevesearch/bleve/issues/1290
	// FuzzySearch won't work the way I'd like the search to behave.
	// This way it will be more "loslly" and actually uses the tokenizers.
	// That we provide within the mappings.go and provide better results.
	timeStart := time.Now()
	sanitzedQuery := bleve.NewMatchQuery(text)

	searchRequest := bleve.NewSearchRequestOptions(sanitzedQuery, 99, 0, false)
	searchRequest.Fields = []string{"*"}

	sr, err := StyleIndex.Search(searchRequest)
	if err != nil {
		return nil, PerformanceMetrics{}, err
	}

	if len(sr.Hits) == 0 {
		return nil, performanceMetrics, ErrSearchNoResults
	}

	returnResult := make([]MinimalStyle, 0, len(sr.Hits))

	performanceMetrics.Hits = len(sr.Hits)

	for _, hit := range sr.Hits {
		if err != nil {
			return nil, PerformanceMetrics{}, err
		}

		created, err := time.Parse(timeFormat, hit.Fields["created_at"].(string))
		if err != nil {
			return nil, PerformanceMetrics{}, err
		}

		updated, err := time.Parse(timeFormat, hit.Fields["updated_at"].(string))
		if err != nil {
			return nil, PerformanceMetrics{}, err
		}

		styleInfo := MinimalStyle{
			CreatedAt:   created,
			UpdatedAt:   updated,
			ID:          int(hit.Fields["id"].(float64)),
			Username:    hit.Fields["username"].(string),
			DisplayName: hit.Fields["display_name"].(string),
			Name:        hit.Fields["name"].(string),
			Description: hit.Fields["description"].(string),
			Preview:     hit.Fields["preview"].(string),
			Notes:       hit.Fields["notes"].(string),
			Views:       int64(hit.Fields["views"].(float64)),
			Installs:    int64(hit.Fields["installs"].(float64)),
			Rating:      hit.Fields["rating"].(float64),
		}

		returnResult = append(returnResult, styleInfo)
	}
	performanceMetrics.TimeSpent = time.Since(timeStart)
	return returnResult, performanceMetrics, nil
}
