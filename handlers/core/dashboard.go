package core

import (
	"fmt"
	"runtime"
	"sort"
	"time"

	"github.com/gofiber/fiber/v2"

	"userstyles.world/handlers/jwt"
	"userstyles.world/models"
	"userstyles.world/modules/charts"
	"userstyles.world/modules/config"
	"userstyles.world/modules/log"
)

var system struct {
	Uptime     string
	GoRoutines int

	MemAllocated uint64
	MemTotal     uint64
	MemSys       uint64
	Lookups      uint64
	MemMallocs   uint64
	MemFrees     uint64

	HeapAlloc    uint64
	HeapSys      uint64
	HeapInuse    uint64
	HeapIdle     uint64
	HeapReleased uint64
	HeapObjects  uint64

	StackInuse uint64
	StackSys   uint64

	MSpanSys    uint64
	MSpanInuse  uint64
	MCacheSys   uint64
	MCacheInuse uint64
	BuckHashSys uint64
	GCSys       uint64
	OtherSys    uint64

	NextGC       uint64
	LastGC       string
	PauseTotalNs string
	PauseNs      string
	NumGC        uint32
}

func getSystemStatus() {
	system.Uptime = time.Since(config.AppUptime).Round(time.Second).String()
	system.GoRoutines = runtime.NumGoroutine()

	m := new(runtime.MemStats)
	runtime.ReadMemStats(m)

	system.MemAllocated = m.Alloc
	system.MemTotal = m.TotalAlloc
	system.MemSys = m.Sys
	system.Lookups = m.Lookups
	system.MemMallocs = m.Mallocs
	system.MemFrees = m.Frees

	system.HeapAlloc = m.HeapAlloc
	system.HeapSys = m.HeapSys
	system.HeapInuse = m.HeapInuse
	system.HeapIdle = m.HeapIdle
	system.HeapReleased = m.HeapReleased
	system.HeapObjects = m.HeapObjects

	system.StackInuse = m.StackInuse
	system.StackSys = m.StackSys

	system.MSpanSys = m.MCacheInuse
	system.MSpanInuse = m.MSpanInuse
	system.MCacheSys = m.MCacheSys
	system.MCacheInuse = m.MCacheInuse
	system.BuckHashSys = m.BuckHashSys
	system.GCSys = m.GCSys
	system.OtherSys = m.OtherSys

	system.NextGC = m.NextGC
	system.LastGC = fmt.Sprintf("%.1fs", float64(time.Now().UnixNano()-int64(m.LastGC))/1000/1000/1000)
	system.PauseTotalNs = fmt.Sprintf("%.1fs", float64(m.PauseTotalNs)/1000/1000/1000)
	system.PauseNs = fmt.Sprintf("%.3fs", float64(m.PauseNs[(m.NumGC+255)%256])/1000/1000/1000)
	system.NumGC = m.NumGC
}

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

	getSystemStatus()

	// Get User statistics.
	userCount, err := new(models.DashStats).GetCounts("users")
	if err != nil {
		log.Info.Println("Failed to get user statistics:", err.Error())
	}

	t := time.Now().Format("2006-01-02")
	latestUser := userCount[len(userCount)-1]
	totalUsers := latestUser.CountSum
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

	latestStyle := styleCount[len(styleCount)-1]
	totalStyles := latestStyle.CountSum
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
		userHistory, err = charts.GetModelHistory(userCount, week, "User history")
		if err != nil {
			log.Info.Println("Failed to render style history:", err.Error())
		}
	}

	// Render style history.
	var styleHistory string
	if len(styleCount) > 0 {
		styleHistory, err = charts.GetModelHistory(styleCount, week, "Style history")
		if err != nil {
			log.Info.Println("Failed to render user history:", err.Error())
		}
	}

	// Get history data.
	history, err := new(models.History).GetStatsForAllStyles()
	if err != nil {
		log.Info.Println("Couldn't find style histories:", err.Error())
	}

	// Render stats history.
	var dailyHistory, totalHistory string
	if len(*history) > 0 {
		dailyHistory, totalHistory, err = charts.GetStatsHistory(*history)
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
		"System":       system,
	})
}
