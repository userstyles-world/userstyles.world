package cache

import (
	"os"
	"path"
	"time"

	"github.com/patrickmn/go-cache"

	"userstyles.world/modules/config"
	"userstyles.world/modules/log"
)

var (
	CacheFile string
	Store     = cache.New(cache.NoExpiration, 5*time.Minute)
	Code      *LRU
)

func Init() {
	dirs := [...]string{
		config.Storage.CacheDir,
		config.Storage.ImageDir,
		config.Storage.ProxyDir,
		config.Storage.PublicDir,
		config.Storage.StyleDir,
	}

	// Create dir if it doesn't exist.
	createIfNotExist := func(dir string) {
		if _, err := os.Stat(dir); os.IsNotExist(err) {
			if err := os.MkdirAll(dir, 0o755); err != nil {
				log.Warn.Fatalf("Failed to create dir %q: %s\n", dir, err)
			}
		}
	}

	// Set up data directories.
	for _, dir := range dirs {
		createIfNotExist(dir)
	}

	// Run install/view stats.
	InstallStats.Run()
	ViewStats.Run()

	CacheFile = path.Join(config.Storage.CacheDir, "cache")
	Code = newLRU(config.App.CodeMaxItems, "code")

	if err := Store.LoadFile(CacheFile); err != nil {
		log.Warn.Println("Failed to read cache:", err)
	}
	log.Info.Println("Loaded cache from disk.")
}

// SaveStore saves the state of in-memory cache to disk.
func SaveStore() {
	if err := Store.SaveFile(CacheFile); err != nil {
		log.Warn.Println("Failed to save cache:", err)
	}
	log.Info.Println("Saved cache to disk.")
}
