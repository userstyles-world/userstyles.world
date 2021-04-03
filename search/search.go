package search

import "github.com/blevesearch/bleve/v2"

func Search(text string) ([]string, error) {
	query := bleve.NewQueryStringQuery(text)
	sr, err := StyleIndex.Search(bleve.NewSearchRequestOptions(query, 5, 0, false))
	if err != nil {
		return nil, err
	}
}
