package models

import (
	"time"

	"gorm.io/gorm"
)

type LogKind = uint8

const (
	LogBanUser LogKind = iota + 1
	LogRemoveStyle
	LogRemoveReview
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

type APILog struct {
	ID             uint
	CreatedAt      time.Time
	UserID         uint
	Username       string
	Reason         string
	Censor         bool
	Kind           LogKind
	TargetData     string
	TargetUserName string
}

// CreateLog inserts a new log entry into the database.
func CreateLog(db *gorm.DB, log *Log) (err error) {
	return db.Model(modelLog).Create(log).Error
}

// GetModLogs tries to return all moderation logs.
func GetModLogs(db *gorm.DB) (l []APILog, err error) {
	err = db.
		Model(modelLog).
		Select("logs.*, (SELECT username FROM users WHERE id = logs.user_id) AS Username").
		Order("created_at DESC").
		Find(&l).
		Error
	return
}
