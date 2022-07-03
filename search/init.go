package search

import (
	"errors"
	"os"
	"path"
	"strconv"
	"time"

	"github.com/blevesearch/bleve/v2"
	"github.com/blevesearch/bleve/v2/index/upsidedown"

	"userstyles.world/models"
	"userstyles.world/modules/config"
	"userstyles.world/modules/log"
)

var (
	StyleIndex bleve.Index
	batchSize  = 500
	indexFile  = path.Join(config.DataDir, "styles.bleve")
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

func Initialize() {
	// We don't have any ms a new style, so we don't need 4 analysis workers
	// for that, we're good by only having 1.
	bleve.Config.SetAnalysisQueueSize(1)

	stylesIndex, err := openBleveIndexFile(indexFile)
	if err != nil {
		log.Info.Println("Creating new index...")
		indexMapping := buildIndexMapping()
		stylesIndex, err = bleve.New(indexFile, indexMapping)
		if err != nil {
			log.Warn.Fatal(err)
		}
	}
	log.Info.Println("Opening existing index...")

	StyleIndex = stylesIndex

	if config.SearchReindex {
		go func() {
			styleEntries, err := models.GetAllStyles()
			if err != nil {
				log.Warn.Fatal(err)
			}
			err = indexStyles(StyleIndex, styleEntries)
			if err != nil {
				log.Warn.Fatal(err)
			}
		}()
	}
}

func indexStyles(index bleve.Index, data []models.StyleSearch) error {
	count := 0
	startTime := time.Now()
	batch := index.NewBatch()
	batchCount := 0
	var err error
	for _, styleEntry := range data {
		id := strconv.Itoa(styleEntry.ID)
		err = batch.Index(id, MinimalStyle{
			ID:          styleEntry.ID,
			CreatedAt:   styleEntry.CreatedAt,
			UpdatedAt:   styleEntry.UpdatedAt,
			Username:    styleEntry.Username,
			DisplayName: styleEntry.DisplayName,
			Name:        styleEntry.Name,
			Description: styleEntry.Description,
			Preview:     styleEntry.Preview,
			Notes:       styleEntry.Notes,
			Installs:    styleEntry.Installs,
			Views:       styleEntry.Views,
			Rating:      styleEntry.Rating,
		})
		if err != nil {
			return err
		}

		batchCount++

		if batchCount >= batchSize {
			err := index.Batch(batch)
			if err != nil {
				return err
			}
			batch = index.NewBatch()
			batchCount = 0
		}
		count++
		if count%1000 == 0 {
			indexDuration := time.Since(startTime)
			indexDurationSeconds := float64(indexDuration) / float64(time.Second)
			timePerDoc := float64(indexDuration) / float64(count)
			log.Info.Printf("Indexed %d documents, in %.2fs (average %.2fms/doc)",
				count, indexDurationSeconds, timePerDoc/float64(time.Millisecond))
		}
	}
	// flush the last batch
	if batchCount > 0 {
		err := index.Batch(batch)
		if err != nil {
			log.Warn.Fatal(err)
		}
	}
	indexDuration := time.Since(startTime)
	indexDurationSeconds := float64(indexDuration) / float64(time.Second)
	timePerDoc := float64(indexDuration) / float64(count)
	log.Info.Printf("Indexed %d documents, in %.2fs (average %.2fms/doc)",
		count, indexDurationSeconds, timePerDoc/float64(time.Millisecond))

	return nil
}
