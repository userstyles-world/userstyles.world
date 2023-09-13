package storage

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"userstyles.world/modules/database"
	"userstyles.world/modules/util"
)

const (
	selectAuthor      = "(SELECT username FROM users WHERE user_id = users.id AND deleted_at IS NULL) AS Username"
	selectInstalls    = "(SELECT total_installs FROM histories h WHERE h.style_id = styles.id ORDER BY id DESC LIMIT 1) AS installs"
	selectViews       = "(SELECT total_views FROM histories h WHERE h.style_id = styles.id ORDER BY id DESC LIMIT 1) AS views"
	selectRatings     = "(SELECT ROUND(AVG(rating), 1) FROM reviews r WHERE r.style_id = styles.id AND r.rating > 0 AND r.deleted_at IS NULL) AS Rating"
	selectReviewCount = "(SELECT COUNT(rating) FROM reviews r WHERE r.style_id = styles.id AND r.rating > 0 AND r.deleted_at IS NULL) AS ReviewCount"
	notDeleted        = "deleted_at IS NULL"
)

var (
	selectCards = strings.Join([]string{
		"id", "updated_at", "name", "preview",
		selectAuthor, selectInstalls, selectViews, selectRatings, selectReviewCount,
	}, ", ")
	selectSearchCards = strings.Join([]string{
		"id", "created_at", "updated_at", "name", "preview",
		selectAuthor, selectInstalls, selectViews, selectRatings, selectReviewCount,
	}, ", ")
)

// StyleCard is a field-aligned struct optimized for style cards.
type StyleCard struct {
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	Preview     string    `json:"preview"`
	Username    string    `json:"username"`
	Name        string    `json:"name"`
	ID          int       `json:"id"`
	Views       int64     `json:"views"`
	Installs    int64     `json:"installs"`
	Rating      float64   `json:"rating"`
	ReviewCount int       `json:"reviewcount"`
}

// TableName returns which table in database to use with GORM.
func (StyleCard) TableName() string { return "styles" }

// Slug returns a user- and SEO-friendly URL.
func (x StyleCard) Slug() string {
	return util.Slug(x.Name)
}

// StyleURL returns an absolute path to a style.
func (x StyleCard) StyleURL() string {
	return fmt.Sprintf("/style/%d/%s", x.ID, x.Slug())
}

// FindStyleCardsForSearch returns style cards for search page.
func FindStyleCardsForSearch(items []int, kind string, size int) ([]StyleCard, error) {
	var b strings.Builder
	b.WriteString("SELECT " + selectSearchCards + " ")
	b.WriteString("FROM styles WHERE id in (")
	for i, item := range items {
		if i == 0 {
			b.WriteString(strconv.Itoa(item))
		} else {
			b.WriteString(", " + strconv.Itoa(item))
		}
	}

	if kind == "" || kind == "default" {
		// NOTE: This is a dynamic/custom ordering implementation, because there's
		// no other way [that I know of] to return results in the order they were
		// selected.  We might need to decrease the amount (96 ATM) of results that
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
		b.WriteString("END")
	} else {
		// TODO: Refactor [probably] using Pagination struct.
		switch kind {
		case "newest":
			kind = "styles.created_at DESC"
		case "oldest":
			kind = "styles.created_at ASC"
		case "recentlyupdated":
			kind = "styles.updated_at DESC"
		case "leastupdated":
			kind = "styles.updated_at ASC"
		case "mostinstalls":
			kind = "installs DESC"
		case "leastinstalls":
			kind = "installs ASC"
		case "mostviews":
			kind = "views DESC"
		case "leastviews":
			kind = "views ASC"
		default:
			kind = "styles.id ASC"
		}

		b.WriteString(") ORDER BY " + kind + " LIMIT " + strconv.Itoa(size))
	}

	res := make([]StyleCard, 0, len(items))
	if err := database.Conn.Raw(b.String()).Scan(&res).Error; err != nil {
		return nil, err
	}

	return res, nil
}

