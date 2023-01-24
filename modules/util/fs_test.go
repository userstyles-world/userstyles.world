package util

import (
	"testing"
	"testing/fstest"
)

var fsys = fstest.MapFS{
	"static/robots.txt":                   {},
	"static/favicon.ico":                  {},
	"static/css/main.css":                 {},
	"static/js/main.js":                   {},
	"static/fonts/Inter-Bold.woff2":       {},
	"static/fonts/Inter-Italic.woff2":     {},
	"static/fonts/Inter-BoldItalic.woff2": {},
	"static/fonts/Inter-Regular.woff2":    {},
}

func TestSubFS(t *testing.T) {
	t.Parallel()

	sub, err := SubFS(fsys, "static")
	if err != nil {
		t.Fatal(err)
	}

	want := []string{"favicon.ico", "js/main.js", "css/main.css"}
	if err := fstest.TestFS(sub, want...); err != nil {
		t.Fatal(err)
	}
}

func TestEmbedFS(t *testing.T) {
	t.Parallel()

	cases := []struct {
		name     string
		env      bool
		root     string
		expected []string
	}{
		{"dev", false, ".", []string{"fs.go", "fs_test.go", "util.go"}},
		{"prod", true, "static", []string{"favicon.ico", "js/main.js", "css/main.css"}},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			got, err := EmbedFS(fsys, c.root, c.env)
			if err != nil {
				t.Fatal(err)
			}

			if err := fstest.TestFS(got, c.expected...); err != nil {
				t.Error(err)
			}
		})
	}
}
