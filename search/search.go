package search

import (
	"github.com/blevesearch/bleve/v2"
)

type MinimalStyle struct {
	ID          string `json:"id"`
	Username    string `json:"username"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Preview     string `json:"preview"`
	Notes       string `json:"notes"`
}

func SearchText(text string) ([]MinimalStyle, error) {
	query := bleve.NewQueryStringQuery(text)
	searchRequest := bleve.NewSearchRequestOptions(query, 99, 0, false)
	searchRequest.Fields = []string{"*"}

	sr, err := StyleIndex.Search(searchRequest)
	if err != nil {
		return nil, err
	}
	var returnResult []MinimalStyle
	for _, hit := range sr.Hits {
		if err != nil {
			return nil, err
		}
		StyleInfo := MinimalStyle{}
		StyleInfo.Name = hit.Fields["name"].(string)
		StyleInfo.ID = hit.Fields["id"].(string)
		StyleInfo.Description = hit.Fields["description"].(string)
		StyleInfo.Preview = hit.Fields["preview"].(string)
		StyleInfo.Notes = hit.Fields["notes"].(string)
		StyleInfo.Username = hit.Fields["username"].(string)

		returnResult = append(returnResult, StyleInfo)
	}
	return returnResult, nil
}
