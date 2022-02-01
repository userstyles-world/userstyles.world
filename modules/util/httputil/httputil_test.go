package httputil

import (
	"io/fs"
	"net/http"
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

func TestProxyHeader(t *testing.T) {
	t.Parallel()

	development := false
	if ProxyHeader(development) != "" {
		t.Fatal("should return an empty string")
	}

	production := true
	if ProxyHeader(production) != "X-Real-IP" {
		t.Fatal("should return X-Real-IP")
	}
}

func TestSubFS(t *testing.T) {
	t.Parallel()

	sub, err := SubFS(fsys, "foo")
	if err != nil {
		t.Fatal(err)
	}

	got, want := 0, 3
	files, err := fs.ReadDir(sub, ".")
	for _, f := range files {
		if !f.IsDir() {
			got++
		}
	}

	if got != want {
		t.Fatalf("Got %d files, wanted %d", got, want)
	}
}

func TestEmbedFS(t *testing.T) {
	t.Parallel()

	development := false
	got, err := EmbedFS(fsys, "foo", development)
	if err != nil {
		t.Fatal(err)
	}

	want := http.Dir("foo")
	if got != want {
		t.Fatalf("Got %q, expected %q", got, want)
	}

	production := true
	_, err = EmbedFS(fsys, "foo", production)
	if err != nil {
		t.Fatal(err)
	}
}
