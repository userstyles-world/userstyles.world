package util

import (
	"io/fs"
	"os"
	"path/filepath"
)

// SubFS returns a new file system starting from directory dir.
func SubFS(fsys fs.FS, dir string) (fs.FS, error) {
	return fs.Sub(fsys, dir)
}

// EmbedFS returns a new file system starting from directory dir depending on
// the environment.  This allows us to have reload for internal files during
// development, and also embed them into executable in production.
func EmbedFS(fsys fs.FS, dir string, embed bool) (fs.FS, error) {
	if embed {
		_, base := filepath.Split(dir)
		return SubFS(fsys, base)
	}

	return os.DirFS(dir), nil
}
