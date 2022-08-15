package search

import (
	"strconv"
	"time"

	"userstyles.world/modules/log"
	"userstyles.world/modules/storage"
)

// indexStyles adds all styles to the indexStyles in batches.
func indexStyles() {
	log.Info.Println("Re-indexing search engine...")

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

// indexMetrics reports the time it took to index all styles.
func indexMetrics(count int, start time.Time) {
	if count%1000 == 0 {
		indexDuration := time.Since(start)
		indexDurationSeconds := float64(indexDuration) / float64(time.Second)
		timePerDoc := float64(indexDuration) / float64(count)
		log.Info.Printf("Indexed %d documents in %.2fs (average %.2fms/doc).",
			count, indexDurationSeconds, timePerDoc/float64(time.Millisecond))
	}
}

// IndexStyle adds a new style to the index.
func IndexStyle(id uint) error {
	res, err := storage.FindStyleForSearch(id)
	if err != nil {
		return err
	}

	return StyleIndex.Index(strconv.FormatUint(uint64(id), 10), res)
}

// DeleteStyle removes a style from the index.
func DeleteStyle(id uint) error {
	return StyleIndex.Delete(strconv.FormatUint(uint64(id), 10))
}
