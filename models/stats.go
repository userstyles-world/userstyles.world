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

var (
	lastDay   = time.Now().AddDate(0, 0, -1)
	lastWeek  = time.Now().AddDate(0, 0, -7)
	lastMonth = time.Now().AddDate(0, -1, 0)
)

type Stats struct {
	gorm.Model
	Hash    string `gorm:"unique"`
	Install bool   `gorm:"default:false"`
	StyleID int
	Style   Style
}

type siteStats struct {
	UserCount, StyleCount, WeeklyViews, TotalViews int64
	WeeklyInstalls, MonthlyInstalls, TotalInstalls int64
}

func AddStatsToStyle(db *gorm.DB, id, ip string, install bool) (Stats, error) {
	s := new(Stats)

	styleID, err := strconv.Atoi(id)
	if err != nil {
		return *s, err
	}

	// TODO: Refactor as GenerateHashedRecord; we have a circular dependency now.
	record := ip + " " + id

	// Generate unique hash.
	h := hmac.New(sha512.New, []byte(config.STATS_KEY))
	h.Write([]byte(record))
	sha := hex.EncodeToString(h.Sum(nil))

	// Set values.
	s.Hash = sha
	s.StyleID = styleID

	if install {
		err = db.
			Debug().
			Model(s).
			Clauses(clause.OnConflict{
				Columns: []clause.Column{{Name: "hash"}},
				DoUpdates: clause.Assignments(map[string]interface{}{
					"updated_at": time.Now(),
					"install":    true,
				}),
			}).
			Create(s).
			Error
	} else {
		err = db.
			Debug().
			Model(s).
			Clauses(clause.OnConflict{
				Columns: []clause.Column{{Name: "hash"}},
				DoUpdates: clause.Assignments(map[string]interface{}{
					"updated_at": time.Now(),
				}),
			}).
			Create(s).
			Error
	}

	if err != nil {
		return *s, err
	}

	return *s, nil
}

func GetWeeklyInstallsForStyle(db *gorm.DB, id string) int64 {
	var weekly int64
	db.
		Model(Stats{}).
		Where("style_id = ? and install = ? and updated_at > ?", id, true, lastWeek).
		Count(&weekly)

	return weekly
}

func GetTotalInstallsForStyle(db *gorm.DB, id string) int64 {
	var total int64
	db.
		Model(Stats{}).
		Where("style_id = ? and install = ?", id, true).
		Count(&total)

	return total
}

func GetTotalViewsForStyle(db *gorm.DB, id string) int64 {
	var total int64
	db.
		Model(Stats{}).
		Where("style_id = ?", id).
		Count(&total)

	return total
}

func GetHomepageStatistics(db *gorm.DB) *siteStats {
	p := new(siteStats)
	i := "install = ?"
	t := "install = ? and updated_at > ?"
	// m := "install = ? and updated_at > ?"

	db.Model(User{}).Where("id").Count(&p.UserCount)
	db.Model(Style{}).Where("id").Count(&p.StyleCount)
	db.Model(Stats{}).Where(i, false).Count(&p.TotalViews)
	// TODO: Replace lastDay with lastWeek on 2021-04-21.
	db.Model(Stats{}).Where(t, false, lastDay).Count(&p.WeeklyViews)
	db.Model(Stats{}).Where(t, true, lastDay).Count(&p.WeeklyInstalls)
	// db.Model(Stats{}).Where(m, true, lastMonth).Count(&p.MonthlyInstalls)
	db.Model(Stats{}).Where(i, true).Count(&p.TotalInstalls)

	return p
}
