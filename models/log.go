package models

import (
	"time"

	"gorm.io/gorm"

	"userstyles.world/modules/errors"
)

type LogKind = uint8

const (
	LogBanUser LogKind = iota + 1
	LogRemoveStyle
	LogRemoveReview
)

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

// AddLog adds a new log to the database.
func (*Log) AddLog(logEntry *Log) (err error) {
	err = db().
		Model(modelLog).
		Create(logEntry).
		Error
	if err != nil {
		return errors.ErrFailedLogAddition
	}
	return nil
}

// GetLogOfKind returns all the logs of the specified kind and
// select the correct user Author.
func GetLogOfKind(kind LogKind) ([]APILog, error) {
	var q []APILog

	err := db().
		Model(modelLog).
		Select("logs.*, (SELECT username FROM users WHERE id = logs.user_id) AS Username").
		Where("kind = ?", kind).
		Order("created_at desc").
		Find(&q).
		Error
	if err != nil {
		return q, errors.ErrFailedLogRetrieval
	}
	return q, nil
}
