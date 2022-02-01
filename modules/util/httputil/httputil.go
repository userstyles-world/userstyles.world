// Package httputil provides helper functions for HTTP facilities.
package httputil

import (
	"io/fs"
	"net/http"
)

// ProxyHeader will return proper IP depending on the environment.
func ProxyHeader(production bool) string {
	if production {
		return "X-Real-IP"
	}

	return ""
}

// SubFS will return subtree of the fsys starting from the prefix.
func SubFS(fsys fs.FS, prefix string) (fs.FS, error) {
	sub, err := fs.Sub(fsys, prefix)
	if err != nil {
		return nil, err
	}

	return sub, nil
}

// EmbedFS will return proper http.FileSystem depending on the environment,
// because we don't want to embed files into the executable during development.
func EmbedFS(fsys fs.FS, dir string, production bool) (http.FileSystem, error) {
	if production {
		sub, err := SubFS(fsys, dir)
		if err != nil {
			return nil, err
		}

		return http.FS(sub), nil
	}

	return http.Dir(dir), nil
}
