package api

import (
	"fmt"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/patrickmn/go-cache"

	"userstyles.world/models"
	"userstyles.world/modules/database"
	"userstyles.world/modules/errors"
	"userstyles.world/modules/log"
)

type USoFormat struct {
	Username       string `json:"an"`
	Name           string `json:"n"`
	Category       string `json:"c"`
	Screenshot     string `json:"sn"`
	UpdatedAt      int64  `json:"u"`
	TotalInstalls  int64  `json:"t"`
	WeeklyInstalls int64  `json:"w"`
	ID             uint   `json:"i"`
}

type USoStyles []USoFormat

func (s *USoStyles) Query() error {
	lastWeek := time.Now().Add(-time.Hour * 24 * 7)

	stmt := "styles.id, styles.name, styles.user_id, u.username, "
	stmt += "styles.category, strftime('%s', styles.created_at) as created_at, "
	stmt += "strftime('%s', styles.updated_at) as updated_at, "
	stmt += "printf('https://userstyles.world/api/style/preview/%d.webp', styles.id) as screenshot, "
	stmt += "(select count(*) from stats where stats.style_id = styles.id) as total_installs, "
	stmt += "(select count(*) from stats where stats.style_id = styles.id and updated_at > ? and created_at < ?) as weekly_installs"

	err := database.Conn.
		Table("styles").
		Select(stmt, lastWeek, lastWeek).
		Joins("join users u on u.id = styles.user_id").
		Find(&s).
		Error
	if err != nil {
		return errors.ErrStylesNotFound
	}

	return nil
}

var mem = cache.New(5*time.Minute, 10*time.Minute)

func convertToUSoFormat(s models.APIStyle) USoFormat {
	id := fmt.Sprintf("%d", s.ID) // Convert uint to string.

	var img string
	if s.Preview != "" {
		img = fmt.Sprintf("https://userstyles.world/api/style/preview/%d.webp", s.ID)
	}

	return USoFormat{
		ID:             s.ID,
		Name:           s.Name,
		Category:       fixCategory(s.Category),
		Username:       s.Username,
		Screenshot:     img,
		UpdatedAt:      s.UpdatedAt.Unix(),
		TotalInstalls:  models.GetTotalInstallsForStyle(id),
		WeeklyInstalls: models.GetWeeklyInstallsForStyle(id),
	}
}

func fixCategory(cat string) string {
	if cat == "unset" {
		return "global"
	}
	cat = strings.ToLower(cat)

	if strings.HasSuffix(cat, ".com") || strings.HasSuffix(cat, ".org") {
		cat = strings.TrimSuffix(cat, ".com")
		cat = strings.TrimSuffix(cat, ".org")
		// Remove any subdomain
		// web.whatsapp -> whatsapp
		if strings.Count(cat, ".") >= 1 {
			cat = strings.Split(cat, ".")[1]
		}
	} else if strings.Count(cat, ".") >= 2 {
		cat = strings.Join(strings.Split(cat, ".")[1:], ".")
	}

	return cat
}

func getUSoIndex(c *fiber.Ctx) error {
Convert:
	cached, found := mem.Get("index")
	if !found {
		styles := new(USoStyles)
		if err := styles.Query(); err != nil {
			log.Warn.Fatal("API/index/uso-format err:", err.Error())
			return c.JSON(fiber.Map{
				"data": "styles not found",
			})
		}

		// TODO: Normalize categories on add/import/edit pages.
		for _, style := range *styles {
			style.Category = fixCategory(style.Category)
		}

		mem.Set("index", styles, 10*time.Minute)
		goto Convert
	}

	return c.JSON(fiber.Map{
		"data": cached,
	})
}

func getFullIndex(c *fiber.Ctx) error {
	styles, err := models.GetAllStylesForIndexAPI()
	if err != nil {
		return c.JSON(fiber.Map{
			"data": "styles not found",
		})
	}

	return c.JSON(fiber.Map{
		"data": styles,
	})
}

func GetStyleIndex(c *fiber.Ctx) error {
	switch c.Params("format") {
	case "uso-format":
		return getUSoIndex(c)
	default:
		return getFullIndex(c)
	}
}
