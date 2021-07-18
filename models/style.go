package models

import (
	"fmt"
	"log"
	"time"

	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"userstyles.world/modules/config"
	"userstyles.world/modules/database"
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
	CreatedAt   time.Time
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

// TruncateCode returns if it should the style, to prevent long loading times.
func (s APIStyle) TruncateCode() bool {
	return len(s.Code) > 100_000
}

func getDBSession() (tx *gorm.DB) {
	var logLevel logger.LogLevel
	switch config.DBDebug {
	case "error":
		logLevel = logger.Error
	case "warn":
		logLevel = logger.Warn
	case "info":
		logLevel = logger.Info
	default:
		logLevel = logger.Silent
	}

	return database.Conn.Session(&gorm.Session{
		Logger: database.Conn.Logger.LogMode(logLevel),
	})
}

func GetAllStyles() ([]StyleSearch, error) {
	q := new([]StyleSearch)

	stmt := `
select
	styles.id, styles.created_at, styles.updated_at, styles.name, styles.description, styles.notes,
	styles.preview, u.username, u.display_name,
	(select count(*) from stats s where s.style_id = styles.id and s.install > 0) installs,
	(select count(*) from stats s where s.style_id = styles.id and s.view > 0) views
from
	styles
join
	users u on u.id = styles.user_id
where
	styles.deleted_at is null
`

	if err := getDBSession().Raw(stmt).Find(q).Error; err != nil {
		return nil, err
	}

	return *q, nil
}

func GetStyleForIndex(id string) (StyleSearch, error) {
	q := new(StyleSearch)

	stmt := `
select
	styles.id, styles.created_at, styles.updated_at, styles.name, styles.description, styles.notes,
	styles.preview, u.username, u.display_name,
	(select count(*) from stats s where s.style_id = styles.id and s.install > 0) installs,
	(select count(*) from stats s where s.style_id = styles.id and s.view > 0) views
from
	styles
join
	users u on u.id = styles.user_id
where
	styles.deleted_at is null and styles.id = ?
`
	if err := getDBSession().Raw(stmt, id).First(q).Error; err != nil {
		return *q, err
	}

	return *q, nil
}

func GetAllStyleIDs() ([]APIStyle, error) {
	t, q := new(Style), new([]APIStyle)
	err := getDBSession().
		Model(t).
		Select("styles.id").
		Find(q).
		Error
	if err != nil {
		return nil, errors.ErrStylesNotFound
	}

	return *q, nil
}

func GetAllStylesForIndexAPI() (*[]APIStyle, error) {
	t, q := new(Style), new([]APIStyle)

	s := "styles.id, styles.name, styles.created_at, styles.updated_at, "
	s += "styles.description, styles.notes, styles.license, styles.homepage, "
	s += "styles.original, styles.category, styles.preview, styles.user_id, "
	s += "styles.homepage, styles.mirror_url, u.username, u.display_name"

	err := getDBSession().
		Model(t).
		Select(s).
		Joins("join users u on u.id = styles.user_id").
		Find(q).
		Error
	if err != nil {
		return nil, errors.ErrStylesNotFound
	}

	return q, nil
}

func GetStyleCount() (i int64, err error) {
	if err := database.Conn.Select("count(id)").Model(Style{}).Count(&i).Error; err != nil {
		return 0, err
	}

	return i, nil
}

func GetAllAvailableStylesPaginated(page int) ([]StyleCard, error) {
	q := new([]StyleCard)
	size := 40
	offset := (page - 1) * size

	s1 := "styles.id, styles.name, styles.created_at, styles.updated_at, styles.preview, u.username, u.display_name, "
	s2 := "(select count(id) from stats s where s.style_id = styles.id and s.install > 0) installs, "
	s3 := "(select count(id) from stats s where s.style_id = styles.id and s.view > 0) views"
	stmt := s1 + s2 + s3

	err := getDBSession().
		Select(stmt).
		Model(Style{}).
		Joins("join users u on u.id = styles.user_id").
		Offset(offset).
		Limit(size).
		Find(q).Error

	if err != nil {
		return nil, err
	}

	return *q, nil
}

func GetAllAvailableStyles() ([]StyleCard, error) {
	q := new([]StyleCard)
	stmt := `
select
	styles.id, styles.name, styles.created_at, styles.updated_at, styles.preview, u.username, u.display_name,
	(select count(*) from stats s where s.style_id = styles.id and s.install > 0) installs,
	(select count(*) from stats s where s.style_id = styles.id and s.view > 0) views
from
	styles
join
	users u on u.id = styles.user_id
where
	styles.deleted_at is null
`

	if err := getDBSession().Raw(stmt).Find(q).Error; err != nil {
		return nil, err
	}

	return *q, nil
}

func GetAllFeaturedStyles() ([]StyleCard, error) {
	q := new([]StyleCard)
	stmt := `
select
	styles.id, styles.name, styles.updated_at, styles.preview, u.username, u.display_name,
	(select count(*) from stats s where s.style_id = styles.id and s.install > 0) installs,
	(select count(*) from stats s where s.style_id = styles.id and s.view > 0) views
from
	styles
join
	users u on u.id = styles.user_id
where
	styles.deleted_at is null and styles.featured = 1
`

	if err := getDBSession().Raw(stmt).Find(q).Error; err != nil {
		return nil, err
	}

	return *q, nil
}

func GetImportedStyles() ([]Style, error) {
	t, q := new(Style), new([]Style)
	err := getDBSession().
		Model(t).
		Find(q, "styles.mirror_url <> '' or styles.original <> '' and styles.mirror_code = ?", true).
		Error
	if err != nil {
		return nil, errors.ErrNoImportedStyles
	}

	return *q, nil
}

// GetStyleByID note: Using ID as a string is fine in this case.
func GetStyleByID(id string) (*APIStyle, error) {
	t, q := new(Style), new(APIStyle)
	err := getDBSession().
		Model(t).
		Select("styles.*,  u.username").
		Joins("join users u on u.id = styles.user_id").
		Find(q, "styles.id = ?", id).
		Error

	if err != nil || q.ID == 0 {
		return nil, errors.ErrStyleNotFound
	}

	return q, nil
}

func GetStylesByUser(username string) ([]StyleCard, error) {
	q := new([]StyleCard)
	stmt := `
select
	styles.id, styles.name, styles.updated_at, styles.preview, u.username, u.display_name,
	(select count(*) from stats s where s.style_id = styles.id and s.install > 0) installs,
	(select count(*) from stats s where s.style_id = styles.id and s.view > 0) views
from
	styles
join
	users u on u.id = styles.user_id
where
	styles.deleted_at is null and u.username = ?
`

	if err := getDBSession().Raw(stmt, username).Find(q).Error; err != nil {
		return nil, err
	}

	return *q, nil
}

func CreateStyle(s *Style) (*Style, error) {
	err := getDBSession().
		Create(&s).
		Error
	if err != nil {
		return s, err
	}

	return s, nil
}

func UpdateStyle(s *Style) error {
	err := getDBSession().
		Model(Style{}).
		Where("id", s.ID).
		Updates(s).
		Error
	if err != nil {
		return err
	}

	return nil
}

func GetStyleSourceCodeAPI(id string) (*APIStyle, error) {
	t, q := new(Style), new(APIStyle)
	err := getDBSession().
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

func CheckDuplicateStyle(s *Style) error {
	q := "styles.name = ? and styles.user_id = ? and styles.code = ?"
	err := getDBSession().
		First(s, q, s.Name, s.UserID, s.Code).
		Error

	if err == nil {
		return errors.ErrDuplicateStyle
	}

	return nil
}

func (s *Style) BanWhereUserID(id interface{}) error {
	return database.Conn.Delete(&Style{}, "user_id = ?", id).Error
}
