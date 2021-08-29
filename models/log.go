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
)

// Log struct has all the relavant information for a log entry.
type Log struct {
	gorm.Model
	UserID   uint
	Username string
	Reason   string

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
func GetLogOfKind(kind LogKind) (q *[]APILog, err error) {
	err = db().
		Model(modelLog).
		Select("logs.*, u.username").
		Joins("join users u on u.id = logs.user_id").
		Where("kind = ?", kind).
		Find(&q).
		Error
	if err != nil {
		return nil, errors.ErrFailedLogRetrieval
	}
	return q, nil
}
