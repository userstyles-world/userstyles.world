package models

import (
	"database/sql"
	"strconv"
	"time"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"

	"userstyles.world/modules/log"
	"userstyles.world/utils/crypto"
)

type Stats struct {
	gorm.Model
	Hash    string `gorm:"unique"`
	Style   Style
	StyleID int       `gorm:"index"`
	Install time.Time `gorm:"default:null"`
	View    time.Time `gorm:"default:null"`
}

type SiteStats struct {
	TotalUsers, TotalStyles       int64
	WeeklyViews, TotalViews       int64
	WeeklyInstalls, TotalInstalls int64
	WeeklyUpdates                 int64
}

type DashStats struct {
	CreatedAt time.Time
	Date      string
	Count     int
	CountSum  int
}

func (_ DashStats) GetCounts(t string) (q []DashStats, err error) {
	stmt := "created_at, date(created_at) Date, count(distinct id) Count,"
	stmt += "sum(count (distinct id)) over (order by date(created_at)) CountSum"

	err = db().
		Select(stmt).Table(t).Group("Date").
		Find(&q, "deleted_at is null").Error

	if err != nil {
		return nil, err
	}

	return q, nil
}

func AddStatsToStyle(id, ip string, install bool) (Stats, error) {
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

	t := time.Now()
	assignment := map[string]interface{}{
		"updated_at": t,
	}
	if install {
		s.Install = t
		assignment["install"] = t
	} else {
		s.View = t
		assignment["view"] = t
	}

	err = db().
		Model(modelStats).
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

func GetWeeklyInstallsForStyle(id string) (weekly int64) {
	lastWeek := time.Now().Add(-time.Hour * 24 * 7)
	q := "style_id = ? and install > 0 and created_at > ?"
	db().
		Model(modelStats).
		Where(q, id, lastWeek).
		Count(&weekly)

	return weekly
}

func GetTotalInstallsForStyle(id string) (total int64) {
	db().
		Model(modelStats).
		Where("style_id = ? and install > 0", id).
		Count(&total)

	return total
}

func GetTotalViewsForStyle(id string) (total int64) {
	db().
		Model(modelStats).
		Where("style_id = ? and view > 0", id).
		Count(&total)

	return total
}

func GetWeeklyUpdatesForStyle(id string) (weekly int64) {
	lastWeek := time.Now().Add(-time.Hour * 24 * 7)
	q := "style_id = ? and install > 0 and updated_at > ? and created_at < ?"
	db().
		Model(modelStats).
		Where(q, id, lastWeek, lastWeek).
		Count(&weekly)

	return weekly
}

func GetHomepageStatistics() *SiteStats {
	p := SiteStats{}
	q := `
SELECT
	(SELECT count(*) FROM users
	 WHERE users.deleted_at IS NULL) total_users,
	(SELECT count(*) FROM styles
	 WHERE styles.deleted_at IS NULL) total_styles,
	(SELECT count(*) FROM stats s
	 WHERE s.view > 0) total_views,
	(SELECT count(*) FROM stats s
	 WHERE s.install > 0) total_installs,
	(SELECT count(*) FROM stats s
	 WHERE s.view > 0 and s.created_at > @d) weekly_views,
	(SELECT count(*) FROM stats s
	 WHERE s.install > 0 and s.created_at > @d) weekly_installs,
	(SELECT count(*) FROM stats s
	 WHERE s.install > 0 and s.updated_at > @d and s.created_at < @d) weekly_updates
`

	// TODO: Replace last day with last week when we get enough data.
	lastDay := time.Now().AddDate(0, 0, -1)

	if err := db().Raw(q, sql.Named("d", lastDay)).Scan(&p).Error; err != nil {
		log.Warn.Println("Failed to get homepage stats:", err.Error())
	}

	return &p
}
