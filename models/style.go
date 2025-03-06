package models

import (
	stderrors "errors"
	"fmt"
	"hash/crc32"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/vednoc/go-usercss-parser"
	"gorm.io/gorm"

	"userstyles.world/modules/config"
	"userstyles.world/modules/database"
	"userstyles.world/modules/errors"
	"userstyles.world/modules/util"
)

type Style struct {
	gorm.Model
	Original       string
	MirrorURL      string
	Homepage       string
	Category       string `validate:"required,min=1,max=255" gorm:"not null"`
	Name           string `validate:"required,min=1,max=50"`
	Description    string `validate:"required,min=1,max=160"`
	Notes          string `validate:"min=0,max=50000"`
	Code           string `validate:"max=10000000"`
	CodeChecksum   string `json:"-"`
	License        string
	Preview        string
	Slug           string `json:"-"`
	User           User   `gorm:"foreignKey:ID"`
	UserID         uint   `gorm:"index"`
	Archived       bool   `gorm:"default:false"`
	Featured       bool   `gorm:"default:false"`
	MirrorCode     bool   `gorm:"default:false"`
	MirrorMeta     bool   `gorm:"default:false"`
	PreviewVersion int    `gorm:"default:0"`
	ImportPrivate  bool   `gorm:"default:false"`
	MirrorPrivate  bool   `gorm:"default:false"`
	CodeSize       uint64 `json:"-"`
}

type APIStyle struct {
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
	Category       string    `json:"category"`
	Name           string    `json:"name"`
	Description    string    `json:"description"`
	Notes          string    `json:"notes"`
	Code           string    `json:"-"`
	License        string    `json:"license"`
	Preview        string    `json:"preview_url"`
	Homepage       string    `json:"homepage"`
	Username       string    `json:"username"`
	Original       string    `json:"original"`
	MirrorURL      string    `json:"mirror_url"`
	DisplayName    string    `json:"display_name"`
	UserID         uint      `json:"user_id"`
	ID             uint      `json:"id"`
	Featured       bool      `json:"-"`
	MirrorCode     bool      `json:"-"`
	MirrorMeta     bool      `json:"-"`
	Archived       bool      `json:"-"`
	PreviewVersion int       `json:"-"`
	ImportPrivate  bool      `json:"-"`
	MirrorPrivate  bool      `json:"-"`
}

// TableName returns which table in database to use with GORM.
func (APIStyle) TableName() string { return "styles" }

// Permalink returns a link to the style page.
func (s Style) Permalink() string {
	return fmt.Sprintf("/style/%d/%s", s.ID, s.Slug)
}

// Prepare sets dynamic fields to their respective values.
func (s *Style) Prepare() {
	s.Slug = util.Slug(s.Name)
	s.CodeSize = uint64(len(s.Code))
	s.CodeChecksum = fmt.Sprintf("%x", crc32.ChecksumIEEE([]byte(s.Code)))
}

type StyleSiteMap struct {
	ID int
}

