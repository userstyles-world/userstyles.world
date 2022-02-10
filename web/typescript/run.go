//go:build script
// +build script

package main

import (
	"os"

	"github.com/evanw/esbuild/pkg/api"
	"userstyles.world/modules/config"
	"userstyles.world/modules/log"
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
		// Ensure we're seeing the error messages in stdout.
		config.Production = false
		log.Initialize()
		watch = &api.WatchMode{
			OnRebuild: func(result api.BuildResult) {
				if len(result.Errors) > 0 {
					log.Info.Printf("Watch build failed: %d errors\n", len(result.Errors))
				} else {
					log.Info.Printf("Watch build succeeded: %d warnings\n", len(result.Warnings))
				}
			},
		}
	}
	buildResult := api.Build(api.BuildOptions{
		EntryPoints:       []string{"./web/typescript/main.ts"},
		Outfile:           "./web/static/js/main.js",
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
	// From Go 1.18 this section will be optimized away into (amd64).
	// any_label:
	//    PCDATA  $1, $0
	//    CALL    runtime.block(SB)
	//    XCHGL   AX, AX
	if shouldWatch {
		select {}
	}
}
