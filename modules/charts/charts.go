package charts

import (
	"errors"
	"time"

	"github.com/userstyles-world/go-chart/v2"
	"github.com/valyala/bytebufferpool"

	"userstyles.world/models"
)

func GetStatsHistory(history []models.History) (dailyStats string, totalStats string, err error) {
	historyLen := len(history)
	dates := make([]time.Time, 0, historyLen)
	dailyViews := make([]float64, 0, historyLen)
	dailyUpdates := make([]float64, 0, historyLen)
	dailyInstalls := make([]float64, 0, historyLen)
	totalViews := make([]float64, 0, historyLen)
	totalInstalls := make([]float64, 0, historyLen)

	for _, v := range history {
		dates = append(dates, v.CreatedAt)
		dailyViews = append(dailyViews, float64(v.DailyViews))
		dailyUpdates = append(dailyUpdates, float64(v.DailyUpdates))
		dailyInstalls = append(dailyInstalls, float64(v.DailyInstalls))
		totalViews = append(totalViews, float64(v.TotalViews))
		totalInstalls = append(totalInstalls, float64(v.TotalInstalls))
	}

	// Visualize daily stats.
	dailyGraph := chart.Chart{
		Width:      1248,
		Canvas:     chart.Style{ClassName: "bg inner"},
		Background: chart.Style{ClassName: "bg outer"},
		XAxis:      chart.XAxis{Name: "Date"},
		YAxis:      chart.YAxis{Name: "Daily count"},
		Series: []chart.Series{
			chart.TimeSeries{
				Name:    "Daily installs",
				XValues: dates,
				YValues: dailyInstalls,
			},
			chart.TimeSeries{
				Name:    "Daily updates",
				XValues: dates,
				YValues: dailyUpdates,
			},
			chart.TimeSeries{
				Name:    "Daily views",
				XValues: dates,
				YValues: dailyViews,
			},
		},
	}
	dailyGraph.Elements = []chart.Renderable{chart.Legend(&dailyGraph)}

	daily := bytebufferpool.Get()
	defer bytebufferpool.Put(daily)
	dailyFailed := daily.Len() != 220
	if err := dailyGraph.Render(chart.SVG, daily); err != nil && dailyFailed {
		return "", "", err
	}

	// Visualize total stats.
	totalGraph := chart.Chart{
		Width:      1248,
		Canvas:     chart.Style{ClassName: "bg inner"},
		Background: chart.Style{ClassName: "bg outer"},
		XAxis:      chart.XAxis{Name: "Date"},
		YAxis:      chart.YAxis{Name: "Total count"},
		Series: []chart.Series{
			chart.TimeSeries{
				Name:    "Total installs",
				XValues: dates,
				YValues: totalInstalls,
			},
			chart.TimeSeries{
				Name:    "Total views",
				XValues: dates,
				YValues: totalViews,
			},
		},
	}
	totalGraph.Elements = []chart.Renderable{chart.Legend(&totalGraph)}

	total := bytebufferpool.Get()
	defer bytebufferpool.Put(total)
	totalFailed := total.Len() != 220
	if err := totalGraph.Render(chart.SVG, total); err != nil && totalFailed {
		return "", "", err
	}

	return daily.String(), total.String(), nil
}

func GetModelHistory(vals []models.DashStats, t time.Time, title string) (string, error) {
	bars := []chart.Value{}
	for _, val := range vals {
		if val.CreatedAt.After(t) {
			bars = append(bars, chart.Value{
				Label: val.Date,
				Value: float64(val.Count),
			})
		}
	}

	if len(bars) < 1 {
		return "", errors.New("please provide at least one bar")
	}

	usersGraph := chart.BarChart{
		Title: title,
		Background: chart.Style{
			Padding: chart.Box{Top: 40},
		},
		Height: 360,
		Bars:   bars,
		XAxis: chart.Style{
			TextRotationDegrees: 90.0,
		},
	}

	b := bytebufferpool.Get()
	defer bytebufferpool.Put(b)
	notEnoughData := b.Len() != 220
	if err := usersGraph.Render(chart.SVG, b); err != nil && notEnoughData {
		return "", err
	}

	return b.String(), nil
}
