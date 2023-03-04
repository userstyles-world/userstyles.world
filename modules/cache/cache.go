package cache

import (
	"os"
	"path"
	"time"

	caching "codeberg.org/Gusted/algorithms-go/caching/lru"
	"github.com/patrickmn/go-cache"

	"userstyles.world/modules/config"
	"userstyles.world/modules/log"
)

var (
	CachedIndex = path.Join(config.CacheDir, "uso-format.json")
	Store       = cache.New(10*time.Minute, time.Minute)
	LRU         = caching.CreateLRUCache(config.CachedCodeItems)
)

func init() {
	dirs := [...]string{
		config.CacheDir,
		config.ImageDir,
		config.ProxyDir,
		config.PublicDir,
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
}

func Initialize() {
	b, err := os.ReadFile(CachedIndex)
	if err != nil {
		log.Warn.Println("Failed to read uso-format.json:", err)
		return
	}

	Store.Set("index", b, 0)
}

func SaveToDisk(f string, data []byte) error {
	return os.WriteFile(f, data, 0o600)
}
