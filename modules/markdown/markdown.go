// Package markdown provides helper functions for rendering Markdown documents.
package markdown

import (
	"bytes"

	"github.com/microcosm-cc/bluemonday"
	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark-meta"
	"github.com/yuin/goldmark/extension"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/renderer/html"

	xhtml "golang.org/x/net/html"

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
	docs = goldmark.New(
		goldmark.WithExtensions(
			extension.Footnote,
			extension.GFM,
			meta.Meta,
		),
		goldmark.WithParserOptions(
			parser.WithAutoHeadingID(),
		),
	)
	fallback    = "<mark>Failed to convert Markdown. Please try again.</mark>"
	unreachable = "<mark>Unreachable. Please contact us if you see this.<mark>"
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
	var buf bytes.Buffer
	if err := md.Convert(text, &buf); err != nil {
		log.Warn.Printf("Failed to convert %16q: %v\n", text, err)
		return fallback
	}

	// TODO(vednoc): Explore moving this up or down.
	doc, err := xhtml.Parse(&buf)
	if err != nil {
		log.Warn.Printf("Failed to parse %16q: %v\n", text, err)
		return fallback
	}

	if body := getContent(doc); body != nil {
		for child := body.FirstChild; child != nil; child = child.NextSibling {
			err := xhtml.Render(&buf, child)
			if err != nil {
				log.Warn.Printf("Failed to render %q: %v\n", child.Data, err)
				return fallback
			}
		}
		return bm.Sanitize(buf.String())
	}

	return unreachable
}

func RenderDocs(text []byte) (string, map[string]interface{}) {
	var buf bytes.Buffer
	ctx := parser.NewContext()
	err := docs.Convert([]byte(text), &buf, parser.WithContext(ctx))
	if err != nil {
		log.Warn.Print(err)
		return "", nil
	}

	m := meta.Get(ctx)

	return buf.String(), m
}

func convert(text []byte) (string, error) {
	var buf bytes.Buffer
	if err := md.Convert(text, &buf); err != nil {
		return "", err
	}

	return buf.String(), nil
}

// getContent traverses the document and tries to return content in body tag.
func getContent(n *xhtml.Node) *xhtml.Node {
	if n.Type == xhtml.ElementNode && n.Data == "body" {
		return n
	}

	for c := n.FirstChild; c != nil; c = c.NextSibling {
		if content := getContent(c); content != nil {
			return content
		}
	}

	return nil
}
