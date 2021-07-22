package search

import (
	"errors"
	"strconv"
	"time"

	"github.com/blevesearch/bleve/v2"

	"userstyles.world/models"
	"userstyles.world/modules/log"
)

var (
	StyleIndex bleve.Index
	batchSize  = 500
)

func Initialize() {
	stylesIndex, err := bleve.Open("data/styles.bleve")
	if errors.Is(err, bleve.ErrorIndexPathDoesNotExist) {
		log.Info.Println("Creating new index...")
		indexMapping := buildIndexMapping()
		stylesIndex, err = bleve.New("data/styles.bleve", indexMapping)
		if err != nil {
			log.Warn.Fatal(err)
		}
	} else if err != nil {
		log.Warn.Fatal(err)
	}
	log.Info.Println("Opening existing index...")

	StyleIndex = stylesIndex

	go func() {
		styleEntries, err := models.GetAllStyles()
		if err != nil {
			log.Warn.Fatal(err)
		}
		err = indexStyles(stylesIndex, styleEntries)
		if err != nil {
			log.Warn.Fatal(err)
		}
	}()
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
