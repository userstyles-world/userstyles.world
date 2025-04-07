package models

import (
	"fmt"

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

// Log struct has all the relavant information for a log entry.
type Log struct {
	gorm.Model
	UserID   uint
	Username string
	Reason   string
	Message  string

	// This isn't the Censor you'd expect.
	// It will only just wrap the style's information into a spoiler.
	// This will be used for disturbing names.
	Censor         bool
	Kind           LogKind
	TargetData     string
	TargetUserName string
}

// Permalink returns a link to the mod log page.
func (l Log) Permalink() string {
	return fmt.Sprintf("/modlog/%d", l.ID)
}

// CreateLog inserts a new log entry into the database.
func CreateLog(db *gorm.DB, log *Log) (err error) {
	return db.Model(modelLog).Create(log).Error
}

// GetModLogs tries to return all moderation logs.
func GetModLogs(db *gorm.DB, page, size, order int) (l []Log, err error) {
	tx := db.
		Model(modelLog).
		Select("logs.*, (SELECT username FROM users WHERE id = logs.user_id) AS Username").
		Order("created_at DESC").
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
	err = db.
		Model(modelLog).
		Select("logs.*, (SELECT username FROM users WHERE id = logs.user_id) AS Username").
		Find(&l, "id = ?", id).Error

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