// TruncateCode returns if it should the style, to prevent long loading times.
func (s Style) TruncateCode() bool {
	return len(s.Code) > 10_000
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

// GetStyleByID tries to fetch a userstyle with id from the database.
func GetStyleByID(id int) (s *Style, err error) {
	err = database.Conn.Joins("User").Find(&s, "styles.id = ?", id).Error
	return s, err
}

func CreateStyle(s *Style) (*Style, error) {
	s.Prepare()

	err := database.Conn.Create(&s).Error
	return s, err
}

func UpdateStyle(s *Style) error {
	s.Prepare()

	return database.Conn.Where("id = ?", s.ID).Updates(s).Error
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
	var i int64
	err := db().
		Model(modelStyle).
		Where("name = ? AND user_id = ? AND id != ?", s.Name, s.UserID, s.ID).
		Count(&i).
		Error

	switch {
	case i > 0:
		return errors.ErrDuplicateStyle
	case err != nil:
		return err
	default:
		return nil
	}
}

// GetStyle tries to fetch a userstyle.
func GetStyle(id string) (Style, error) {
	var s Style
	err := db().
		Select("styles.*, u.username").
		Joins("JOIN users u ON u.id = styles.user_id").
		First(&s, "styles.id = ?", id).Error
	return s, err
}

func (s Style) GetSourceCodeSize() uint64 {
	return uint64(len(s.Code))
}

func AbleToReview(uid, sid uint) (string, bool) {
	reviewSpam := new(Review)
	// Collecting of the error is not needed.
	// As we simply check "valid" data by checking if ID is a positive integer.
	if _ = reviewSpam.FindLastFromUser(sid, uid); reviewSpam.ID > 0 {
		t := time.Now().Sub(reviewSpam.CreatedAt)
		if t < 7*24*time.Hour {
			t = -7*24*time.Hour + t
			return util.RelDuration(t), false
		}
	}
	return "", true
}

func TempGetStyleByID(id int) (s *Style, err error) {
	err = database.Conn.
		Select("styles.*, u.username").
		Joins("JOIN users u ON u.id = styles.user_id").
		First(&s, "styles.id = ?", id).
		Error
	return s, err
}

// GetStyleFromAuthor tries to fetch a userstyle made by logged in user.
func GetStyleFromAuthor(id, uid int) (Style, error) {
	var s Style
	err := db().
		Select("styles.*, u.username").
		Joins("JOIN users u ON u.id = styles.user_id").
		First(&s, "styles.id = ? AND styles.user_id = ?", id, uid).Error
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
	s.Preview = fmt.Sprintf("%s/preview/%d/%dt.webp", config.App.BaseURL, s.ID, s.PreviewVersion)
}

var (
	ErrStyleAsErrs   = stderrors.New("unexpected error during validation")
	ErrStyleNoFields = stderrors.New("incorrect mandatory UserCSS fields")
	ErrStyleNoGlobal = stderrors.New("bare global styles are forbidden")
	ErrStyleNoCode   = stderrors.New("incorrect userstyle source code")
)

// Validate makes sure input data is correct.
func (s Style) Validate(v *validator.Validate, addPage bool) (map[string]any, error) {
	m := make(map[string]any)

	fields := []string{"Name", "Description", "Notes", "Category", "Code"}
	err := v.StructPartial(s, fields...)
	if err != nil {
		var errs validator.ValidationErrors
		if ok := stderrors.As(err, &errs); !ok {
			return nil, ErrStyleAsErrs
		}

		for _, e := range errs {
			switch e.Field() {
			case "Name":
				if len(s.Name) == 0 {
					m[e.Field()] = "Name is required and it cannot be empty."
				} else {
					m[e.Field()] = "Name is too long. It must be up to 50 characters."
				}
			case "Description":
				if len(s.Description) == 0 {
					m[e.Field()] = "Description is required and it cannot be empty."
				} else {
					m[e.Field()] = "Description is too long. It must be up to 160 characters."
				}
			case "Notes":
				m[e.Field()] = "Notes are too long. They must be up to 50K characters."
			case "Category":
				if len(s.Category) == 0 {
					m[e.Field()] = "Category is required and it cannot be empty."
				} else {
					m[e.Field()] = "Category is too long. It must be up to 255 characters."
				}
			case "Code":
				m[e.Field()] = "Code is too long. It must be up to 10M characters."
			}
		}
	}

	// TODO: Improve in UserCSS parser.
	var uc usercss.UserCSS
	if err := uc.Parse(s.Code); err != nil {
		msg := strings.ToUpper(string(err.Error()[0])) + err.Error()[1:] + "."
		m["Code"] = msg
		return m, ErrStyleNoCode
	}
	if errs := uc.Validate(); errs != nil {
		m["UserCSS"] = errs
		return m, ErrStyleNoFields
	}
	if addPage && len(uc.MozDocument) == 0 {
		m["Stylus"] = "Your userstyle might be affected by a bug."
		return m, ErrStyleNoGlobal
	}

	return m, err
}

// ValidateCode makes sure source code is correct.
func (s Style) ValidateCode(v *validator.Validate, addPage bool) (string, error) {
	// TODO: Improve in UserCSS parser.
	var uc usercss.UserCSS
	if err := uc.Parse(s.Code); err != nil {
		return "Error: Userstyle source code is incorrect.", ErrStyleNoCode
	}
	if errs := uc.Validate(); errs != nil {
		return "Error: Source code is missing mandatory fields: name, namespace, and/or version.", ErrStyleNoFields
	}
	if addPage && len(uc.MozDocument) == 0 {
		return "Error: Bad style format (visit https://userstyles.world/docs/faq#bad-style-format-error)", ErrStyleNoGlobal
	}

	return "", nil
}

// SetPreview will set preview image URL.
func (s *APIStyle) SetPreview() {
	s.Preview = fmt.Sprintf("%s/preview/%d/%dt.webp", config.App.BaseURL, s.ID, s.PreviewVersion)
}

// SelectUpdateStyle will update specific fields in the styles table.
func SelectUpdateStyle(s Style) error {
	s.Prepare()

	const f = "name, description, notes, code, homepage, code_size, " +
		"license, category, slug, preview, preview_version, code_checksum, " +
		"mirror_url, mirror_code, mirror_meta, import_private, mirror_private"

	return database.Conn.Select(f).Where("id = ?", s.ID).UpdateColumns(s).Error
}

func SaveStyleCode(id, s string) error {
	return os.WriteFile(filepath.Join(config.Storage.StyleDir, id), []byte(s), 0o644)
}

func RemoveStyleCode(id string) error {
	return os.Remove(filepath.Join(config.Storage.StyleDir, id))
}

// mirrorEnabled returns whether or not mirroring is enabled.
func (s Style) mirrorEnabled() bool {
	return s.MirrorCode || s.MirrorMeta
}

// sameMirrorURL returns whether or not mirror URL matches import URL.
func (s Style) sameMirrorURL() bool {
	return s.MirrorURL == "" || s.Original == s.MirrorURL
}

// isImportedAndMirrored returns whether or not a userstyle is imported and
// mirrored from the same URLs.
func (s Style) isImportedAndMirrored() bool {
	return s.isImported() && s.mirrorEnabled() && s.sameMirrorURL()
}

// ImportedAndMirrored returns from which location a userstyle is imported and
// mirrored.
func (s Style) ImportedAndMirrored() string {
	if !s.isImportedAndMirrored() {
		return ""
	}
	if s.ImportPrivate || s.MirrorPrivate {
		return "Imported and mirrored from a private source"
	}
	return "Imported and mirrored from <code>" + s.Original + "</code>"
}

// isImported returns whether or not a userstyle is isImported.
func (s Style) isImported() bool {
	return s.Original != ""
}

// Imported returns from which location a userstyle is imported.
func (s Style) Imported() string {
	if !s.isImported() {
		return ""
	}
	if s.ImportPrivate {
		return "Imported from a private source"
	}
	return "Imported from <code>" + s.Original + "</code>"
}

// isMirrored returns whether or not a userstyle is isMirrored.
func (s Style) isMirrored() bool {
	return s.MirrorURL != "" && s.mirrorEnabled()
}

// Mirrored returns from which location a userstyle is mirrored.
func (s Style) Mirrored() string {
	if !s.isMirrored() {
		return ""
	}
	if s.MirrorPrivate {
		return "Mirrored from a private source"
	}
	return "Mirrored from <code>" + s.MirrorURL + "</code>"
}
