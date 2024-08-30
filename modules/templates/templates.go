package templates

import (
	"encoding/json"
	"fmt"
	"html/template"
	"math"
	"net/http"
	"runtime"
	"time"

	"github.com/dustin/go-humanize"
	"github.com/gofiber/template/html"

	"userstyles.world/modules/config"
	"userstyles.world/modules/markdown"
	"userstyles.world/modules/util"
)

type sys struct {
	Uptime     string
	GoRoutines int
	LastGC     string
	NumGC      int
	AverageGC  string
}

func status() sys {
	m := new(runtime.MemStats)
	runtime.ReadMemStats(m)
	uptime := time.Since(config.App.Started).Round(time.Second)

	return sys{
		Uptime:     uptime.String(),
		GoRoutines: runtime.NumGoroutine(),
		LastGC:     fmt.Sprintf("%.1fs", float64(time.Now().UnixNano()-int64(m.LastGC))/1000/1000/1000),
		NumGC:      int(m.NumGC),
		AverageGC:  fmt.Sprintf("%.1fs", float64(uptime.Seconds())/float64(m.NumGC)),
	}
}

func New(views http.FileSystem) *html.Engine {
	engine := html.NewFileSystem(views, ".tmpl")

	engine.AddFunc("sys", status)

	engine.AddFunc("comma", humanize.Comma)

	engine.AddFunc("size", humanize.Bytes)

	engine.AddFunc("num", util.RelNumber)

	engine.AddFunc("proxy", func(src, kind string, id uint) string {
		return util.ProxyResources(src, kind, id)
	})

	engine.AddFunc("MarkdownSafe", func(text string) string {
		return markdown.RenderSafe([]byte(text))
	})

	engine.AddFunc("MarkdownUnsafe", func(text string) string {
		return markdown.RenderUnsafe([]byte(text))
	})

	engine.AddFunc("descMax", func(s string) string {
		if len(s) > 160 {
			return s[:160] + "…"
		}

		return s
	})

	engine.AddFunc("rel", util.RelTime)

	engine.AddFunc("iso", func(t time.Time) string {
		return t.Format(time.RFC3339)
	})

	engine.AddFunc("toTime", func(s string) time.Time {
		t, err := time.Parse(time.RFC3339, s)
		if err != nil {
			return time.Time{}
		}
		return t
	})

	engine.AddFunc("add", func(a, b int) int {
		return a + b
	})

	engine.AddFunc("floor", math.Floor)

	engine.AddFunc("sub", func(a, b int) int {
		return a - b
	})

	engine.AddFunc("unescape", func(s string) template.HTML {
		return template.HTML(s)
	})

	engine.AddFunc("BaseCodeTemplate", func() template.HTML {
		return `/* ==UserStyle==
@name           Example
@namespace      example.com
@version        1.0.0
@description    A new userstyle
@author         Me
==/UserStyle== */

@-moz-document domain("example.com") {
    /**
        Your code goes here!
        More info in the official documentation at Stylus' wiki:
        https://github.com/openstyles/stylus/wiki/Writing-UserCSS
    */
}`
	})

	engine.AddFunc("canonical", func(url any) template.HTML {
		if url == nil {
			return template.HTML(config.App.BaseURL)
		}
		return template.HTML(config.App.BaseURL + "/" + url.(string))
	})

	engine.AddFunc("Elapsed", func(dur time.Duration) template.HTML {
		// Normalize duration.
		dur = dur.Round(time.Microsecond)
		return template.HTML(dur.String())
	})

	engine.AddFunc("fullImage", func(url string) string {
		return url[:len(url)-6] + url[len(url)-5:]
	})

	engine.AddFunc("toJPEG", func(url string) string {
		return url[:len(url)-4] + "jpeg"
	})

	engine.AddFunc("jsonify", func(a any) string {
		b, err := json.Marshal(a)
		if err != nil {
			return "Failed to serialize data to JSON."
		}
		return string(b)
	})

	if !config.App.Production {
		engine.Reload(true)
	}

	return engine
}
