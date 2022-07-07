package storage

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"userstyles.world/modules/config"
	"userstyles.world/modules/database"
	"userstyles.world/utils/strutils"
)

var (
	selectAuthor   = "(SELECT username FROM users WHERE user_id = users.id AND deleted_at IS NULL) AS Username"
	selectInstalls = "(SELECT COUNT(*) FROM stats s WHERE s.style_id = styles.id AND s.install > 0) AS Installs"
	selectViews    = "(SELECT COUNT(*) FROM stats s WHERE s.style_id = styles.id AND s.view > 0) AS Views"
	selectRatings  = "(SELECT ROUND(AVG(rating), 1) FROM reviews r WHERE r.style_id = styles.id AND r.deleted_at IS NULL) AS Rating"
)

// StyleCard is a field-aligned struct optimized for style cards.
type StyleCard struct {
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Preview   string    `json:"preview"`
	Username  string    `json:"username"`
	Name      string    `json:"name"`
	ID        int       `json:"id"`
	Views     int64     `json:"views"`
	Installs  int64     `json:"installs"`
	Rating    float64   `json:"rating"`
}

// TableName returns which table in database to use with GORM.
func (StyleCard) TableName() string { return "styles" }

// Slug returns a URL- and SEO-friendly string.
func (x StyleCard) Slug() string {
	return strutils.SlugifyURL(x.Name)
}

// StyleURL returns an absolute path to a style.
func (x StyleCard) StyleURL() string {
	return fmt.Sprintf("/style/%d/%s", x.ID, x.Slug())
}

// FindStyleCardsForSearch returns style cards for search page.
func FindStyleCardsForSearch(items []int) ([]StyleCard, error) {
	fields := []string{
		"id", "created_at", "updated_at", "name", "preview",
		selectAuthor, selectInstalls, selectViews, selectRatings,
	}
	var b strings.Builder
	b.WriteString("SELECT " + strings.Join(fields, ", ") + " ")
	b.WriteString("FROM styles WHERE id in (")
	for i, item := range items {
		if i == 0 {
			b.WriteString(strconv.Itoa(item))
		} else {
			b.WriteString(", " + strconv.Itoa(item))
		}
	}

	// NOTE: This is a dynamic/custom ordering implementation, because there's
	// no other way [that I know of] to return results in the order they were
	// selected.  We might need to decrease the amount (99 ATM) of results that
	// we return, because it could be too much for Pluto (our VPS) to process.
	//
	// We want to keep "ordering by relevance" by default, which is returned by
	// our search engine, and we'll use sort package for ordering in other ways
	// for the time being.  In the future, especially if we consider adding
	// pagination to results, we might want to do it all in here.
	b.WriteString(") ORDER BY CASE id ")
	for i, num := range items {
		b.WriteString("WHEN " + strconv.Itoa(num) + " THEN " + strconv.Itoa(i) + " ")
	}
	b.WriteString("END;")

	var res []StyleCard
	if err := database.Conn.Raw(b.String()).Scan(&res).Error; err != nil {
		return nil, err
	}

	return res, nil
}

// FindStyleCardsForUsername returns style cards for a user.
func FindStyleCardsForUsername(username string) ([]StyleCard, error) {
	var res []StyleCard

	fields := []string{
		"id", "updated_at", "name", "preview",
		selectAuthor, selectInstalls, selectViews, selectRatings,
	}
	err := database.Conn.
		Select(strings.Join(fields, ", ")).
		Find(&res, "deleted_at IS NULL AND username = ?", username).Error
	if err != nil {
		return nil, err
	}

	return res, nil
}

// FindStyleCardsFeatured returns style cars that are featured.
func FindStyleCardsFeatured() ([]StyleCard, error) {
	var res []StyleCard

	fields := []string{
		"id", "updated_at", "name", "preview",
		selectAuthor, selectInstalls, selectViews, selectRatings,
	}
	err := database.Conn.
		Select(strings.Join(fields, ", ")).
		Find(&res, "deleted_at IS NULL AND featured = 1").Error
	if err != nil {
		return nil, err
	}

	return res, nil
}

// FindStyleCardsPaginated returns style cards for paginated pages.
func FindStyleCardsPaginated(page int, order string) ([]StyleCard, error) {
	var res []StyleCard

	size := config.AppPageMaxItems
	offset := (page - 1) * size

	// Reflection go brrrr.
	nums := []struct {
		ID, Views, Installs int
	}{}

	var stmt string
	switch {
	case strings.HasPrefix(order, "styles"):
		stmt = "id, created_at, updated_at"
	case strings.HasPrefix(order, "views"):
		stmt = "id, (SELECT SUM(daily_views) FROM histories h WHERE h.style_id = styles.id) AS views"
	case strings.HasPrefix(order, "installs"):
		stmt = "id, (SELECT SUM(daily_installs) FROM histories h WHERE h.style_id = styles.id) AS installs"
	}

	err := database.Conn.
		Select(stmt).Table("styles").
		Order(order).Offset(offset).Limit(size).Find(&nums).Error
	if err != nil {
		return nil, err
	}

	styleIDs := make([]int, 0, len(nums))
	for _, partial := range nums {
		styleIDs = append(styleIDs, partial.ID)
	}

	fields := []string{
		"id", "updated_at", "name", "preview",
		selectAuthor, selectInstalls, selectViews, selectRatings,
	}
	err = database.Conn.
		Select(strings.Join(fields, ", ")).Order(order).Find(&res, styleIDs).Error
	if err != nil {
		return nil, err
	}

	return res, nil
}
