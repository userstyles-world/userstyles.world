package models

import (
	"gorm.io/gorm"

	"userstyles.world/modules/database"
)

var (
	modelUser    = User{}
	modelStyle   = Style{}
	modelStats   = Stats{}
	modelHistory = History{}
	modelLog     = Log{}
	modelOAuth   = OAuth{}
	modelReview  = Review{}
)

func db() (tx *gorm.DB) {
	return database.Conn.Session(&gorm.Session{})
}
