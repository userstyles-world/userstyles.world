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
	Original    string
	MirrorURL   string
	Homepage    string
	Category    string `gorm:"not null"`
	Name        string
	Description string
	Notes       string
	Code        string
	License     string
	Preview     string
	User        User `gorm:"foreignKey:ID"`
	UserID      uint
	Archived    bool `gorm:"default:false"`
	Featured    bool `gorm:"default:false"`
	MirrorCode  bool `gorm:"default:false"`
	MirrorMeta  bool `gorm:"default:false"`
}

type APIStyle struct {
	CreatedAt   time.Time
	UpdatedAt   time.Time
	Category    string
	Name        string
	Description string
	Notes       string
	Code        string
	License     string
	Preview     string
	Homepage    string
	Username    string
	Original    string
	MirrorURL   string
	DisplayName string
	UserID      uint
	ID          uint
	Featured    bool
	MirrorCode  bool
	MirrorMeta  bool
	Archived    bool
}

type StyleCard struct {
	gorm.Model
	Views       int64
	Installs    int64
	Name        string
	Preview     string
	DisplayName string
	Username    string
	User        User `gorm:"foreignKey:ID"`
	UserID      uint
}

func getDBSession(db *gorm.DB) (tx *gorm.DB) {
	var log logger.LogLevel
	switch config.DB_DEBUG {
	case "error":
		log = logger.Error
	case "warn":
		log = logger.Warn
	case "info":
		log = logger.Info
	default:
		log = logger.Silent
	}

	return db.Session(&gorm.Session{
		Logger: db.Logger.LogMode(log),
	})
}

func GetAllStyles(db *gorm.DB) (*[]APIStyle, error) {
	t, q := new(Style), new([]APIStyle)
	err := getDBSession(db).
		Model(t).
		Select("styles.id, styles.name, styles.description, styles.notes, " +
			"styles.category, styles.preview, u.username, u.display_name").
		Joins("join users u on u.id = styles.user_id").
		Find(q).
		Error
	if err != nil {
		return nil, errors.New("styles not found")
	}

	return q, nil
}

func GetAllStyleIDs(db *gorm.DB) ([]APIStyle, error) {
	t, q := new(Style), new([]APIStyle)
	err := getDBSession(db).
		Model(t).
		Select("styles.id").
		Find(q).
		Error
	if err != nil {
		return nil, errors.New("styles not found")
	}

	return *q, nil
}

func GetAllStylesForIndexAPI(db *gorm.DB) (*[]APIStyle, error) {
	t, q := new(Style), new([]APIStyle)

	s := "styles.id, styles.name, styles.created_at, styles.updated_at, "
	s += "styles.description, styles.notes, styles.category, styles.preview, u.username"

	err := getDBSession(db).
		Model(t).
		Select(s).
		Joins("join users u on u.id = styles.user_id").
		Find(q).
		Error
	if err != nil {
		return nil, errors.New("styles not found")
	}

	return q, nil
}

func GetAllAvailableStyles(db *gorm.DB) ([]StyleCard, error) {
	q := new([]StyleCard)
	stmt := `
select
	styles.id, styles.name, styles.updated_at, styles.preview, u.username, u.display_name,
	(select count(*) from stats s where s.style_id = styles.id and s.install = 1) installs,
	(select count(*) from stats s where s.style_id = styles.id and s.view = 1) views
from
	styles
join
	users u on u.id = styles.user_id
where
	styles.deleted_at is null
`

	if err := getDBSession(db).Raw(stmt).Find(q).Error; err != nil {
		return nil, err
	}

	return *q, nil
}

func GetAllFeaturedStyles(db *gorm.DB) ([]StyleCard, error) {
	q := new([]StyleCard)
	stmt := `
select
	styles.id, styles.name, styles.updated_at, styles.preview, u.username, u.display_name,
	(select count(*) from stats s where s.style_id = styles.id and s.install = 1) installs,
	(select count(*) from stats s where s.style_id = styles.id and s.view = 1) views
from
	styles
join
	users u on u.id = styles.user_id
where
	styles.deleted_at is null and styles.featured = 1
`

	if err := getDBSession(db).Raw(stmt).Find(q).Error; err != nil {
		return nil, err
	}

	return *q, nil
}

func GetImportedStyles(db *gorm.DB) ([]Style, error) {
	t, q := new(Style), new([]Style)
	err := getDBSession(db).
		Model(t).
		Find(q, "styles.mirror_url <> '' or styles.original <> '' and styles.mirror_code = ?", true).
		Error
	if err != nil {
		return nil, errors.New("no imported styles")
	}

	return *q, nil
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
		return nil, errors.New("style not found")
	}

	return q, nil
}

func GetStylesByUser(db *gorm.DB, username string) ([]StyleCard, error) {
	q := new([]StyleCard)
	stmt := `
select
	styles.id, styles.name, styles.updated_at, styles.preview, u.username, u.display_name,
	(select count(*) from stats s where s.style_id = styles.id and s.install = 1) installs,
	(select count(*) from stats s where s.style_id = styles.id and s.view = 1) views
from
	styles
join
	users u on u.id = styles.user_id
where
	styles.deleted_at is null and u.username = ?
`

	if err := getDBSession(db).Raw(stmt, username).Find(q).Error; err != nil {
		return nil, err
	}


	return *q, nil
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

func UpdateStyle(db *gorm.DB, s *Style) error {
	err := getDBSession(db).
		Model(Style{}).
		Where("id", s.ID).
		Updates(s).
		Error
	if err != nil {
		return err
	}

	return nil
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

func CheckDuplicateStyle(db *gorm.DB, s *Style) error {
	q := "styles.name = ? and styles.user_id = ? and styles.code = ?"
	err := getDBSession(db).
		First(s, q, s.Name, s.UserID, s.Code).
		Error

	if err == nil {
		return errors.New("duplicate style")
	}

	return nil
}
