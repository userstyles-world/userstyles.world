// Package web provides files for the front-end facilities.
package web

import (
	"embed"
	"io/fs"
	"log"

	"userstyles.world/modules/config"
	"userstyles.world/modules/util"
)

var (
	//go:embed docs static views
	files embed.FS

	// Directories.
	DocsDir   fs.FS
	StaticDir fs.FS
	ViewsDir  fs.FS
)

func Init() {
	var err error
	DocsDir, err = util.EmbedFS(files, "web/docs", config.App.Production)
	if err != nil {
		log.Fatalf("Failed to set docs directory: %s\n", err)
	}

	StaticDir, err = util.EmbedFS(files, "web/static", config.App.Production)
	if err != nil {
		log.Fatalf("Failed to set static directory: %s\n", err)
	}

	ViewsDir, err = util.EmbedFS(files, "web/views", config.App.Production)
	if err != nil {
		log.Fatalf("Failed to set views directory: %s\n", err)
	}
}
