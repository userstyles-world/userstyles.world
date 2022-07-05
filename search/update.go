package search

import (
	"strconv"

	"userstyles.world/models"
	"userstyles.world/modules/storage"
)

func IndexStyle(id uint) error {
	styleID := strconv.FormatUint(uint64(id), 10)

	style, err := models.GetStyleForIndex(styleID)
	if err != nil {
		return err
	}

	return StyleIndex.Index(styleID, storage.StyleSearch{
		ID:          style.ID,
		Username:    style.Username,
		Name:        style.Name,
		Description: style.Description,
		Notes:       style.Notes,
	})
}

func DeleteStyle(id uint) error {
	return StyleIndex.Delete(strconv.FormatUint(uint64(id), 10))
}
