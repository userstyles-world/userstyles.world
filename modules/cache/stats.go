package cache

import (
	"strings"
	"sync"
	"time"

	"userstyles.world/modules/database"
	"userstyles.world/modules/log"
	"userstyles.world/modules/util"
)

// InstallStats stores stats for installs.
var InstallStats = newStats("install")

// ViewStats stores stats for views.
var ViewStats = newStats("view")

// stats stores moving parts of a stats cache.
type stats struct {
	sync.Mutex
	name   string
	done   chan bool
	m      map[string]string
	ticker *time.Ticker
}

// newStats initializes a specific stats cache.
func newStats(name string) *stats {
	return &stats{
		name:   name,
		done:   make(chan bool),
		m:      make(map[string]string),
		ticker: time.NewTicker(time.Minute),
	}
}

// Run starts a store in a separate goroutine.
func (s *stats) Run() {
	go func() {
		for {
			select {
			case <-s.done:
				return
			case <-s.ticker.C:
				s.UpsertAndEvict()
			}
		}
	}()
}

// Close stops the cache and runs a database query.
func (s *stats) Close() {
	log.Info.Printf("Stopping %q cache.\n", s.name)
	s.done <- true
	s.ticker.Stop()
	s.UpsertAndEvict()
}

// Upsert generates database query, executes it, and returns its status.
func (s *stats) Upsert(count int) error {
	var b strings.Builder
	b.WriteString("INSERT INTO stats(created_at, updated_at, ")
	b.WriteString(s.name + ", hash, style_id) VALUES ")

	ending := ", "
	for k, v := range s.m {
		count--
		if count == 0 {
			ending = " "
		}

		// HACK: Might want to rethink that string split.
		b.WriteString("(DATETIME('now'), DATETIME('now'), DATETIME('now'), '")
		b.WriteString(v + "', " + strings.Split(k, " ")[1] + ")" + ending)
	}

	b.WriteString("ON CONFLICT(hash) DO UPDATE SET ")
	b.WriteString("updated_at = excluded.updated_at, ")
	b.WriteString(s.name + " = excluded." + s.name)

	return database.Conn.Exec(b.String()).Error
}

// UpsertAndEvict runs a database query and resets the map if it isn't empty.
func (s *stats) UpsertAndEvict() {
	s.Lock()
	defer s.Unlock()

	count := len(s.m)
	if count > 0 {
		if err := s.Upsert(count); err != nil {
			log.Database.Printf("Failed to upsert %q: %s\n", s.name, err)
		} else {
			log.Info.Printf("Evicting %d items from %q.\n", count, s.name)
			s.m = make(map[string]string)
		}
	}
}

// Add inserts a key in the map as well as its hashed value.
func (s *stats) Add(key string) {
	s.Lock()
	defer s.Unlock()

	_, found := s.m[key]
	if !found {
		val, err := util.HashIP(key)
		if err != nil {
			log.Info.Printf("Failed to create hash for %q: %s\n", key, err)
			return
		}

		s.m[key] = val
	}
}
