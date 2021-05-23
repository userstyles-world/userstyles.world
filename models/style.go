package models

import (
	"fmt"
	"log"
	"time"

	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"userstyles.world/config"
	"userstyles.world/modules/errors"
	"userstyles.world/utils/strings"
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
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	Category    string    `json:"category"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Notes       string    `json:"notes"`
	Code        string    `json:"-"`
	License     string    `json:"license"`
	Preview     string    `json:"preview_url"`
	Homepage    string    `json:"homepage"`
	Username    string    `json:"username"`
	Original    string    `json:"original"`
	MirrorURL   string    `json:"mirror_url"`
	DisplayName string    `json:"display_name"`
	UserID      uint      `json:"user_id"`
	ID          uint      `json:"id"`
	Featured    bool      `json:"-"`
	MirrorCode  bool      `json:"-"`
	MirrorMeta  bool      `json:"-"`
	Archived    bool      `json:"-"`
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

type StyleSearch struct {
	ID          int
	UpdatedAt   time.Time
	Name        string
	Description string
	Notes       string
	Preview     string
	DisplayName string
	Username    string
	Views       int
	Installs    int
	User        User `gorm:"foreignKey:ID"`
	UserID      uint
}

func (s StyleCard) Slug() string {
	return strings.SlugifyURL(s.Name)
}

func (s StyleCard) StyleURL() string {
	return fmt.Sprintf("/style/%d/%s", s.ID, s.Slug())
}

func (s StyleCard) Author() string {
	if s.DisplayName != "" {
		return s.DisplayName
	}

	return s.Username
}

// Truncate large styles to prevent long loading times.
func (s APIStyle) TruncateCode() bool {
	if len(s.Code) > 150_000 {
		return true
	}

	return false
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

func GetAllStyles(db *gorm.DB) ([]StyleSearch, error) {
	q := new([]StyleSearch)

	stmt := `
select
	styles.id, styles.updated_at, styles.name, styles.description, styles.notes,
	styles.preview, u.username, u.display_name,
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

func GetAllStyleIDs(db *gorm.DB) ([]APIStyle, error) {
	t, q := new(Style), new([]APIStyle)
	err := getDBSession(db).
		Model(t).
		Select("styles.id").
		Find(q).
		Error
	if err != nil {
		return nil, errors.StylesNotFound
	}

	return *q, nil
}

func GetAllStylesForIndexAPI(db *gorm.DB) (*[]APIStyle, error) {
	t, q := new(Style), new([]APIStyle)

	s := "styles.id, styles.name, styles.created_at, styles.updated_at, "
	s += "styles.description, styles.notes, styles.license, styles.homepage, "
	s += "styles.original, styles.category, styles.preview, styles.user_id, "
	s += "styles.homepage, styles.mirror_url, u.username, u.display_name"

	err := getDBSession(db).
		Model(t).
		Select(s).
		Joins("join users u on u.id = styles.user_id").
		Find(q).
		Error
	if err != nil {
		return nil, errors.StylesNotFound
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
		return nil, errors.NoImportedStyles
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
		return nil, errors.StyleNotFound
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
		return errors.DuplicateStyle
	}

	return nil
}
