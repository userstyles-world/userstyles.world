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

	engine.AddFunc("MarkdownNew", func(s string) template.HTML {
		gen := md.Run([]byte(s), md.WithExtensions(ext))
		return template.HTML(gen)
	})

	engine.AddFunc("Markdown", func(s string) template.HTML {
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

	if !config.Production {
		engine.Reload(true)
	}

	return engine
}
