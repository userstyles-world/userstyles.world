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
	Install bool   `gorm:"default:false"`
	View    bool   `gorm:"default:false"`
	StyleID int
	Style   Style
}

type siteStats struct {
	UserCount, StyleCount, WeeklyViews, TotalViews int64
	WeeklyInstalls, MonthlyInstalls, TotalInstalls int64
}

func generateHashedRecord(id, ip string) string {
	// Merge it here before using.
	record := ip + " " + id

	// Generate unique hash.
	h := hmac.New(sha512.New, []byte(config.STATS_KEY))
	h.Write([]byte(record))
	s := hex.EncodeToString(h.Sum(nil))

	return s
}

func AddStatsToStyle(db *gorm.DB, id, ip string, install bool) (Stats, error) {
	s := new(Stats)

	styleID, err := strconv.Atoi(id)
	if err != nil {
		return *s, err
	}

	// Set values.
	s.Hash = generateHashedRecord(id, ip)
	s.StyleID = styleID

	if install {
		s.Install = install // Initial install.
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
		s.View = true // Initial view.
		err = db.
			Debug().
			Model(s).
			Clauses(clause.OnConflict{
				Columns: []clause.Column{{Name: "hash"}},
				DoUpdates: clause.Assignments(map[string]interface{}{
					"updated_at": time.Now(),
					"view":       true,
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

func GetHomepageStatistics(db *gorm.DB) *siteStats {
	p := new(siteStats)
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
