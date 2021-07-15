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

type stats struct {
	NewUsers  int
	NewStyles int
}

func Dashboard(c *fiber.Ctx) error {
	u, _ := jwt.User(c)
	t := time.Now().Format("2006-01-02")
	s := stats{}

	// Don't allow regular users to see this page.
	if u.Role < models.Moderator {
		return c.Render("err", fiber.Map{
			"Title": "Page not found",
			"User":  u,
		})
	}

	// Get styles.
	styles, err := models.GetAllAvailableStyles()
	if err != nil {
		return c.Render("err", fiber.Map{
			"Title": "Styles not found",
			"User":  u,
		})
	}

	// Summary of new styles.
	for _, v := range styles {
		if v.CreatedAt.Format("2006-01-02") == t {
			s.NewStyles++
		}
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

	// Summary of new users.
	for _, v := range users {
		if v.CreatedAt.Format("2006-01-02") == t {
			s.NewUsers++
		}
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
		"Styles":       styles,
		"Users":        users,
		"Summary":      s,
		"DailyHistory": dailyHistory,
		"TotalHistory": totalHistory,
		"UserHistory":  userHistory,
	})
}
