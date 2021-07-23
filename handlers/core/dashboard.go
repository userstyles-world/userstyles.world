package core

import (
	"sort"
	"time"

	"github.com/gofiber/fiber/v2"

	"userstyles.world/handlers/jwt"
	"userstyles.world/models"
	"userstyles.world/modules/charts"
	"userstyles.world/modules/log"
)

func Dashboard(c *fiber.Ctx) error {
	u, _ := jwt.User(c)
	week := time.Now().AddDate(0, -1, 0)

	// Don't allow regular users to see this page.
	if u.Role < models.Moderator {
		return c.Render("err", fiber.Map{
			"Title": "Page not found",
			"User":  u,
		})
	}

	// Get User statistics.
	userCount, err := new(models.DashStats).GetCounts("users")
	if err != nil {
		log.Info.Println("Failed to get user statistics:", err.Error())
	}

	totalUsers := userCount[len(userCount)-1].CountSum

	t := time.Now().Format("2006-01-02")
	latestUser := userCount[len(userCount)-1]
	if t != latestUser.Date {
		latestUser = models.DashStats{
			CreatedAt: time.Now(),
			Date:      t,
			Count:     0,
			CountSum:  0,
		}
	}

	// Get Style statistics.
	styleCount, err := new(models.DashStats).GetCounts("styles")
	if err != nil {
		log.Info.Println("Failed to get style statistics:", err.Error())
	}

	totalStyles := styleCount[len(styleCount)-1].CountSum

	latestStyle := styleCount[len(styleCount)-1]
	if t != latestStyle.Date {
		latestStyle = models.DashStats{
			CreatedAt: time.Now(),
			Date:      t,
			Count:     0,
			CountSum:  0,
		}
	}

	// Get styles.
	var styles []models.StyleCard
	if c.Query("data") == "styles" {
		s, err := models.GetAllAvailableStyles()
		if err != nil {
			return c.Render("err", fiber.Map{
				"Title": "Styles not found",
				"User":  u,
			})
		}

		sort.Slice(styles, func(i, j int) bool {
			return styles[i].ID > styles[j].ID
		})

		styles = s
	}

	// Get users.
	var users []models.User
	if c.Query("data") == "users" {
		u, err := models.FindAllUsers()
		if err != nil {
			return c.Render("err", fiber.Map{
				"Title": "Users not found",
				"User":  u,
			})
		}

		sort.Slice(users, func(i, j int) bool {
			return users[i].ID > users[j].ID
		})

		users = u
	}

	// Render user history.
	var userHistory string
	if len(userCount) > 0 {
		userHistory, err = charts.GetUserHistory(userCount, week, "User history")
		if err != nil {
			log.Info.Println("Failed to render style history:", err.Error())
		}
	}

	// Render style history.
	var styleHistory string
	if len(styleCount) > 0 {
		styleHistory, err = charts.GetUserHistory(styleCount, week, "Style history")
		if err != nil {
			log.Info.Println("Failed to render user history:", err.Error())
		}
	}

	// Get history data.
	history, err := new(models.History).GetStatsForAllStyles()
	if err != nil {
		log.Info.Println("Couldn't find style histories:", err.Error())
	}

	// Render style history.
	var dailyHistory, totalHistory string
	if len(*history) > 0 {
		dailyHistory, totalHistory, err = charts.GetStyleHistory(*history)
		if err != nil {
			log.Info.Println("Failed to render style history:", err.Error())
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
		"StyleHistory": styleHistory,
	})
}
