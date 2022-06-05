// Package markdown provides helper functions for rendering Markdown documents.
package markdown

import (
	"bytes"

	"github.com/microcosm-cc/bluemonday"
	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/extension"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/renderer/html"

	"userstyles.world/modules/log"
)

var (
	bm = bluemonday.UGCPolicy()
	md = goldmark.New(
		goldmark.WithExtensions(
			extension.Footnote,
			extension.GFM,
		),
		goldmark.WithParserOptions(
			parser.WithAutoHeadingID(),
		),
		goldmark.WithRendererOptions(
			html.WithHardWraps(),
			html.WithUnsafe(),
		),
	)
	fallback = "<mark>Failed to convert Markdown. Please try again.</mark>"
)

// RenderSafe is used for rendering internal documents.
func RenderSafe(text []byte) string {
	s, err := convert(text)
	if err != nil {
		log.Warn.Printf("Failed to render %16q: %v\n", text, err)
		return fallback
	}

	return s
}

// RenderUnsafe is used for rendering user-generated input.
func RenderUnsafe(text []byte) string {
	s, err := convert(text)
	if err != nil {
		log.Warn.Printf("Failed to render %16q: %v\n", text, err)
		return fallback
	}

	return bm.Sanitize(s)
}

func convert(text []byte) (string, error) {
	var buf bytes.Buffer
	if err := md.Convert(text, &buf); err != nil {
		return "", err
	}

	return buf.String(), nil
}
