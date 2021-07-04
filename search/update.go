package search

import (
	"strconv"

	"userstyles.world/models"
)

func IndexStyle(id uint) error {
	styleID := strconv.FormatUint(uint64(id), 10)

	style, err := models.GetStyleForIndex(styleID)
	if err != nil {
		return err
	}

	if err := StyleIndex.Index(styleID, MinimalStyle{
		ID:          style.ID,
		CreatedAt:   style.CreatedAt,
		UpdatedAt:   style.UpdatedAt,
		Username:    style.Username,
		DisplayName: style.DisplayName,
		Name:        style.Name,
		Description: style.Description,
		Preview:     style.Preview,
		Notes:       style.Notes,
		Installs:    style.Installs,
		Views:       style.Views,
	}); err != nil {
		return err
	}

	return nil
}

func DeleteStyle(id uint) error {
	return StyleIndex.Delete(strconv.FormatUint(uint64(id), 10))
}
