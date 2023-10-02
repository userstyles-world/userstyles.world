package storage

import (
	"userstyles.world/modules/database"
)

// Category represents a single category.
type Category struct {
	Category string
	Count    int
}

// URL returns a link to a category page.
func (c *Category) URL() string {
	return "/category/" + c.Category
}

// Categories is an alias for a slice of Category structs.
type Categories []*Category

// GetStyleCategories returns a slice of categories for category page.
func GetStyleCategories(page, size int) (c Categories, err error) {
	offset := (page - 1) * size

	err = database.Conn.
		Select("LOWER(category) AS category, COUNT(LOWER(category)) AS count").
		Table("styles").
		Group("LOWER(category)").
		Order("count DESC").
		Offset(offset).
		Limit(size).
		Find(&c, notDeleted).
		Error
	if err != nil {
		return nil, err
	}

	return c, err
}

// CountStyleCategories returns a count of unique categories for pagination.
func CountStyleCategories() (i int, err error) {
	err = database.Conn.
		Select("COUNT(DISTINCT LOWER(category))").
		Table("styles").
		Where(notDeleted).
		Scan(&i).
		Error
	if err != nil {
		return 0, err
	}

	return i, nil
}