// FindStyleCardsForUsername returns style cards for a user.
func FindStyleCardsForUsername(username string) ([]StyleCard, error) {
	var res []StyleCard

	err := database.Conn.
		Select(selectCards).
		Find(&res, "deleted_at IS NULL AND username = ?", username).Error
	if err != nil {
		return nil, err
	}

	return res, nil
}

// FindStyleCardsFeatured returns style cars that are featured.
func FindStyleCardsFeatured() ([]StyleCard, error) {
	var res []StyleCard

	err := database.Conn.
		Select(selectCards).
		Find(&res, "deleted_at IS NULL AND featured = 1").Error
	if err != nil {
		return nil, err
	}

	return res, nil
}

// FindStyleCardsPaginated returns style cards for paginated pages.
func FindStyleCardsPaginated(page, size int, order string) ([]StyleCard, error) {
	var res []StyleCard
	offset := (page - 1) * size

	var stmt string
	switch {
	case strings.HasPrefix(order, "styles"):
		stmt = "id"
	case strings.HasPrefix(order, "views"):
		stmt = "id, (SELECT total_views FROM histories h WHERE h.style_id = styles.id ORDER BY id DESC LIMIT 1) AS views"
	case strings.HasPrefix(order, "installs"):
		stmt = "id, (SELECT total_installs FROM histories h WHERE h.style_id = styles.id ORDER BY id DESC LIMIT 1) AS installs"
	case strings.HasPrefix(order, "rating"):
		stmt = "id, (SELECT ROUND(AVG(rating), 1) FROM reviews r WHERE r.style_id = styles.id) AS rating"
	}

	var nums []struct{ ID int }
	err := database.Conn.
		Select(stmt).Table("styles").
		Order(order).Offset(offset).Limit(size).Find(&nums, notDeleted).Error
	if err != nil {
		return nil, err
	}

	items := make([]int, 0, len(nums))
	for _, num := range nums {
		items = append(items, num.ID)
	}

	err = database.Conn.
		Select(selectCards).Order(order).Find(&res, items).Error
	if err != nil {
		return nil, err
	}

	return res, nil
}

// FindStyleCardsPaginatedForUserID returns user's style cards for paginated pages.
func FindStyleCardsPaginatedForUserID(page, size int, order string, id uint) ([]StyleCard, error) {
	var stmt string
	switch {
	case strings.HasPrefix(order, "styles"):
		stmt = "id"
	case strings.HasPrefix(order, "views"):
		stmt = "id, (SELECT total_views FROM histories h WHERE h.style_id = styles.id ORDER BY id DESC LIMIT 1) AS views"
	case strings.HasPrefix(order, "installs"):
		stmt = "id, (SELECT total_installs FROM histories h WHERE h.style_id = styles.id ORDER BY id DESC LIMIT 1) AS installs"
	case strings.HasPrefix(order, "rating"):
		stmt = "id, (SELECT ROUND(AVG(rating), 1) FROM reviews r WHERE r.style_id = styles.id) AS rating"
	}

	offset := (page - 1) * size
	var nums []struct{ ID int }
	err := database.Conn.
		Select(stmt).Table("styles").
		Order(order).Offset(offset).Limit(size).
		Find(&nums, notDeleted+" AND user_id = ?", id).Error
	if err != nil {
		return nil, err
	}

	items := make([]int, 0, len(nums))
	for _, num := range nums {
		items = append(items, num.ID)
	}

	var res []StyleCard
	err = database.Conn.
		Select(selectCards).Order(order).
		Find(&res, "id in ?", items).Error
	if err != nil {
		return nil, err
	}

	return res, nil
}

// FindStyleCardsCreatedOn returns style cards created on a specific date.
func FindStyleCardsCreatedOn(date time.Time) ([]StyleCard, error) {
	var res []StyleCard

	cond := notDeleted + " AND DATE(?) == DATE(created_at)"
	err := database.Conn.
		Select(selectCards).
		Find(&res, cond, date).Error
	if err != nil {
		return nil, err
	}

	return res, nil
}
