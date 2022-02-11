// Package web provides files for the front-end facilities.
package web

import (
	"embed"
	"log"
	"net/http"

	"userstyles.world/modules/config"
	"userstyles.world/modules/util/httputil"
)

var (
	//go:embed docs static views
	files embed.FS

	// Directories.
	DocsDir   http.FileSystem
	StaticDir http.FileSystem
	ViewsDir  http.FileSystem
)

func init() {
	var err error
	DocsDir, err = httputil.EmbedFS(files, "web/docs", config.Production)
	if err != nil {
		log.Fatalf("Failed to set docs directory: %s\n", err)
	}

	StaticDir, err = httputil.EmbedFS(files, "web/static", config.Production)
	if err != nil {
		log.Fatalf("Failed to set static directory: %s\n", err)
	}

	ViewsDir, err = httputil.EmbedFS(files, "web/views", config.Production)
	if err != nil {
		log.Fatalf("Failed to set views directory: %s\n", err)
	}
}
