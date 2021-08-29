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
func (s *Stats) CreateRecord(id, ip string) error {
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

	return nil
}

// UpsertInstall updates or inserts style install date.
func (s *Stats) UpsertInstall() error {
	t := time.Now()
	s.Install = t

	if err := db().
		Debug().
		Model(modelStats).
		Clauses(clause.OnConflict{
			Columns: []clause.Column{{Name: "hash"}},
			DoUpdates: clause.Assignments(map[string]interface{}{
				"updated_at": t,
				"install":    t,
			}),
		}).
		Create(s).Error; err != nil {
		return err
	}

	return nil
}

// UpsertView updates or inserts style view date.
func (s *Stats) UpsertView() error {
	t := time.Now()
	s.View = t

	if err := db().
		Model(modelStats).
		Clauses(clause.OnConflict{
			Columns: []clause.Column{{Name: "hash"}},
			DoUpdates: clause.Assignments(map[string]interface{}{
				"updated_at": t,
				"view":       t,
			}),
		}).
		Create(s).Error; err != nil {
		return err
	}

	return nil
}

// Delete will remove stats for a given style ID.
func (s *Stats) Delete(id interface{}) error {
	return db().Debug().Delete(&modelStats, "style_id = ?", id).Error
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
	(SELECT count(*) FROM users
	 WHERE users.deleted_at IS NULL) total_users,
	(SELECT count(*) FROM styles
	 WHERE styles.deleted_at IS NULL) total_styles,
	(SELECT count(*) FROM stats s
	 WHERE s.deleted_at IS NULL AND s.view > 0) total_views,
	(SELECT count(*) FROM stats s
	 WHERE s.deleted_at IS NULL AND s.install > 0) total_installs,
	(SELECT count(*) FROM stats s
	 WHERE s.deleted_at IS NULL AND s.view > 0 AND s.created_at > @d) weekly_views,
	(SELECT count(*) FROM stats s
	 WHERE s.deleted_at IS NULL AND s.install > 0 AND s.created_at > @d) weekly_installs,
	(SELECT count(*) FROM stats s
	 WHERE s.deleted_at IS NULL AND s.install > 0 AND s.updated_at > @d AND s.created_at < @d) weekly_updates
`

	// TODO: Replace last day with last week when we get enough data.
	lastDay := time.Now().AddDate(0, 0, -1)

	if err := db().Raw(q, sql.Named("d", lastDay)).Scan(&p).Error; err != nil {
		log.Warn.Println("Failed to get homepage stats:", err.Error())
	}

	return &p
}
