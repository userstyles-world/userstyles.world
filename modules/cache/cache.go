package cache

import (
	"encoding/json"
	"os"
	"time"

	caching "codeberg.org/Gusted/algorithms-go/caching/lru"
	"github.com/patrickmn/go-cache"

	"userstyles.world/modules/config"
	"userstyles.world/modules/log"
)

var (
	CachedIndex = "./data/cache/uso-format.json"
	Store       = cache.New(10*time.Minute, time.Minute)
	LRU         = caching.CreateLRUCache(config.CachedCodeItems)
)

func Initialize() {
	b, err := os.ReadFile(CachedIndex)
	if err != nil {
		log.Warn.Println("Failed to read uso-format.json:", err)
		return
	}

	Store.Set("index", b, 0)
}

func SaveToDisk(f string, data interface{}) error {
	b, err := json.Marshal(data)
	if err != nil {
		log.Warn.Println("Failed to marshal JSON:", err)
		return err
	}

	return os.WriteFile(f, b, 0o600)
}
