package models

import (
	"errors"

	"gorm.io/gorm"
)

type Style struct {
	gorm.Model
	UserID      uint
	User        User `gorm:"foreignKey:ID"`
	Name        string
	Summary     string
	Description string
	Code        string
	Preview     string
	Archived    bool   `gorm:"default:false"`
	Featured    bool   `gorm:"default:false"`
	Category    string `gorm:"not null"`
}

type APIStyle struct {
	ID          uint
	Name        string
	Summary     string
	Description string
	Code        string
	Preview     string
	Archived    bool
	Featured    bool
	Category    string
	UserID      uint
	Username    string
}

func GetAllStyles(db *gorm.DB) (*[]APIStyle, error) {
	t := &Style{}
	q := &[]APIStyle{}
	err := db.
		Debug().
		Model(t).
		Select("styles.*, u.username").
		Joins("join users u on u.id = styles.user_id").
		Find(q).
		Error

	if err != nil {
		return nil, errors.New("Styles not found.")
	}

	return q, nil
}

// Using ID as a string is fine in this case.
func GetStyleByID(db *gorm.DB, id string) (*APIStyle, error) {
	t := &Style{}
	q := &APIStyle{}
	err := db.
		Debug().
		Model(t).
		Select("styles.*,  u.username").
		Joins("join users u on u.id = styles.user_id").
		Find(q, "styles.id = ?", id).
		Error

	if err != nil || q.ID == 0 {
		return nil, errors.New("Style not found.")
	}

	return q, nil
}
