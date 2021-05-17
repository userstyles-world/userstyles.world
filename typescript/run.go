// +build script

package main

import (
	"log"
	"os"

	"github.com/evanw/esbuild/pkg/api"
)

var (
	isDebug      = getEnv("DEBUG", "false") == "true"
	isProduction = !isDebug
	shouldWatch  = getEnv("WATCH", "false") == "true"
)

func getEnv(name, fallback string) string {
	if val, set := os.LookupEnv(name); set {
		return val
	}
	return fallback
}

func main() {
	sourceMap := api.SourceMapInline
	var watch *api.WatchMode
	if isProduction {
		sourceMap = api.SourceMapNone
	}
	if shouldWatch {
		watch = &api.WatchMode{
			OnRebuild: func(result api.BuildResult) {
				if len(result.Errors) > 0 {
					log.Printf("watch build failed: %d errors\n", len(result.Errors))
				} else {
					log.Printf("watch build succeeded: %d warnings\n", len(result.Warnings))
				}
			},
		}
	}
	buildResult := api.Build(api.BuildOptions{
		EntryPoints:       []string{"./typescript/main.ts"},
		Outfile:           "./static/js/main.js",
		Bundle:            true,
		Write:             true,
		LogLevel:          api.LogLevelInfo,
		Platform:          api.PlatformBrowser,
		MinifySyntax:      isProduction,
		MinifyWhitespace:  isProduction,
		MinifyIdentifiers: isProduction,
		Sourcemap:         sourceMap,
		Target:            api.ES2017,
		Charset:           api.CharsetUTF8,
		Format:            api.FormatIIFE,
		Watch:             watch,
		Banner: map[string]string{
			"js": "\"use strict\";",
		},
	})

	if len(buildResult.Errors) > 0 {
		os.Exit(1)
	}

	// When we are watching we shouldn't exit program.
	// So a quick and dirty hack to let an never end loop run.
	// From Go 1.15 this section will be optimized away into (amd64).
	// any_label:
	// XCHGL   AX, AX
	// JMP     any_label
	// Whereby XCHGL acts as a "write block".
	if shouldWatch {
		for {
		}
	}
}
