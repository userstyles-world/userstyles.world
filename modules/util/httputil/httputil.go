// Package httputil provides helper functions for HTTP facilities.
package httputil

import "io/fs"

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
