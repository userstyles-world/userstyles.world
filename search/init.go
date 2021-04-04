package search

import (
	"log"
	"strconv"
	"time"

	"github.com/blevesearch/bleve/v2"
	"github.com/gofiber/fiber/v2"

	"userstyles.world/database"
	"userstyles.world/models"
)

var (
	StyleIndex bleve.Index
	batchSize  = 25
)

func Initialize() {
	stylesIndex, err := bleve.Open("styles.bleve")
	if err == bleve.ErrorIndexPathDoesNotExist {
		log.Printf("Creating new index...")
		// create a mapping
		indexMapping, err := buildIndexMapping()
		if err != nil {
			log.Fatal(err)
		}
		stylesIndex, err = bleve.New("styles.bleve", indexMapping)
		if err != nil {
			log.Fatal(err)
		}
	} else if err != nil {
		log.Fatal(err)
	} else {
		log.Printf("Opening existing index...")
	}

	StyleIndex = stylesIndex

	go func() {
		err = indexStyles(stylesIndex)
		if err != nil {
			log.Fatal(err)
		}
	}()
}

func indexStyles(index bleve.Index) error {
	styleEntries, err := models.GetAllStyles(database.DB)

	log.Printf("Indexing...")
	count := 0
	startTime := time.Now()
	batch := index.NewBatch()
	batchCount := 0
	for _, styleEntry := range *styleEntries {
		ID := strconv.FormatUint(uint64(styleEntry.ID), 10)
		batch.Index(ID, fiber.Map{
			"name":        styleEntry.Name,
			"description": styleEntry.Description,
			"notes":       styleEntry.Notes,
			"user":        styleEntry.Username,
		})

		batchCount++

		if batchCount >= batchSize {
			err = index.Batch(batch)
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
			log.Printf("Indexed %d documents, in %.2fs (average %.2fms/doc)", count, indexDurationSeconds, timePerDoc/float64(time.Millisecond))
		}
	}
	// flush the last batch
	if batchCount > 0 {
		err = index.Batch(batch)
		if err != nil {
			log.Fatal(err)
		}
	}
	indexDuration := time.Since(startTime)
	indexDurationSeconds := float64(indexDuration) / float64(time.Second)
	timePerDoc := float64(indexDuration) / float64(count)
	log.Printf("Indexed %d documents, in %.2fs (average %.2fms/doc)", count, indexDurationSeconds, timePerDoc/float64(time.Millisecond))
	return nil
}
