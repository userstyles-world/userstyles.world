package util

import (
	"io/fs"
	"net/http"
	"os"
	"path"
)

// ProxyHeader returns proper IP depending on the environment.
func ProxyHeader(production bool) string {
	if production {
		return "X-Real-IP"
	}

	return ""
}

// SubFS returns subtree of the fsys starting from the prefix.
func SubFS(fsys fs.FS, prefix string) (fs.FS, error) {
	return fs.Sub(fsys, prefix)
}

// EmbedFS returns proper http.FileSystem depending on the environment, because
// we don't want to embed files into the executable during development.
func EmbedFS(fsys fs.FS, dir string, production bool) (http.FileSystem, error) {
	if production {
		_, base := path.Split(dir)
		sub, err := SubFS(fsys, base)
		if err != nil {
			return nil, err
		}

		return http.FS(sub), nil
	}

	return http.FS(os.DirFS(dir)), nil
}
