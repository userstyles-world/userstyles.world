package search

import (
	"github.com/blevesearch/bleve/v2"
)

func SearchText(text string) ([]string, error) {
	query := bleve.NewQueryStringQuery(text)
	searchRequest := bleve.NewSearchRequestOptions(query, 5, 0, false)
	// searchRequest.Fields = []string{"*"}

	sr, err := StyleIndex.Search(searchRequest)
	if err != nil {
		return nil, err
	}
	var returnResult []string
	for _, a := range sr.Hits {
		if err != nil {
			return nil, err
		}
		returnResult = append(returnResult, a.ID)
	}
	return returnResult, nil
}
