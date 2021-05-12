package models

import (
	"strconv"
	"time"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"

	"userstyles.world/utils/crypto"
)

const (
	totalViews     = "view = 1"
	weeklyViews    = "view = 1 and created_at > ?"
	totalInstalls  = "install = 1"
	weeklyInstalls = "install = 1 and created_at > ?"
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
	UserCount, StyleCount         int64
	WeeklyViews, TotalViews       int64
	WeeklyInstalls, TotalInstalls int64
}

func AddStatsToStyle(db *gorm.DB, id, ip string, install bool) (Stats, error) {
	s := new(Stats)

	styleID, err := strconv.Atoi(id)
	if err != nil {
		return *s, err
	}

	// Set values.
	s.StyleID = styleID
	s.Hash, err = crypto.CreateHashedRecord(id, ip)
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
	lastWeek := time.Now().Add(-time.Hour * 24 * 7)
	q := "style_id = ? and install = 1 and created_at > ?"
	db.
		Model(Stats{}).
		Where(q, id, lastWeek).
		Count(&weekly)

	return weekly
}

func GetTotalInstallsForStyle(db *gorm.DB, id string) (total int64) {
	db.
		Model(Stats{}).
		Where("style_id = ? and install = 1", id).
		Count(&total)

	return total
}

func GetTotalViewsForStyle(db *gorm.DB, id string) (total int64) {
	db.
		Model(Stats{}).
		Where("style_id = ? and view = 1", id).
		Count(&total)

	return total
}

func GetHomepageStatistics(db *gorm.DB) *SiteStats {
	p, s := SiteStats{}, Stats{}

	// TODO: Replace last day with last week when we get enough data.
	lastDay := time.Now().AddDate(0, 0, -1)

	db.Model(User{}).Where("id").Count(&p.UserCount)
	db.Model(Style{}).Where("id").Count(&p.StyleCount)
	db.Model(&s).Where(totalViews).Count(&p.TotalViews)
	db.Model(&s).Where(weeklyViews, lastDay).Count(&p.WeeklyViews)
	db.Model(&s).Where(weeklyInstalls, lastDay).Count(&p.WeeklyInstalls)
	db.Model(&s).Where(totalInstalls).Count(&p.TotalInstalls)

	return &p
}
