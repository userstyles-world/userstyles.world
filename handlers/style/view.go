package style

import (
	"bytes"
	"fmt"
	"log"
	"time"

	"github.com/gofiber/fiber/v2"
	chart "github.com/wcharczuk/go-chart/v2"

	"userstyles.world/handlers/jwt"
	"userstyles.world/models"
	"userstyles.world/utils/strings"
)

func GetStylePage(c *fiber.Ctx) error {
	u, _ := jwt.User(c)
	id, name := c.Params("id"), c.Params("name")

	// Check if style exists.
	data, err := models.GetStyleByID(id)
	if err != nil {
		return c.Render("err", fiber.Map{
			"Title": "Style not found",
			"User":  u,
		})
	}

	// Create slugged URL.
	// TODO: Refactor after GetStyleByID switches away from APIStyle.
	slug := strings.SlugifyURL(data.Name)

	// Always redirect to correct slugged URL.
	if name != slug {
		url := fmt.Sprintf("/style/%s/%s", id, slug)
		return c.Redirect(url, fiber.StatusSeeOther)
	}

	// Count views.
	_, err = models.AddStatsToStyle(id, c.IP(), false)
	if err != nil {
		log.Println("Failed to add stats to style, err:", err)
		return c.Render("err", fiber.Map{
			"Title": "Internal server error",
			"User":  u,
		})
	}

	// Get history data.
	history, err := new(models.History).GetStatsForStyle(id)
	if err != nil {
		log.Printf("No style stats for style %s, err: %s", id, err.Error())
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
		Width: 1248,
		XAxis: chart.XAxis{Name: "Date"},
		YAxis: chart.YAxis{Name: "Count"},
		Series: []chart.Series{
			chart.TimeSeries{
				XValues: dates,
				YValues: views,
			},
			chart.TimeSeries{
				Style:   chart.Style{StrokeColor: chart.ColorRed},
				XValues: dates,
				YValues: updates,
			},
			chart.TimeSeries{
				Style:   chart.Style{StrokeColor: chart.ColorGreen},
				XValues: dates,
				YValues: installs,
			},
		},
	}

	buffer := bytes.NewBuffer([]byte{})
	err = graph.Render(chart.SVG, buffer)
	if err != nil {
		log.Printf("Failed to render SVG, err: %s\n", err.Error())
		return c.Render("err", fiber.Map{
			"Title": "Internal server error",
			"User":  u,
		})
	}

	return c.Render("style/view", fiber.Map{
		"Title":          data.Name,
		"User":           u,
		"Style":          data,
		"TotalViews":     models.GetTotalViewsForStyle(id),
		"TotalInstalls":  models.GetTotalInstallsForStyle(id),
		"WeeklyInstalls": models.GetWeeklyInstallsForStyle(id),
		"WeeklyUpdates":  models.GetWeeklyUpdatesForStyle(id),
		"Url":            fmt.Sprintf("https://userstyles.world/style/%d", data.ID),
		"Slug":           c.Path(),
		"History":        buffer,
	})
}
