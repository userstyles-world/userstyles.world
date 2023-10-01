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
func GetStyleCategories() (c Categories, err error) {
	err = database.Conn.
		Select("category, COUNT(category) AS count").
		Table("styles").
		Group("category").
		Order("count DESC").
		Find(&c, notDeleted).
		Error
	if err != nil {
		return nil, err
	}

	return c, err
}
