package models

import (
	"fmt"
	"strings"
	"time"

	"gorm.io/gorm"

	"userstyles.world/modules/config"
	"userstyles.world/modules/errors"
	"userstyles.world/utils/strutils"
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
	UserID      uint `gorm:"index"`
	Archived    bool `gorm:"default:false"`
	Featured    bool `gorm:"default:false"`
	MirrorCode  bool `gorm:"default:false"`
	MirrorMeta  bool `gorm:"default:false"`

	PreviewVersion int `gorm:"default:0"`
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

	PreviewVersion int `json:"-"`
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
	Rating      float64
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
	Views       int64
	Installs    int64
	User        User `gorm:"foreignKey:ID"`
	UserID      uint
	Rating      float64
}

type StyleSiteMap struct {
	ID int
}

func (s StyleCard) Slug() string {
	return strutils.SlugifyURL(s.Name)
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

func GetAllStyles() ([]StyleSearch, error) {
	q := new([]StyleSearch)

	stmt := `
select
	styles.id, styles.created_at, styles.updated_at, styles.name, styles.description, styles.notes,
	styles.preview, u.username, u.display_name,
	(SELECT ROUND(AVG(rating), 1) FROM reviews r WHERE r.style_id = styles.id AND r.deleted_at IS NULL) AS rating,
	(select count(*) from stats s where s.style_id = styles.id and s.install > 0) installs,
	(select count(*) from stats s where s.style_id = styles.id and s.view > 0) views
from
	styles
join
	users u on u.id = styles.user_id
where
	styles.deleted_at is null
`

	if err := db().Raw(stmt).Find(q).Error; err != nil {
		return nil, err
	}

	return *q, nil
}

func GetAllSitesSiteMap() ([]StyleSiteMap, error) {
	q := new([]StyleSiteMap)

	err := db().
		Table("styles").
		Select("id").
		Where("deleted_at is null").
		Scan(q).Error
	if err != nil {
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
	if err := db().Raw(stmt, id).First(q).Error; err != nil {
		return *q, err
	}

	return *q, nil
}

func GetAllStyleIDs() ([]APIStyle, error) {
	q := new([]APIStyle)
	err := db().
		Model(modelStyle).
		Select("styles.id").
		Find(q).
		Error
	if err != nil {
		return nil, errors.ErrStylesNotFound
	}

	return *q, nil
}

func GetAllStylesForIndexAPI() (*[]APIStyle, error) {
	q := new([]APIStyle)

	s := "styles.id, styles.name, styles.created_at, styles.updated_at, "
	s += "styles.description, styles.notes, styles.license, styles.homepage, "
	s += "styles.original, styles.category, styles.preview, styles.user_id, "
	s += "styles.homepage, styles.mirror_url, u.username, u.display_name"

	err := db().
		Model(modelStyle).
		Select(s).
		Joins("join users u on u.id = styles.user_id").
		Find(q).
		Error
	if err != nil {
		return nil, errors.ErrStylesNotFound
	}

	return q, nil
}

func GetStyleCount() (int, error) {
	var c int64
	if err := db().Select("count(id)").Model(modelStyle).Count(&c).Error; err != nil {
		return 0, err
	}

	return int(c), nil
}

func GetAllAvailableStylesPaginated(page int, order string) ([]StyleCard, error) {
	q := new([]StyleCard)
	size := config.AppPageMaxItems
	offset := (page - 1) * size

	// Reflection go brrrr.
	nums := []struct {
		ID, Views, Installs int
	}{}

	var stmt string
	switch {
	case strings.HasPrefix(order, "styles"):
		stmt += "styles.id, styles.created_at, styles.updated_at"
	case strings.HasPrefix(order, "views"):
		stmt += "styles.id, (SELECT sum(daily_views) FROM histories h WHERE h.style_id = styles.id) AS views"
	default:
		stmt += "styles.id, (SELECT sum(daily_installs) FROM histories h WHERE h.style_id = styles.id) AS installs"
	}

	err := db().
		Select(stmt).Model(modelStyle).Order(order).Offset(offset).
		Limit(size).Find(&nums).Error
	if err != nil {
		return nil, err
	}

	styleIDs := make([]int, 0, len(nums))
	for _, partial := range nums {
		styleIDs = append(styleIDs, partial.ID)
	}

	stmt = "styles.id, styles.name, styles.created_at, styles.updated_at, styles.preview, u.username, u.display_name, "
	stmt += "(SELECT ROUND(AVG(rating), 1) FROM reviews r WHERE r.style_id = styles.id AND r.deleted_at IS NULL) AS rating, "
	stmt += "(select count(id) from stats s where s.style_id = styles.id and s.install > 0) installs, "
	stmt += "(select count(id) from stats s where s.style_id = styles.id and s.view > 0) views"

	err = db().
		Select(stmt).Model(modelStyle).Joins("join users u on u.id = styles.user_id").
		Order(order).Find(&q, styleIDs).Error
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

	if err := db().Raw(stmt).Find(q).Error; err != nil {
		return nil, err
	}

	return *q, nil
}

func GetAllFeaturedStyles() ([]StyleCard, error) {
	q := new([]StyleCard)
	stmt := `
select
	styles.id, styles.name, styles.updated_at, styles.preview, u.username, u.display_name,
	(SELECT ROUND(AVG(rating), 1) FROM reviews r WHERE r.style_id = styles.id AND r.deleted_at IS NULL) AS rating,
	(select count(*) from stats s where s.style_id = styles.id and s.install > 0) installs,
	(select count(*) from stats s where s.style_id = styles.id and s.view > 0) views
from
	styles
join
	users u on u.id = styles.user_id
where
	styles.deleted_at is null and styles.featured = 1
`

	if err := db().Raw(stmt).Find(q).Error; err != nil {
		return nil, err
	}

	return *q, nil
}

func GetImportedStyles() ([]Style, error) {
	q := new([]Style)
	err := db().
		Model(modelStyle).
		Find(q, "styles.mirror_url <> '' or styles.original <> '' and styles.mirror_code = ?", true).
		Error
	if err != nil {
		return nil, errors.ErrNoImportedStyles
	}

	return *q, nil
}

// GetStyleByID note: Using ID as a string is fine in this case.
func GetStyleByID(id string) (*APIStyle, error) {
	q := new(APIStyle)
	err := db().
		Model(modelStyle).
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
	(SELECT ROUND(AVG(rating), 1) FROM reviews r WHERE r.style_id = styles.id AND r.deleted_at IS NULL) AS rating,
	(select count(*) from stats s where s.style_id = styles.id and s.install > 0) installs,
	(select count(*) from stats s where s.style_id = styles.id and s.view > 0) views
from
	styles
join
	users u on u.id = styles.user_id
where
	styles.deleted_at is null and u.username = ?
`

	if err := db().Raw(stmt, username).Find(q).Error; err != nil {
		return nil, err
	}

	return *q, nil
}

func CreateStyle(s *Style) (*Style, error) {
	if err := db().Create(&s).Error; err != nil {
		return s, err
	}

	return s, nil
}

func UpdateStyle(s *Style) error {
	err := db().
		Model(modelStyle).
		Where("id", s.ID).
		Updates(s).
		Error
	if err != nil {
		return err
	}

	return nil
}

func GetStyleSourceCodeAPI(id string) (*APIStyle, error) {
	q := new(APIStyle)
	err := db().
		Model(modelStyle).
		Select("styles.*, u.username").
		Joins("join users u on u.id = styles.user_id").
		First(q, "styles.id = ?", id).
		Error
	if err != nil {
		return q, err
	}

	return q, nil
}

func CheckDuplicateStyle(s *Style) error {
	q := "styles.name = ? and styles.user_id = ?"

	if err := db().First(s, q, s.Name, s.UserID).Error; err == nil {
		return errors.ErrDuplicateStyle
	}

	return nil
}

func (*Style) BanWhereUserID(id any) error {
	return db().Delete(&Style{}, "user_id = ?", id).Error
}

// MirrorStyle will update fields depending on which mirror option is used.
func (*Style) MirrorStyle(f map[string]any) error {
	err := db().Model(modelStyle).Where("id", f["id"]).Updates(f).Error
	if err != nil {
		return err
	}

	return nil
}

func (s *Style) UpdateColumn(col string, val any) error {
	return db().Model(modelStyle).Where("id", s.ID).UpdateColumn(col, val).Error
}

// SetPreview will set preview image URL.
func (s *Style) SetPreview() {
	s.Preview = fmt.Sprintf("%s/preview/%d/%dt.webp", config.BaseURL, s.ID, s.PreviewVersion)
}

// SetPreview will set preview image URL.
func (s *APIStyle) SetPreview() {
	s.Preview = fmt.Sprintf("%s/preview/%d/%dt.webp", config.BaseURL, s.ID, s.PreviewVersion)
}

// SelectUpdateStyle will update specific fields in the styles table.
func SelectUpdateStyle(s *APIStyle) error {
	fields := []string{"name", "description", "notes", "code", "homepage",
		"license", "category", "preview", "preview_version", "mirror_url",
		"mirror_code", "mirror_meta"}

	return db().
		Model(modelStyle).
		Select(fields).
		Where("id = ?", s.ID).
		Updates(s).Error
}
