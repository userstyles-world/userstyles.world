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
	StyleID int       `gorm:"index:idx_stats"`
	Install time.Time `gorm:"default:null; index:idx_stats"`
	View    time.Time `gorm:"default:null; index:idx_stats"`
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

// CreateRecord prepares style stats for upsert queries.
func (s *Stats) CreateRecord(field, id, ip string) error {
	hash, err := crypto.CreateHashedRecord(id, ip)
	if err != nil {
		return err
	}

	styleID, err := strconv.Atoi(id)
	if err != nil {
		return err
	}

	s.StyleID = styleID
	s.Hash = hash

	switch field {
	case "install":
		s.Install = time.Now()
	case "view":
		s.View = time.Now()
	}

	return nil
}

// UpsertInstall updates or inserts new install date.
func (s *Stats) UpsertInstall(id, ip string) error {
	if err := s.CreateRecord("install", id, ip); err != nil {
		return err
	}

	t := time.Now()
	s.Install = t

	return db().
		Model(modelStats).
		Clauses(clause.OnConflict{
			Columns: []clause.Column{{Name: "hash"}},
			DoUpdates: clause.Assignments(map[string]interface{}{
				"updated_at": t,
				"install":    t,
			}),
		}).
		Create(s).Error
}

// UpsertView updates or inserts new viewed date.
func (s *Stats) UpsertView(id, ip string) error {
	if err := s.CreateRecord("view", id, ip); err != nil {
		return err
	}

	t := time.Now()
	s.View = t

	return db().
		Model(modelStats).
		Clauses(clause.OnConflict{
			Columns: []clause.Column{{Name: "hash"}},
			DoUpdates: clause.Assignments(map[string]interface{}{
				"updated_at": t,
				"view":       t,
			}),
		}).
		Create(s).Error
}

// Delete will remove stats for a given style ID.
func (*Stats) Delete(id interface{}) error {
	return db().Delete(&modelStats, "style_id = ?", id).Error
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

func GetWeeklyViewsForStyle(id string) (weekly int64) {
	lastWeek := time.Now().Add(-time.Hour * 24 * 7)
	q := "style_id = ? and view > 0 and created_at > ?"
	db().Model(modelStats).Where(q, id, lastWeek).Count(&weekly)

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
