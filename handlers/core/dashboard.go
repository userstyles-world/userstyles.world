package core

import (
	"log"
	"sort"
	"time"

	"github.com/gofiber/fiber/v2"

	"userstyles.world/handlers/jwt"
	"userstyles.world/models"
	"userstyles.world/modules/charts"
	"userstyles.world/modules/database"
)

type count struct {
	CreatedAt time.Time
	Date      string
	Count     int
}

func getCount(t string) (c []count, err error) {
	err = database.Conn.Debug().
		Select("created_at, date(created_at) Date, count(distinct id) Count").
		Table(t).Group("Date").Find(&c, "deleted_at is null").Error

	if err != nil {
		return nil, err
	}

	return c, nil
}

func Dashboard(c *fiber.Ctx) error {
	u, _ := jwt.User(c)

	// Don't allow regular users to see this page.
	if u.Role < models.Moderator {
		return c.Render("err", fiber.Map{
			"Title": "Page not found",
			"User":  u,
		})
	}

	// Get User statistics.
	userCount, err := getCount("users")
	if err != nil {
		log.Printf("Failed to get user statistics, err: %v\n", err)
	}

	var totalUsers int
	for _, v := range userCount {
		totalUsers += v.Count
	}

	t := time.Now().Format("2006-01-02")
	latestUser := userCount[len(userCount)-1]
	if t != latestUser.Date {
		latestUser = count{
			CreatedAt: time.Now(),
			Date:      t,
			Count:     0,
		}
	}

	// Get Style statistics.
	styleCount, err := getCount("styles")
	if err != nil {
		log.Printf("Failed to get style statistics, err: %v\n", err)
	}

	var totalStyles int
	for _, v := range styleCount {
		totalStyles += v.Count
	}

	latestStyle := styleCount[len(styleCount)-1]
	if t != latestStyle.Date {
		latestStyle = count{
			CreatedAt: time.Now(),
			Date:      t,
			Count:     0,
		}
	}

	// TODO: Refactor.
	// Get styles.
	styles, err := models.GetAllAvailableStyles()
	if err != nil {
		return c.Render("err", fiber.Map{
			"Title": "Styles not found",
			"User":  u,
		})
	}

	sort.Slice(styles, func(i, j int) bool {
		return styles[i].ID > styles[j].ID
	})

	// Get users.
	users, err := models.FindAllUsers()
	if err != nil {
		return c.Render("err", fiber.Map{
			"Title": "Users not found",
			"User":  u,
		})
	}

	// Render user history.
	var userHistory string
	if len(users) > 0 {
		userHistory, err = charts.GetUserHistory(users)
		if err != nil {
			log.Printf("Failed to render user history, err: %s\n", err.Error())
		}
	}

	sort.Slice(users, func(i, j int) bool {
		return users[i].ID > users[j].ID
	})

	// Get history data.
	history, err := new(models.History).GetStatsForAllStyles()
	if err != nil {
		log.Printf("Couldn't find style histories, err: %s", err.Error())
	}

	// Render style history.
	var dailyHistory, totalHistory string
	if len(*history) > 0 {
		dailyHistory, totalHistory, err = charts.GetStyleHistory(*history)
		if err != nil {
			log.Printf("Failed to render history for all styles, err: %s\n", err.Error())
		}
	}

	return c.Render("core/dashboard", fiber.Map{
		"Title":        "Dashboard",
		"User":         u,
		"TotalStyles":  totalStyles,
		"LatestStyle":  latestStyle,
		"TotalUsers":   totalUsers,
		"LatestUser":   latestUser,
		"Styles":       styles,
		"Users":        users,
		"DailyHistory": dailyHistory,
		"TotalHistory": totalHistory,
		"UserHistory":  userHistory,
	})
}
