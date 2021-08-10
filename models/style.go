package models

import (
	"fmt"
	"time"

	"gorm.io/gorm"

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

	if err := db().Raw(stmt).Find(q).Error; err != nil {
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

func GetStyleCount() (i int64, err error) {
	if err := db().Select("count(id)").Model(modelStyle).Count(&i).Error; err != nil {
		return 0, err
	}

	return i, nil
}

func GetAllAvailableStylesPaginated(page int, orderStatement string) ([]StyleCard, error) {
	q := new([]StyleCard)
	size := 40
	offset := (page - 1) * size

	// Reflection go brrrr.
	nums := []struct {
		ID, Views, Installs int
	}{}

	var stmt string
	if strings.HasPrefix(orderStatement, "styles") {
		stmt += "styles.id, styles.created_at, styles.updated_at"
	} else if strings.HasPrefix(orderStatement, "views") {
		stmt += "styles.id, (select count(*) from stats s where s.view > 0 and s.style_id = styles.id) views"
	} else {
		stmt += "styles.id, (select count(*) from stats s where s.install > 0 and s.style_id = styles.id) installs"
	}

	err := db().
		Select(stmt).Model(modelStyle).Order(orderStatement).Offset(offset).
		Limit(size).Find(&nums, "styles.deleted_at is null").Error
	if err != nil {
		return nil, err
	}

	var styleIDs []int
	for _, partial := range nums {
		styleIDs = append(styleIDs, int(partial.ID))
	}

	stmt = "styles.id, styles.name, styles.created_at, styles.updated_at, styles.preview, u.username, u.display_name, "
	stmt += "(select count(id) from stats s where s.style_id = styles.id and s.install > 0) installs, "
	stmt += "(select count(id) from stats s where s.style_id = styles.id and s.view > 0) views"

	err = db().
		Select(stmt).Model(modelStyle).Joins("join users u on u.id = styles.user_id").
		Order(orderStatement).Find(&q, styleIDs).Error
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
	q := "styles.name = ? and styles.user_id = ? and styles.code = ?"

	if err := db().First(s, q, s.Name, s.UserID, s.Code).Error; err != nil {
		return errors.ErrDuplicateStyle
	}

	return nil
}

func (s *Style) BanWhereUserID(id interface{}) error {
	return db().Delete(&Style{}, "user_id = ?", id).Error
}
