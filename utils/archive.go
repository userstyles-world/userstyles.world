package utils

import (
	"userstyles.world/models"
)

const ArchiveURL = "https://raw.githubusercontent.com/33kk/uso-archive/flomaster/data/"

func ImportFromArchive(url string, u models.APIUser) *models.Style {
	s := new(models.Style)
	s.UserID = u.ID
	s.Name = "test"
	return s
}
