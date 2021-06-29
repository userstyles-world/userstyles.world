package templates

import (
	"html/template"
	"time"

	"github.com/gofiber/template/html"
	"github.com/markbates/pkger"
	"github.com/microcosm-cc/bluemonday"
	md "github.com/russross/blackfriday/v2"

	"userstyles.world/config"
)

var ext = md.CommonExtensions | md.AutoHeadingIDs

func New() *html.Engine {
	engine := html.NewFileSystem(pkger.Dir("/views"), ".html")

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

	engine.AddFunc("GitCommit", func() template.HTML {
		if !config.Production {
			return template.HTML("dev")
		}

		return template.HTML(config.GIT_COMMIT)
	})

	engine.AddFunc("Date", func(time time.Time) template.HTML {
		return template.HTML(time.Format("January 02, 2006 15:04"))
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

	if !config.Production {
		engine.Reload(true)
	}

	return engine
}