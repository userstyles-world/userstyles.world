package core

import (
	"fmt"
	"runtime"
	"sync"
	"time"

	"github.com/dustin/go-humanize"
	"github.com/gofiber/fiber/v2"

	"userstyles.world/handlers/jwt"
	"userstyles.world/models"
	"userstyles.world/modules/charts"
	"userstyles.world/modules/config"
	"userstyles.world/modules/log"
	"userstyles.world/modules/storage"
)

var (
	systemMutex sync.Mutex

	system struct {
		Uptime     string
		GoRoutines int

		MemAllocated string
		MemTotal     string
		MemSys       string
		Lookups      string
		MemMallocs   string
		MemFrees     string

		HeapAlloc    string
		HeapSys      string
		HeapInuse    string
		HeapIdle     string
		HeapReleased string
		HeapObjects  string

		StackInuse  string
		StackSys    string
		MSpanSys    string
		MSpanInuse  string
		MCacheSys   string
		MCacheInuse string
		BuckHashSys string
		GCSys       string
		OtherSys    string

		NextGC       string
		AverageGC    string
		LastGC       string
		PauseTotalNs string
		PauseNs      string
		NumGC        string
	}
)

func getSystemStatus() {
	systemMutex.Lock()
	defer systemMutex.Unlock()

	uptime := time.Since(config.App.Started).Round(time.Second)

	system.Uptime = uptime.String()
	system.GoRoutines = runtime.NumGoroutine()

	m := new(runtime.MemStats)
	runtime.ReadMemStats(m)

	system.MemAllocated = humanize.Bytes(m.Alloc)
	system.MemTotal = humanize.Bytes(m.TotalAlloc)
	system.MemSys = humanize.Bytes(m.Sys)
	system.Lookups = humanize.Comma(int64(m.Lookups))
	system.MemMallocs = humanize.Comma(int64(m.Mallocs))
	system.MemFrees = humanize.Comma(int64(m.Frees))

	system.HeapAlloc = humanize.Bytes(m.HeapAlloc)
	system.HeapSys = humanize.Bytes(m.HeapSys)
	system.HeapInuse = humanize.Bytes(m.HeapInuse)
	system.HeapIdle = humanize.Bytes(m.HeapIdle)
	system.HeapReleased = humanize.Bytes(m.HeapReleased)
	system.HeapObjects = humanize.Comma(int64(m.HeapObjects))

	system.StackInuse = humanize.Bytes(m.StackInuse)
	system.StackSys = humanize.Bytes(m.StackSys)

	system.MSpanSys = humanize.Bytes(m.MSpanSys)
	system.MSpanInuse = humanize.Bytes(m.MSpanInuse)
	system.MCacheSys = humanize.Bytes(m.MCacheSys)
	system.MCacheInuse = humanize.Bytes(m.MCacheInuse)
	system.BuckHashSys = humanize.Bytes(m.BuckHashSys)
	system.GCSys = humanize.Bytes(m.GCSys)
	system.OtherSys = humanize.Bytes(m.OtherSys)

	system.NextGC = humanize.Bytes(m.NextGC)
	system.AverageGC = fmt.Sprintf("%.1fs", float64(uptime.Seconds())/float64(m.NumGC))
	system.LastGC = fmt.Sprintf("%.1fs", float64(time.Now().UnixNano()-int64(m.LastGC))/1000/1000/1000)
	system.PauseTotalNs = fmt.Sprintf("%.1fs", float64(m.PauseTotalNs)/1000/1000/1000)
	system.PauseNs = fmt.Sprintf("%.3fs", float64(m.PauseNs[(m.NumGC+255)%256])/1000/1000/1000)
	system.NumGC = humanize.Comma(int64(m.NumGC))
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

	// Get styles.
	if c.Query("data") == "styles" {
		styles, err := storage.FindStyleCardsCreatedOn(time.Now())
		if err != nil {
			return c.Render("err", fiber.Map{
				"Title": "Styles not found",
				"User":  u,
			})
		}

		return c.Render("core/dashboard", fiber.Map{
			"Title":        "Dashboard",
			"User":         u,
			"Styles":       styles,
			"RenderStyles": true,
		})
	}

	// Get users.
	if c.Query("data") == "users" {
		users, err := storage.FindUsersCreatedOn(time.Now())
		if err != nil {
			return c.Render("err", fiber.Map{
				"Title": "Users not found",
				"User":  u,
			})
		}

		return c.Render("core/dashboard", fiber.Map{
			"Title":       "Dashboard",
			"User":        u,
			"Users":       users,
			"RenderUsers": true,
		})
	}

	// Get System statistics.
	getSystemStatus()

	// Get User statistics.
	dashUsers := new(models.DashStats)
	userCount, err := dashUsers.GetCounts("users")
	if err != nil {
		log.Info.Println("Failed to get user statistics:", err.Error())
	}

	t := time.Now().Format("2006-01-02")
	// userCount = []models.DashStats{}
	var latestUser models.DashStats
	var totalUsers int
	if len(userCount) > 0 {
		latestUser = userCount[len(userCount)-1]
		totalUsers = latestUser.CountSum
		if t != latestUser.Date {
			latestUser = models.DashStats{
				CreatedAt: time.Now(),
				Date:      t,
				Count:     0,
				CountSum:  0,
			}
		}
	}

	// Get Style statistics.
	dashStyles := new(models.DashStats)
	styleCount, err := dashStyles.GetCounts("styles")
	if err != nil {
		log.Info.Println("Failed to get style statistics:", err.Error())
	}

	var latestStyle models.DashStats
	var totalStyles int
	if len(styleCount) > 0 {
		latestStyle = styleCount[len(styleCount)-1]
		totalStyles = latestStyle.CountSum
		if t != latestStyle.Date {
			latestStyle = models.DashStats{
				CreatedAt: time.Now(),
				Date:      t,
				Count:     0,
				CountSum:  0,
			}
		}
	}

	// Render user history.
	var userHistory string
	if len(userCount) > 1 {
		userHistory, err = charts.GetModelHistory(userCount, week, "User history")
		if err != nil {
			log.Info.Println("Failed to render style history:", err.Error())
		}
	}

	// Render style history.
	var styleHistory string
	if len(styleCount) > 1 {
		styleHistory, err = charts.GetModelHistory(styleCount, week, "Style history")
		if err != nil {
			log.Info.Println("Failed to render user history:", err.Error())
		}
	}

	// Get history data.
	history, err := models.GetAllStyleHistories()
	if err != nil {
		log.Info.Println("Couldn't find style histories:", err.Error())
	}

	// Render stats history.
	var dailyHistory, totalHistory string
	if len(history) > 2 {
		dailyHistory, totalHistory, err = charts.GetStatsHistory(history)
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
		"DailyHistory": dailyHistory,
		"TotalHistory": totalHistory,
		"UserHistory":  userHistory,
		"StyleHistory": styleHistory,
		"System":       system,
	})
}
