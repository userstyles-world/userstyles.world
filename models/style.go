package models

import (
	stderrors "errors"
	"fmt"
	"time"

	"github.com/go-playground/validator/v10"
	"gorm.io/gorm"

	"userstyles.world/modules/config"
	"userstyles.world/modules/errors"
	"userstyles.world/utils"
)

type Style struct {
	gorm.Model
	Original    string
	MirrorURL   string
	Homepage    string
	Category    string `validate:"required,min=1,max=255" gorm:"not null"`
	Name        string `validate:"required,min=1,max=50"`
	Description string `validate:"required,min=1,max=160"`
	Notes       string `validate:"min=0,max=50000"`
	Code        string `validate:"max=10000000"`
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

type StyleSiteMap struct {
	ID int
}

// TruncateCode returns if it should the style, to prevent long loading times.
func (s APIStyle) TruncateCode() bool {
	return len(s.Code) > 100_000
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

// GetStyle tries to fetch a userstyle from database.
func GetStyle(id string) (Style, error) {
	var s Style
	err := db().
		Select("styles.*, u.username").
		Joins("JOIN users u ON u.id = styles.user_id").
		First(&s, "styles.id = ?", id).Error
	return s, err
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

var (
	ErrStyleNoData = stderrors.New("missing mandatory userstyle data")
	ErrStyleAsErrs = stderrors.New("unexpected error during validation")
)

// Validate makes sure input data is correct.
func (s Style) Validate() (map[string]any, string, error) {
	fields := []string{"Name", "Description", "Notes", "Category", "Code"}
	err := utils.Validate().StructPartial(s, fields...)
	if err == nil {
		return nil, "", nil
	}

	var errs validator.ValidationErrors
	if ok := stderrors.As(err, &errs); !ok {
		return nil, "Unexpected error during validation.", ErrStyleAsErrs
	}

	m := make(map[string]any, len(errs))
	for _, e := range errs {
		var msg string

		switch e.Field() {
		case "Name":
			msg = "Name must be up to 50 characters."
		case "Description":
			msg = "Description must be up to 160 characters."
		case "Notes":
			msg = "Notes must be up to 50K characters."
		case "Category":
			msg = "Category must be up to 255 characters."
		case "Code":
			msg = "Code must be up to 10M characters."
		}

		m[e.Field()] = msg
	}

	return m, "Missing mandatory userstyle data.", ErrStyleNoData
}

// SetPreview will set preview image URL.
func (s *APIStyle) SetPreview() {
	s.Preview = fmt.Sprintf("%s/preview/%d/%dt.webp", config.BaseURL, s.ID, s.PreviewVersion)
}

// SelectUpdateStyle will update specific fields in the styles table.
func SelectUpdateStyle(s Style) error {
	fields := []string{"name", "description", "notes", "code", "homepage",
		"license", "category", "preview", "preview_version", "mirror_url",
		"mirror_code", "mirror_meta"}

	return db().
		Model(modelStyle).
		Select(fields).
		Where("id = ?", s.ID).
		Updates(s).Error
}

// CompareMirrorURL will return true if a style is being imported and mirrored
// from the same URL.
func (s *APIStyle) CompareMirrorURL() bool {
	if s.Original != "" &&
		(s.MirrorCode || s.MirrorMeta) &&
		(s.MirrorURL == "" || s.MirrorURL == s.Original) {
		return true
	}

	return false
}
