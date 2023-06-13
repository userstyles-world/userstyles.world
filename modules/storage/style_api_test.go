package storage

import (
	"strconv"
	"testing"
	"time"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

	"userstyles.world/models"
)

func initDB() (*gorm.DB, error) {
	db, err := gorm.Open(sqlite.Open(":memory:"))
	if err != nil {
		return nil, err
	}

	t := []any{models.Style{}, models.Stats{}, models.User{}, models.Review{}}
	if err = db.AutoMigrate(t...); err != nil {
		return nil, err
	}

	return db, nil
}

func TestGetStyleCompactIndex(t *testing.T) {
	db, err := initDB()
	if err != nil {
		t.Fatal(err)
	}

	var s []models.Style
	for i := 1; i <= 2; i++ {
		s = append(s, models.Style{
			Model: gorm.Model{
				UpdatedAt: time.Date(1970, 1, 1, 1, 0, 0, 0, time.UTC),
			},
			Name: "test " + strconv.Itoa(i),
		})
	}

	if err = db.Create(s).Error; err != nil {
		t.Fatal(err)
	}

	got, err := GetStyleCompactIndex(db)
	if err != nil {
		t.Fatal(err)
	}

	exp := `{"data":[{"n":"test 1","an":"","sn":"","c":"","i":1,"u":3600,"t":0,"w":0,"r":0},{"n":"test 2","an":"","sn":"","c":"","i":2,"u":3600,"t":0,"w":0,"r":0}]}`
	if string(got) != exp {
		t.Errorf("got: %s\n", got)
		t.Errorf("exp: %s\n", exp)
	}
}
