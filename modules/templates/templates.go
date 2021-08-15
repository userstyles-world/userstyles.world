package templates

import (
	"fmt"
	"html/template"
	"strconv"
	"time"

	"github.com/gofiber/template/html"
	"github.com/markbates/pkger"
	"github.com/microcosm-cc/bluemonday"
	md "github.com/russross/blackfriday/v2"

	"userstyles.world/modules/config"
	"userstyles.world/utils/strings"
)

var ext = md.CommonExtensions | md.AutoHeadingIDs

var appConfig = map[string]interface{}{
	"copyright":       time.Now().Year(),
	"appName":         config.AppName,
	"appCodeName":     config.AppCodeName,
	"appVersion":      config.GitVersion,
	"appSourceCode":   config.AppSourceCode,
	"appLatestCommit": config.AppLatestCommit,
	"allowedEmailsRe": `^[a-zA-Z0-9.!#$%&’*+/=?^_\x60{|}~-]+@[a-zA-Z0-9-]+(?:\.[a-zA-Z0-9-]+)+$`,
	"allowedImagesRe": `^https:\/\/(www\.)?(userstyles.world/api|(user-images|raw)\.githubusercontent\.com|github\.com|gitlab\.com|codeberg\.org|cdn\.jsdelivr\.net/gh/33kk)\/.*\.(jpe?g|png|webp)$`,
}

func New() *html.Engine {
	engine := html.NewFileSystem(pkger.Dir("/views"), ".html")

	engine.AddFunc("config", func(key string) template.HTML {
		return template.HTML(fmt.Sprintf("%v", appConfig[key]))
	})

	engine.AddFunc("MarkdownSafe", func(s string) template.HTML {
		gen := md.Run([]byte(s), md.WithExtensions(ext))
		return template.HTML(gen)
	})

	engine.AddFunc("MarkdownUnsafe", func(s string) template.HTML {
		// Generate Markdown then sanitize it before returning HTML.
		gen := md.Run(
			[]byte(s),
			md.WithExtensions(md.HardLineBreak),
		)
		out := bluemonday.UGCPolicy().SanitizeBytes(gen)

		return template.HTML(out)
	})

	engine.AddFunc("RemoveLastNewLine", func(s template.HTML) template.HTML {
		sLen := len(s)
		if sLen > 0 && s[sLen-1] == '\n' {
			s = s[:sLen-1]
		}

		return template.HTML(s)
	})

	engine.AddFunc("GitCommit", func() template.HTML {
		if !config.Production {
			return template.HTML("dev")
		}

		return template.HTML(config.GitCommit)
	})

	engine.AddFunc("Date", func(time time.Time) template.HTML {
		return template.HTML(time.Format("January 02, 2006 15:04"))
	})

	engine.AddFunc("ymdhmsDate", func(time time.Time) template.HTML {
		return template.HTML(time.Format("2006-02-01 15:04:03"))
	})

	engine.AddFunc("subtract", func(a, b int) template.HTML {
		return template.HTML(strconv.FormatInt(int64(a-b), 10))
	})

	engine.AddFunc("paginate", func(page int, sort string) template.HTML {
		s := fmt.Sprintf("/explore?page=%v", page)
		if sort != "" {
			s += fmt.Sprintf("&sort=%v", sort)
		}

		return template.HTML(strings.QueryUnescape(s))
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
			return template.HTML(config.BaseURL())
		}
		return template.HTML(config.BaseURL() + "/" + url.(string))
	})

	if !config.Production {
		engine.Reload(true)
	}

	return engine
}
