package models

import (
	"fmt"
	"net/url"
	"strconv"

	"userstyles.world/modules/config"
)

// Pagination is a field-aligned struct optimized for pagination.
type Pagination struct {
	Path     string
	Sort     string
	Query    string
	Category string

	Prev3 int
	Prev2 int
	Prev1 int
	Now   int
	Next1 int
	Next2 int
	Next3 int
	Max   int
}

// NewPagination is a convenience function that initializes pagination struct.
func NewPagination(page, total int, sort, path string) Pagination {
	p := Pagination{
		Now:  page,
		Path: path,
		Sort: sort,
	}

	p.calcItems(total)

	return p
}

// URL generates a dynamic path from available items.
func (p Pagination) URL(page int) string {
	s := fmt.Sprintf("%s?page=%d", p.Path, page)

	if p.Sort != "" {
		s += fmt.Sprintf("&sort=%s", p.Sort)
	}

	if p.Query != "" {
		s += fmt.Sprintf("&q=%s", url.QueryEscape(p.Query))
	}

	if p.Category != "" {
		s += fmt.Sprintf("&category=%s", url.QueryEscape(p.Category))
	}

	return s
}

// IsValidPage checks whether a passed parameter is a valid number.
func IsValidPage(s string) (int, error) {
	if s == "" {
		return 1, nil
	}

	i, err := strconv.Atoi(s)
	if err != nil {
		return 0, err
	}

	return i, err
}

// calcItems calculates values for pagination fields.
func (p *Pagination) calcItems(total int) {
	if p.Now < 1 {
		p.Now = 1
	}

	if total == 0 {
		p.Max = 1
		return
	}

	// Calculate max page and remainder.
	p.Max = total / config.AppPageMaxItems
	if total%config.AppPageMaxItems > 0 {
		p.Max++
	}

	// Set prev/next.
	p.Prev1 = p.Now - 1
	p.Prev2 = p.Now - 2
	p.Prev3 = p.Now - 3
	p.Next1 = p.Now + 1
	p.Next2 = p.Now + 2
	p.Next3 = p.Now + 3
}

// OutOfBounds checks whether current page exceeds the bounds.
func (p *Pagination) OutOfBounds() bool {
	// Display last page if requested page is greater than max page.
	if p.Now > p.Max {
		p.Now = p.Max
		return true
	}

	// Display first page if requested page is less than 1.
	if p.Now < 1 {
		p.Now = 1
		return true
	}

	return false
}

// Show renders pagination when page count is more than one.
func (p Pagination) Show() bool {
	return p.Max > 1
}

func (p *Pagination) SortStyles() string {
	switch p.Sort {
	case "newest":
		return "styles.created_at DESC"
	case "oldest":
		return "styles.created_at ASC"
	case "recentlyupdated":
		return "styles.updated_at DESC"
	case "leastupdated":
		return "styles.updated_at ASC"
	case "mostinstalls":
		return "installs DESC"
	case "leastinstalls":
		return "installs ASC"
	case "mostviews":
		return "views DESC"
	case "leastviews":
		return "views ASC"
	case "ratinghigh":
		return "rating DESC"
	case "ratinglow":
		return "rating ASC"
	default:
		return "styles.id ASC"
	}
}
