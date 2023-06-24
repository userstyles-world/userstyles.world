package models

import (
	"database/sql"
	"time"

	"gorm.io/gorm"

	"userstyles.world/modules/log"
)

type Stats struct {
	ID        uint      `gorm:"primarykey"`
	CreatedAt time.Time `gorm:"index:idx_stats_weekly_installs,priority:20"`
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
	Hash      string         `gorm:"unique"`
	Style     Style
	StyleID   int       `gorm:"index:idx_stats_installed; index:idx_stats_viewed; index:idx_stats_weekly_installs"`
	Install   time.Time `gorm:"default:null; index:idx_stats_installed; index:idx_stats_weekly_installs"`
	View      time.Time `gorm:"default:null; index:idx_stats_viewed"`
}

type SiteStats struct {
	DailyUsers, WeeklyUsers, TotalUsers          int64
	DailyStyles, WeeklyStyles, TotalStyles       int64
	DailyViews, WeeklyViews, TotalViews          int64
	DailyInstalls, WeeklyInstalls, TotalInstalls int64
	DailyUpdates, WeeklyUpdates                  int64
}

type DashStats struct {
	CreatedAt time.Time
	Date      string
	Count     int
	CountSum  int
}

func (DashStats) GetCounts(t string) (q []DashStats, err error) {
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

// Delete will remove stats for a given style ID.
func (*Stats) Delete(id any) error {
	return db().Delete(&modelStats, "style_id = ?", id).Error
}

func GetHomepageStatistics() *SiteStats {
	p := SiteStats{}
	q := `
SELECT
	(SELECT count(*) FROM users u
	 WHERE u.deleted_at IS NULL AND u.created_at > @d) DailyUsers,
	(SELECT count(*) FROM users u
	 WHERE u.deleted_at IS NULL AND u.created_at > @w) WeeklyUsers,
	(SELECT count(*) FROM users u
	 WHERE u.deleted_at IS NULL) TotalUsers,

	(SELECT count(*) FROM styles s
	 WHERE s.deleted_at IS NULL AND s.created_at > @d) DailyStyles,
	(SELECT count(*) FROM styles s
	 WHERE s.deleted_at IS NULL AND s.created_at > @w) WeeklyStyles,
	(SELECT count(*) FROM styles s
	 WHERE s.deleted_at IS NULL) TotalStyles,

	(SELECT count(*) FROM stats s
	 WHERE s.deleted_at IS NULL AND s.install > 0 AND
	       s.created_at > @d) DailyInstalls,
	(SELECT count(*) FROM stats s
	 WHERE s.deleted_at IS NULL AND s.install > 0 AND
	       s.created_at > @w) WeeklyInstalls,
	(SELECT count(*) FROM stats s
	 WHERE s.deleted_at IS NULL AND s.install > 0) TotalInstalls,

	(SELECT count(*) FROM stats s
	 WHERE s.deleted_at IS NULL AND s.view > 0 AND
	       s.created_at > @d) DailyViews,
	(SELECT count(*) FROM stats s
	 WHERE s.deleted_at IS NULL AND s.view > 0 AND
	       s.created_at > @w) WeeklyViews,
	(SELECT count(*) FROM stats s
	 WHERE s.deleted_at IS NULL AND s.view > 0) TotalViews,

	(SELECT count(*) FROM stats s
	 WHERE s.deleted_at IS NULL AND s.install > 0 AND
	       s.created_at > @w) WeeklyInstalls,

	(SELECT count(*) FROM stats s
	 WHERE s.deleted_at IS NULL AND s.install > 0 AND
	       s.updated_at > @d AND
	       s.created_at < @d) DailyUpdates,
	(SELECT count(*) FROM stats s
	 WHERE s.deleted_at IS NULL AND s.install > 0 AND
	       s.updated_at > @w AND
	       s.created_at < @w) WeeklyUpdates
`

	day := sql.Named("d", time.Now().AddDate(0, 0, -1))
	week := sql.Named("w", time.Now().AddDate(0, 0, -7))
	if err := db().Raw(q, day, week).Scan(&p).Error; err != nil {
		log.Warn.Println("Failed to get homepage stats:", err.Error())
	}

	return &p
}
