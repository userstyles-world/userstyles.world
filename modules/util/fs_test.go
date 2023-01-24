package util

import (
	"testing"
	"testing/fstest"
)

var fsys = fstest.MapFS{
	"foo.html":         {},
	"bar.html":         {},
	"baz.html":         {},
	"foo/foo.html":     {},
	"foo/bar.html":     {},
	"foo/baz.html":     {},
	"foo/bar/baz.html": {},
}

func TestSubFS(t *testing.T) {
	t.Parallel()

	sub, err := SubFS(fsys, "foo")
	if err != nil {
		t.Fatal(err)
	}

	want := []string{"foo.html", "bar.html", "baz.html", "bar/baz.html"}
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
		{"prod", true, "foo", []string{"foo.html", "bar.html", "bar/baz.html"}},
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
