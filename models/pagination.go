package models

import (
	"strconv"
)

type Pagination struct {
	Min  int
	Prev int
	Now  int
	Next int
	Max  int
	Rem  int
	Sort string
}

func (p *Pagination) ConvPage(s string) error {
	if s != "" {
		i, err := strconv.Atoi(s)
		if err != nil {
			return err
		}
		p.Now = i
	} else {
		p.Now = 1
	}

	return nil
}

func (p *Pagination) CalcItems(c, i int) {
	p.Min = 1

	// Calculate max page and remainder.
	p.Max = c / i
	p.Rem = c % i
	if p.Rem > 0 {
		p.Max++
	}

	// Set prev/next.
	p.Prev = p.Now - 1
	p.Next = p.Now + 1
}

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
	default:
		return "styles.id ASC"
	}
}
