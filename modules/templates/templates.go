package templates

import (
	"embed"
	"fmt"
	"html/template"
	"net/http"
	"runtime"
	"time"

	"github.com/dustin/go-humanize"
	"github.com/gofiber/template/html"
	"github.com/microcosm-cc/bluemonday"
	md "github.com/russross/blackfriday/v2"

	"userstyles.world/modules/config"
	"userstyles.world/utils/strutils"
)

var ext = md.CommonExtensions | md.AutoHeadingIDs

var appConfig = map[string]interface{}{
	"copyright":       time.Now().Year(),
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

func New(views embed.FS) *html.Engine {
	// Embed templates.
	var fsys http.FileSystem
	if config.Production {
		fsys = http.FS(views)
	} else {
		fsys = http.Dir("views")
	}
	engine := html.NewFileSystem(fsys, ".html")

	engine.AddFunc("config", func(key string) template.HTML {
		return template.HTML(fmt.Sprintf("%v", appConfig[key]))
	})

	engine.AddFunc("sys", status)

	engine.AddFunc("comma", humanize.Comma)

	engine.AddFunc("num", func(i int64) string {
		return strutils.HumanizeNumber(int(i))
	})

	engine.AddFunc("proxy", func(s, t string, id uint) template.HTML {
		return template.HTML(strutils.ProxyResources(s, t, id))
	})

	engine.AddFunc("MarkdownSafe", func(s string) template.HTML {
		gen := md.Run([]byte(s), md.WithExtensions(ext))
		return template.HTML(gen)
	})

	engine.AddFunc("MarkdownUnsafe", func(s string) string {
		// Generate Markdown then sanitize it before returning HTML.
		gen := md.Run(
			[]byte(s),
			md.WithExtensions(md.HardLineBreak),
		)
		out := bluemonday.UGCPolicy().SanitizeBytes(gen)

		return string(out)
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
		return time.Format("2006-02-01 15:04")
	})

	engine.AddFunc("DateISO8601", func(time time.Time) string {
		return time.Format("2006-02-01T15:04:05-0700")
	})

	engine.AddFunc("add", func(a, b int) int {
		return a + b
	})

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

	engine.AddFunc("canonical", func(url interface{}) template.HTML {
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

	if !config.Production {
		engine.Reload(true)
	}

	return engine
}
