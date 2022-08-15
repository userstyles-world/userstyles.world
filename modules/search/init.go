package search

import (
	"errors"
	"os"
	"path"

	"github.com/blevesearch/bleve/v2"
	"github.com/blevesearch/bleve/v2/index/upsidedown"

	"userstyles.world/modules/config"
	"userstyles.world/modules/log"
)

func openBleveIndexFile(path string) (bleve.Index, error) {
	_, err := os.Stat(path)
	if err != nil {
		return nil, err
	}

	index, err := bleve.Open(path)
	if err != nil && errors.Is(err, upsidedown.IncompatibleVersion) {
		return nil, os.RemoveAll(path)
	} else if err != nil {
		return nil, err
	}
	return index, nil
}

// Initialize sets up search engine.
func Initialize() {
	// Limiting query size to 1 due to not having to process a lot of data.
	bleve.Config.SetAnalysisQueueSize(1)

	indexFile := path.Join(config.DataDir, "styles.bleve")
	stylesIndex, err := openBleveIndexFile(indexFile)
	if err != nil {
		log.Info.Println("Creating a new search index.")
		indexMapping := buildIndexMapping()
		stylesIndex, err = bleve.New(indexFile, indexMapping)
		if err != nil {
			log.Warn.Fatal(err)
		}
	}
	log.Info.Println("Opening search index.")

	StyleIndex = stylesIndex

	if config.SearchReindex {
		go indexStyles()
	}
}
