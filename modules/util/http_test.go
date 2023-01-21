package util

import (
	"net/http"
	"os"
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

	want := []string{"foo.html", "bar.html", "baz.html", "bar/baz.html"}
	if err := fstest.TestFS(sub, want...); err != nil {
		t.Fatal(err)
	}
}

func TestEmbedFS(t *testing.T) {
	t.Parallel()

	development := false
	got, err := EmbedFS(fsys, "foo", development)
	if err != nil {
		t.Fatal(err)
	}

	want := http.FS(os.DirFS("foo"))
	if got != want {
		t.Fatalf("Got %#v, expected %#v", got, want)
	}

	production := true
	_, err = EmbedFS(fsys, "foo", production)
	if err != nil {
		t.Fatal(err)
	}
}
