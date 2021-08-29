package models

import (
	"gorm.io/gorm"
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

func (Review) FindAllForStyle(id interface{}) (q []Review, err error) {
	err = db().
		// Preload(clause.Associations).
		Preload("User").
		Model(modelReview).
		Order("id DESC").
		Find(&q, "style_id = ? ", id).
		Error

	if err != nil {
		return nil, err
	}

	return q, nil
}

func (r Review) CreateForStyle() error {
	return db().Create(&r).Error
}

func (r *Review) FindLastForStyle(styleID, userID interface{}) error {
	return db().Last(&r, "style_id = ? and user_id = ?", styleID, userID).Error
}

func (r *Review) FindLastFromUser(styleID, userID interface{}) error {
	return db().
		Preload("User").
		Model(modelReview).
		Order("id DESC").
		Last(&r, "style_id = ? and user_id = ?", styleID, userID).
		Error
}
