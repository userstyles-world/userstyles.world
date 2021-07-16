package models

import (
	"gorm.io/gorm"

	"userstyles.world/modules/database"
)

type Review struct {
	gorm.Model
	Rating  int
	Comment string

	User   User
	UserID int

	Style   Style
	StyleID int
}

func (r Review) FindAllForStyle(id interface{}) (q []Review, err error) {
	err = database.Conn.
		Debug().
		// Preload(clause.Associations).
		Preload("User").
		Model(Review{}).
		Order("id DESC").
		Find(&q, "style_id = ? ", id).
		Error

	if err != nil {
		return nil, err
	}

	return q, nil
}

func (r Review) CreateForStyle(id interface{}) error {
	return database.Conn.Debug().Create(&r).Error
}
