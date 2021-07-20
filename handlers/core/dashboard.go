package core

import (
	"log"
	"sort"
	"time"

	"github.com/gofiber/fiber/v2"

	"userstyles.world/handlers/jwt"
	"userstyles.world/models"
	"userstyles.world/modules/charts"
)

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
	userCount, err := new(models.DashStats).GetCounts("users")
	if err != nil {
		log.Printf("Failed to get user statistics, err: %v\n", err)
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
		log.Printf("Failed to get style statistics, err: %v\n", err)
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

	// TODO: Refactor.
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
	var userHistory string
	var users []models.User
	if c.Query("data") == "users" {
		u, err := models.FindAllUsers()
		if err != nil {
			return c.Render("err", fiber.Map{
				"Title": "Users not found",
				"User":  u,
			})
		}

		// Render user history.
		if len(users) > 0 {
			userHistory, err = charts.GetUserHistory(users)
			if err != nil {
				log.Printf("Failed to render user history, err: %s\n", err.Error())
			}
		}

		sort.Slice(users, func(i, j int) bool {
			return users[i].ID > users[j].ID
		})

		users = u
	}

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
