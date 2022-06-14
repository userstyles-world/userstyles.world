package templates

import (
	"fmt"
	"html/template"
	"io/fs"
	"math"
	"runtime"
	"time"

	"github.com/dustin/go-humanize"
	"github.com/gofiber/template/html"

	"userstyles.world/modules/config"
	"userstyles.world/modules/log"
	"userstyles.world/modules/markdown"
	"userstyles.world/modules/util/httputil"
	"userstyles.world/utils/strutils"
)

var appConfig = map[string]string{
	"copyright":       time.Now().Format("2006"),
	"appName":         config.AppName,
	"appCodeName":     config.AppCodeName,
	"appVersion":      config.GitVersion,
	"appSourceCode":   config.AppSourceCode,
	"appLatestCommit": config.AppLatestCommit,
	"appCommitSHA":    config.AppCommitSHA,
	"allowedEmailsRe": config.AllowedEmailsRe,
	"allowedImagesRe": config.AllowedImagesRe,
}

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
	uptime := time.Since(config.AppUptime).Round(time.Second)

	return sys{
		Uptime:     uptime.String(),
		GoRoutines: runtime.NumGoroutine(),
		LastGC:     fmt.Sprintf("%.1fs", float64(time.Now().UnixNano()-int64(m.LastGC))/1000/1000/1000),
		NumGC:      int(m.NumGC),
		AverageGC:  fmt.Sprintf("%.1fs", float64(uptime.Seconds())/float64(m.NumGC)),
	}
}

func New(views fs.FS, dir string) *html.Engine {
	// Embed templates.
	fsys, err := httputil.EmbedFS(views, dir, config.Production)
	if err != nil {
		log.Warn.Fatal(err)
	}
	engine := html.NewFileSystem(fsys, ".html")

	engine.AddFunc("config", func(key string) string {
		return appConfig[key]
	})

	engine.AddFunc("sys", status)

	engine.AddFunc("comma", humanize.Comma)

	engine.AddFunc("num", func(i int64) string {
		return strutils.HumanizeNumber(int(i))
	})

	engine.AddFunc("proxy", func(src, kind string, id uint) string {
		return strutils.ProxyResources(src, kind, id)
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

	engine.AddFunc("Date", func(time time.Time) string {
		return time.Format("January 2, 2006 15:04")
	})

	engine.AddFunc("shortDate", func(time time.Time) string {
		return time.Format("2006-01-02 15:04")
	})

	engine.AddFunc("DateISO8601", func(time time.Time) string {
		return time.Format("2006-01-02T15:04:05-0700")
	})

	engine.AddFunc("add", func(a, b int) int {
		return a + b
	})

	engine.AddFunc("floor", math.Floor)

	engine.AddFunc("sub", func(a, b int) int {
		return a - b
	})

	engine.AddFunc("paginate", func(page int, sort string) template.HTML {
		s := fmt.Sprintf("/explore?page=%v", page)
		if sort != "" {
			s += fmt.Sprintf("&sort=%v", sort)
		}

		return template.HTML(strutils.QueryUnescape(s))
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
			return template.HTML(config.BaseURL)
		}
		return template.HTML(config.BaseURL + "/" + url.(string))
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

	if !config.Production {
		engine.Reload(true)
	}

	return engine
}
