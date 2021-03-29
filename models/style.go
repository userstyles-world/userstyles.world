package models

import (
	"errors"
	"log"
	"time"

	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"userstyles.world/config"
)

type Style struct {
	gorm.Model
	UserID      uint
	User        User `gorm:"foreignKey:ID"`
	Name        string
	Description string
	Notes       string
	Code        string
	License     string
	Preview     string
	Homepage    string
	Archived    bool   `gorm:"default:false"`
	Featured    bool   `gorm:"default:false"`
	Category    string `gorm:"not null"`
	Mirror      bool   `gorm:"default:false"`
	Original    string
}

type APIStyle struct {
	ID          uint
	CreatedAt   time.Time
	UpdatedAt   time.Time
	Name        string
	Description string
	Notes       string
	Code        string
	License     string
	Preview     string
	Homepage    string
	Archived    bool
	Featured    bool
	Category    string
	Original    string
	Mirror      bool
	UserID      uint
	Username    string
}

func getDBSession(db *gorm.DB) (tx *gorm.DB) {
	if config.DB_DEBUG == "info" {
		return db.Session(&gorm.Session{
			Logger: db.Logger.LogMode(logger.Info),
		})
	} else {
		return db.Session(&gorm.Session{
			Logger: db.Logger.LogMode(logger.Silent),
		})
	}
}

func GetAllStyles(db *gorm.DB) (*[]APIStyle, error) {
	t, q := new(Style), new([]APIStyle)
	err := getDBSession(db).
		Model(t).
		Select("styles.id, styles.name, styles.preview, u.username").
		Joins("join users u on u.id = styles.user_id").
		Find(q).
		Error

	if err != nil {
		return nil, errors.New("Styles not found.")
	}

	return q, nil
}

func GetAllFeaturedStyles(db *gorm.DB) (*[]APIStyle, error) {
	t, q := new(Style), new([]APIStyle)
	err := getDBSession(db).
		Model(t).
		Joins("join users u on u.id = styles.user_id").
		Select("styles.id, styles.name, styles.preview, u.username").
		Find(q, "styles.featured = ?", true).
		Error

	if err != nil {
		return nil, errors.New("No featured styles.")
	}

	return q, nil
}

func GetAllImportedStyles(db *gorm.DB) (*[]APIStyle, error) {
	t, q := new(Style), new([]APIStyle)
	err := getDBSession(db).
		Model(t).
		Joins("join users u on u.id = styles.user_id").
		Select("styles.id, styles.original, styles.user_id, styles.code").
		Find(q, "styles.original <> '' and styles.mirror = ?", true).
		Error

	if err != nil {
		return nil, errors.New("No imported styles.")
	}

	return q, nil
}

// Using ID as a string is fine in this case.
func GetStyleByID(db *gorm.DB, id string) (*APIStyle, error) {
	t, q := new(Style), new(APIStyle)
	err := getDBSession(db).
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
	t, q := new(Style), new([]APIStyle)
	err := getDBSession(db).
		Model(t).
		Select("styles.id, styles.name, styles.preview, u.username").
		Joins("join users u on u.id = styles.user_id").
		Find(q, "u.username = ?", username).
		Error

	if err != nil {
		return nil, errors.New("Styles not found.")
	}

	return q, nil
}

func CreateStyle(db *gorm.DB, s *Style) (*Style, error) {
	err := getDBSession(db).
		Create(&s).
		Error

	if err != nil {
		return s, err
	}

	return s, nil
}

func GetStyleSourceCodeAPI(db *gorm.DB, id string) (*APIStyle, error) {
	t, q := new(Style), new(APIStyle)
	err := getDBSession(db).
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
