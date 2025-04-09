package models

import (
	"fmt"
	"time"

	"gorm.io/gorm"
)

type LogKind int

const (
	LogBanUser LogKind = iota + 1
	LogRemoveStyle
	LogRemoveReview
	LogCount
)

func (x LogKind) String() string {
	return []string{"Remove user", "Remove style", "Remove review"}[x-1]
}

// Log represents schema for logs table.
type Log struct {
	ID        uint      `gorm:"column:id; primaryKey"`
	ByUserID  uint      `gorm:"column:by_user_id"`
	ToUserID  uint      `gorm:"column:to_user_id"`
	StyleID   uint      `gorm:"column:style_id; default:null"`
	ReviewID  uint      `gorm:"column:review_id; default:null"`
	Kind      LogKind   `gorm:"column:kind"`
	CreatedAt time.Time `gorm:"column:created_at"`
	UpdatedAt time.Time `gorm:"column:updated_at"`
	DeletedAt deletedAt `gorm:"column:deleted_at; index"`
	Reason    string    `gorm:"column:reason"`
	Message   string    `gorm:"column:message"`
	// This isn't the Censor you'd expect.
	// It will only just wrap the style's information into a spoiler.
	// This will be used for disturbing names.
	Censor bool `gorm:"column:censor"`
	ByUser *User
	ToUser *User
	Style  *Style
	Review *Review
}

// TableName returns a table to be used with GORM.
func (Log) TableName() string { return "logs" }

// Permalink returns a link to the mod log page.
func (l Log) Permalink() string {
	return fmt.Sprintf("/modlog/%d", l.ID)
}

// ShowMessage returns whether or not to display a private message if it's set.
func (l Log) ShowMessage(u APIUser) bool {
	return l.Message != "" && (l.ToUserID == u.ID || u.IsModOrAdmin())
}

// CreateLog inserts a new log entry into the database.
func CreateLog(db *gorm.DB, log *Log) (err error) {
	return db.Model(modelLog).Create(log).Error
}

// GetModLogs tries to return all moderation logs.
func GetModLogs(db *gorm.DB, page, size, order int) (l []Log, err error) {
	const q = "logs.*, ByUser.username AS ByUser__username, ToUser.username AS ToUser__username"
	tx := db.
		Model(modelLog).
		Joins("LEFT JOIN users ByUser ON logs.by_user_id = ByUser.id").
		Joins("LEFT JOIN users ToUser ON logs.to_user_id = ToUser.id").
		Select(q).
		Order("id DESC").
		Offset((page - 1) * size).
		Limit(size)

	if order > 0 {
		tx.Where("kind = ?", order)
	}

	err = tx.Find(&l).Error

	return l, err
}

// GetModLog tries to return a specific moderation log.
func GetModLog(db *gorm.DB, id int) (l Log, err error) {
	q := "logs.*, ByUser.username AS ByUser__username, "
	q += "ToUser.username AS ToUser__username, styles.name AS Style__name"

	err = db.
		Select(q).
		Model(modelLog).
		Joins("LEFT JOIN users ByUser ON logs.by_user_id = ByUser.id").
		Joins("LEFT JOIN users ToUser ON logs.to_user_id = ToUser.id").
		Joins("LEFT JOIN styles ON logs.style_id = styles.id").
		Find(&l, "logs.id = ?", id).Error

	return l, err
}

// GetModLogCount tries to count mod log entries depending on their kind.
func GetModLogCount(db *gorm.DB, kind int) (i int64, err error) {
	tx := db.Model(modelLog)
	if kind > 0 {
		tx.Where("kind = ?", kind)
	}

	err = tx.Count(&i).Error

	return i, err
}
