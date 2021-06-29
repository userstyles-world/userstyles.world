package core

import (
	"bytes"
	"log"
	"sort"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/userstyles-world/go-chart/v2"

	"userstyles.world/handlers/jwt"
	"userstyles.world/models"
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

	sort.Slice(users, func(i, j int) bool {
		return users[i].ID > users[j].ID
	})

	// Get history.
	history, err := new(models.History).GetStatsForAllStyles()
	if err != nil {
		return c.Render("err", fiber.Map{
			"Title": "Users not found",
			"User":  u,
		})
	}

	var dates []time.Time
	var views []float64
	var updates []float64
	var installs []float64
	for _, v := range *history {
		// fmt.Printf("%v | %3v | %3v | %3v\n", v.CreatedAt.Format("2006-01-02"),
		// 	v.DailyInstalls, v.DailyUpdates, v.DailyViews)
		dates = append(dates, v.CreatedAt)
		views = append(views, float64(v.DailyViews))
		updates = append(updates, float64(v.DailyUpdates))
		installs = append(installs, float64(v.DailyInstalls))
	}

	// Visualize data.
	graph := chart.Chart{
		Width:      1248,
		Canvas:     chart.Style{ClassName: "bg inner"},
		Background: chart.Style{ClassName: "bg outer"},
		XAxis:      chart.XAxis{Name: "Date"},
		YAxis:      chart.YAxis{Name: "Count"},
		Series: []chart.Series{
			chart.TimeSeries{
				Name:    "Installs",
				XValues: dates,
				YValues: installs,
			},
			chart.TimeSeries{
				Name:    "Updates",
				XValues: dates,
				YValues: updates,
			},
			chart.TimeSeries{
				Name:    "Views",
				XValues: dates,
				YValues: views,
			},
		},
	}
	graph.Elements = []chart.Renderable{chart.Legend(&graph)}

	buffer := bytes.NewBuffer([]byte{})
	err = graph.Render(chart.SVG, buffer)
	if err != nil && buffer.Len() != 220 {
		log.Printf("Failed to render SVG, err: %s\n", err.Error())
		return c.Render("err", fiber.Map{
			"Title": "Internal server error",
			"User":  u,
		})
	}

	return c.Render("core/dashboard", fiber.Map{
		"Title":   "Dashboard",
		"User":    u,
		"Styles":  styles,
		"Users":   users,
		"History": buffer.String(),
	})
}
