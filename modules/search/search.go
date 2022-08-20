// Package search provides helper functions for Bleve.
package search

import (
	"errors"
	"time"

	"github.com/blevesearch/bleve/v2"

	"userstyles.world/modules/log"
	"userstyles.world/modules/storage"
)

var (
	// ErrSearchNoResults errors that search engine couldn't find results.
	ErrSearchNoResults = errors.New("no search results found")

	// ErrSearchBadRequest errors that search engine had an internal error.
	ErrSearchBadRequest = errors.New("bad search request")

	// StyleIndex holds the connection to search engine.
	StyleIndex bleve.Index
)

// engineMetrics returns basic metrics for search queries.
type engineMetrics struct {
	Hits      int
	TimeSpent time.Duration
}

// FindStylesByText searches for text and returns styles from search index and
// performance metrics, or an error if it fails to find anything.
func FindStylesByText(text, kind string, size int) ([]storage.StyleCard, engineMetrics, error) {
	metrics := engineMetrics{}
	// See https://github.com/blevesearch/bleve/issues/1290
	// FuzzySearch won't work the way I'd like the search to behave.
	// This way it will be more "loslly" and actually uses the tokenizers.
	// That we provide within the mappings.go and provide better results.
	timeStart := time.Now()
	sanitzedQuery := bleve.NewMatchQuery(text)

	searchRequest := bleve.NewSearchRequestOptions(sanitzedQuery, size, 0, true)
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

	res, err := storage.FindStyleCardsForSearch(nums(), kind, 96)
	if err != nil {
		metrics.Hits = 0
		return res, metrics, err
	}
	metrics.TimeSpent = time.Since(timeStart)

	return res, metrics, nil
}
