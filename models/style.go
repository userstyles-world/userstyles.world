package models

import (
	"errors"
	"log"
	"time"

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
	Homepage    string
	Archived    bool   `gorm:"default:false"`
	Featured    bool   `gorm:"default:false"`
	Category    string `gorm:"not null"`
}

type APIStyle struct {
	ID          uint
	CreatedAt   time.Time
	UpdatedAt   time.Time
	Name        string
	Summary     string
	Description string
	Code        string
	Preview     string
	Homepage    string
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

func GetStylesByUser(db *gorm.DB, username string) (*[]APIStyle, error) {
	t := &Style{}
	q := &[]APIStyle{}
	err := db.
		Debug().
		Model(t).
		Select("styles.*,  u.username").
		Joins("join users u on u.id = styles.user_id").
		Find(q, "u.username = ?", username).
		Error

	if err != nil {
		return nil, errors.New("Style not found.")
	}

	return q, nil
}

func GetStyleSourceCodeAPI(db *gorm.DB, id string) (*APIStyle, error) {
	t, q := new(Style), new(APIStyle)
	err := db.
		Debug().
		Model(t).
		Select("styles.*, u.username").
		Joins("join users u on u.id = styles.user_id").
		First(q, "styles.id = ?", id).
		Error

	if err != nil {
		log.Printf("Problem with style id %s, err: %v", id, err)
		return q, err
	}

	return q, nil
}
