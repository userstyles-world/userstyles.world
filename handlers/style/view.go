package style

import (
	"bytes"
	"fmt"
	"log"

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

	// TODO: Get data from database.

	// Visualize data.
	var b = 1000.
	graph := chart.Chart{
		Width: 1248,
		XAxis: chart.XAxis{
			ValueFormatter: chart.TimeDateValueFormatter,
		},
		YAxis: chart.YAxis{
			Name: "Installs/views",
		},
		Series: []chart.Series{
			chart.ContinuousSeries{
				Name:    "Time Series",
				XValues: []float64{10 * b, 20 * b, 30 * b, 40 * b, 50 * b, 60 * b, 70 * b, 80 * b},
				YValues: []float64{1.0, 2.0, 30.0, 4.0, 50.0, 6.0, 7.0, 88.0},
			},
			chart.ContinuousSeries{
				Style: chart.Style{
					StrokeColor: chart.GetDefaultColor(1),
				},
				XValues: []float64{10 * b, 20 * b, 30 * b, 40 * b, 50 * b, 60 * b, 70 * b, 80 * b},
				YValues: []float64{15.0, 52.0, 30.0, 42.0, 50.0, 26.0, 77.0, 38.0},
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
