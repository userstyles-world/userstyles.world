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
	p.Rem = c % 40
	if p.Rem > 0 {
		p.Max++
	}

	// Display last page if requested page is greater than max page.
	if p.Now > p.Max {
		p.Now = p.Max
	}

	// Display first page if requested page is less than 1.
	if p.Now < 1 {
		p.Now = 1
	}

	// Set prev/next.
	p.Prev = p.Now - 1
	p.Next = p.Now + 1
}
