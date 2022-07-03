package search

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/blevesearch/bleve/v2"

	"userstyles.world/modules/database"
	"userstyles.world/modules/log"
	"userstyles.world/utils/strutils"
)

const (
	timeFormat = "2006-01-02T15:04:05Z"
)

var (
	// ErrSearchNoResults errors that search engine couldn't find results.
	ErrSearchNoResults = errors.New("no search results found")

	// ErrSearchBadRequest errors that search engine had an internal error.
	ErrSearchBadRequest = errors.New("bad search request")
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
	metrics := PerformanceMetrics{}
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
		log.Warn.Printf("Failed to find results for %q: %s\n", text, err)
		return nil, metrics, ErrSearchBadRequest
	}

	hits := len(sr.Hits)
	if hits == 0 {
		return nil, metrics, ErrSearchNoResults
	}
	metrics.Hits = hits

	nums := func() []int {
		hits := make([]int, metrics.Hits)
		for i, hit := range sr.Hits {
			hits[i] = int(hit.Fields["id"].(float64))
		}
		return hits
	}

	var res []MinimalStyle
	installed := "(SELECT COUNT(*) FROM stats s WHERE s.style_id = styles.id AND s.install > 0) AS Installs"
	viewed := "(SELECT COUNT(*) FROM stats s WHERE s.style_id = styles.id AND s.view > 0) AS Views"
	author := "(SELECT username FROM users WHERE styles.user_id = users.id) AS Username"
	rating := "(SELECT ROUND(AVG(rating), 1) FROM reviews WHERE reviews.style_id = styles.id AND reviews.deleted_at IS NULL) AS Rating"
	fields := []string{"id", "created_at", "updated_at", "name", "preview", installed, viewed, author, rating}
	tx := database.Conn.Debug().Table("styles").Select(strings.Join(fields, ", "))
	// TODO: Fix up ordering of results.
	if err := tx.Find(&res, "id in ?", nums()).Error; err != nil {
		metrics.Hits = 0
		return res, metrics, err
	}
	metrics.TimeSpent = time.Since(timeStart)
	return res, metrics, nil
}
