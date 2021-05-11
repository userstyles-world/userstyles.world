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

func FindStylesByText(text string) ([]MinimalStyle, error) {
	query := bleve.NewFuzzyQuery(text)
	searchRequest := bleve.NewSearchRequestOptions(query, 99, 0, false)
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
		styleInfo := MinimalStyle{}
		styleInfo.Name = hit.Fields["name"].(string)
		styleInfo.ID = hit.Fields["id"].(string)
		styleInfo.Description = hit.Fields["description"].(string)
		styleInfo.Preview = hit.Fields["preview"].(string)
		styleInfo.Notes = hit.Fields["notes"].(string)
		styleInfo.Username = hit.Fields["username"].(string)

		returnResult = append(returnResult, styleInfo)
	}
	return returnResult, nil
}
