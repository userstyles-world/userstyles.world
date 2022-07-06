package search

import (
	"errors"
	"os"
	"path"
	"strconv"
	"time"

	"github.com/blevesearch/bleve/v2"
	"github.com/blevesearch/bleve/v2/index/upsidedown"

	"userstyles.world/modules/config"
	"userstyles.world/modules/log"
	"userstyles.world/modules/storage"
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
	// We don't have any ms a new style, so we don't need 4 analysis workers
	// for that, we're good by only having 1.
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
		go index()
	}
}

func index() {
	log.Info.Println("Re-indexing search index.")

	count := 0
	start := time.Now()
	action := func(ss []storage.StyleSearch) error {
		b := StyleIndex.NewBatch()

		for _, s := range ss {
			if err := b.Index(strconv.Itoa(s.ID), s); err != nil {
				return err
			}

			count++
			indexMetrics(count, start)
		}

		return StyleIndex.Batch(b)
	}

	if err := storage.FindStylesForSearch(action); err != nil {
		log.Warn.Fatal(err)
	}
}

func indexMetrics(count int, start time.Time) {
	if count%1000 == 0 {
		indexDuration := time.Since(start)
		indexDurationSeconds := float64(indexDuration) / float64(time.Second)
		timePerDoc := float64(indexDuration) / float64(count)
		log.Info.Printf("Indexed %d documents in %.2fs (average %.2fms/doc).",
			count, indexDurationSeconds, timePerDoc/float64(time.Millisecond))
	}
}
