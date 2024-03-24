package storage

import (
	"strconv"
	"testing"
	"time"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

	"userstyles.world/models"
	"userstyles.world/modules/config"
)

func initDB() (*gorm.DB, error) {
	db, err := gorm.Open(sqlite.Open(":memory:"))
	if err != nil {
		return nil, err
	}

	t := []any{models.Style{}, models.Stats{}, models.User{}, models.Review{}, models.History{}}
	if err = db.AutoMigrate(t...); err != nil {
		return nil, err
	}

	return db, nil
}

func TestGetStyleCompactIndex(t *testing.T) {
	cases := []struct {
		name string
		size int
		exp  string
	}{
		{"one", 1, `{"data":[{"n":"test 1","an":"","sn":"","c":"","i":1,"u":3600,"t":0,"w":0,"r":0}]}`},
		{"many", 2, `{"data":[{"n":"test 1","an":"","sn":"","c":"","i":1,"u":3600,"t":0,"w":0,"r":0},{"n":"test 2","an":"","sn":"","c":"","i":2,"u":3600,"t":0,"w":0,"r":0}]}`},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			db, err := initDB()
			if err != nil {
				t.Fatal(err)
			}

			var s []models.Style
			for i := 1; i <= c.size; i++ {
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

			if string(got) != c.exp {
				t.Errorf("got: %s\n", got)
				t.Errorf("exp: %s\n", c.exp)
			}
		})
	}
}

func BenchmarkGetStyleCompactIndex(b *testing.B) {
	cases := []struct {
		name string
		size int
	}{
		{"10", 10},
		{"100", 100},
		{"1000", 1000},
		{"10000", 10000},
	}

	for _, c := range cases {
		b.Run(c.name, func(b *testing.B) {
			db, err := initDB()
			if err != nil {
				b.Fatal(err)
			}

			var s []models.Style
			for i := 1; i <= c.size; i++ {
				id := strconv.Itoa(i)
				s = append(s, models.Style{
					Model: gorm.Model{
						UpdatedAt: time.Date(1970, 1, 1, 1, 0, 0, 0, time.UTC),
					},
					Name:    "test " + id,
					Preview: config.Config.BaseURL + "/preview/" + id + "/0.webp",
				})
			}

			if err = db.CreateInBatches(s, 100).Error; err != nil {
				b.Fatal(err)
			}

			for i := 0; i < b.N; i++ {
				GetStyleCompactIndex(db)
			}
		})
	}
}
