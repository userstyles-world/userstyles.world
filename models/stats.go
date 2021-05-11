package models

import (
	"crypto/hmac"
	"crypto/sha512"
	"encoding/hex"
	"strconv"
	"time"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"

	"userstyles.world/config"
)

type Stats struct {
	gorm.Model
	Hash    string `gorm:"unique"`
	Style   Style
	StyleID int
	Install bool `gorm:"default:false"`
	View    bool `gorm:"default:false"`
}

type SiteStats struct {
	UserCount, StyleCount, WeeklyViews, TotalViews int64
	WeeklyInstalls, MonthlyInstalls, TotalInstalls int64
}

func generateHashedRecord(id, ip string) (string, error) {
	// Merge it here before using.
	record := ip + " " + id

	// Generate unique hash.
	h := hmac.New(sha512.New, []byte(config.STATS_KEY))
	if _, err := h.Write([]byte(record)); err != nil {
		return "", err
	}
	s := hex.EncodeToString(h.Sum(nil))

	return s, nil
}

func AddStatsToStyle(db *gorm.DB, id, ip string, install bool) (Stats, error) {
	s := new(Stats)

	styleID, err := strconv.Atoi(id)
	if err != nil {
		return *s, err
	}

	// Set values.
	s.Hash, err = generateHashedRecord(id, ip)
	s.StyleID = styleID
	if err != nil {
		return *s, err
	}

	assignment := map[string]interface{}{
		"updated_at": time.Now(),
	}
	if install {
		s.Install = true
		assignment["install"] = true
	} else {
		s.View = true
		assignment["view"] = true
	}

	err = db.
		Debug().
		Model(s).
		Clauses(clause.OnConflict{
			Columns:   []clause.Column{{Name: "hash"}},
			DoUpdates: clause.Assignments(assignment),
		}).
		Create(s).
		Error
	if err != nil {
		return *s, err
	}

	return *s, nil
}

func GetWeeklyInstallsForStyle(db *gorm.DB, id string) (weekly int64) {
	lastWeek := time.Now().AddDate(0, -1, 0)
	q := "style_id = ? and install = ? and updated_at > ?"
	db.
		Model(Stats{}).
		Where(q, id, true, lastWeek).
		Count(&weekly)

	return weekly
}

func GetTotalInstallsForStyle(db *gorm.DB, id string) (total int64) {
	db.
		Model(Stats{}).
		Where("style_id = ? and install = ?", id, true).
		Count(&total)

	return total
}

func GetTotalViewsForStyle(db *gorm.DB, id string) (total int64) {
	db.
		Model(Stats{}).
		Where("style_id = ?", id).
		Count(&total)

	return total
}

func GetHomepageStatistics(db *gorm.DB) *SiteStats {
	p := new(SiteStats)
	i, t := "install = ?", "install = ? and updated_at > ?"

	// TODO: Replace last day with last week when we get enough data.
	lastDay := time.Now().AddDate(0, 0, -1)

	db.Model(User{}).Where("id").Count(&p.UserCount)
	db.Model(Style{}).Where("id").Count(&p.StyleCount)
	db.Model(Stats{}).Where(i, false).Count(&p.TotalViews)
	db.Debug().Model(Stats{}).Where(t, false, lastDay).Count(&p.WeeklyViews)
	db.Debug().Model(Stats{}).Where(t, true, lastDay).Count(&p.WeeklyInstalls)
	// db.Model(Stats{}).Where(t, true, lastMonth).Count(&p.MonthlyInstalls)
	db.Model(Stats{}).Where(i, true).Count(&p.TotalInstalls)

	return p
}
