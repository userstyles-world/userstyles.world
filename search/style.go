package search

import (
	"strconv"

	"userstyles.world/modules/storage"
)

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
