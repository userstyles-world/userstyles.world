package util

import (
	"io/fs"
	"os"
	"path"
)

// SubFS returns subtree of the fsys starting from the prefix.
func SubFS(fsys fs.FS, prefix string) (fs.FS, error) {
	return fs.Sub(fsys, prefix)
}

// EmbedFS returns proper http.FileSystem depending on the environment, because
// we don't want to embed files into the executable during development.
func EmbedFS(fsys fs.FS, dir string, production bool) (fs.FS, error) {
	if production {
		_, base := path.Split(dir)
		return SubFS(fsys, base)
	}

	return os.DirFS(dir), nil
}
