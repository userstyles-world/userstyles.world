package cache

import (
	"container/list"
	"fmt"
	"sync"
)

// LRU represents a least recently used cache.
type LRU struct {
	cap   int
	list  *list.List
	cache map[int]*list.Element
	mu    sync.Mutex
}

// item represents an element in the linked list.
type item struct {
	k int
	v []byte
}

// Add checks whether a key exists and moves it to the front of the list while
// updating its value, otherwise saves the key and its associated value as a new
// item and moves it to the front of the list, then saves the pointer to the new
// element in the list as a value in the cache for fast lookup.
func (lru *LRU) Add(k int, v []byte) {
	lru.mu.Lock()
	defer lru.mu.Unlock()

	el, ok := lru.cache[k]
	if ok {
		lru.list.MoveToFront(el)
		el.Value.(*item).v = v
		return
	}

	if lru.list.Len() >= lru.cap {
		delete(lru.cache, lru.list.Back().Value.(*item).k)
		lru.list.Remove(lru.list.Back())
	}

	el = lru.list.PushFront(&item{k, v})
	lru.cache[k] = el
}

// Remove checks whether a key exists and removes it, otherwise does nothing.
func (lru *LRU) Remove(k int) {
	lru.mu.Lock()
	defer lru.mu.Unlock()

	el, ok := lru.cache[k]
	if ok {
		delete(lru.cache, k)
		lru.list.Remove(el)
	}
}

// Get checks whether a key exists and returns its value, otherwise returns nil.
func (lru *LRU) Get(k int) []byte {
	lru.mu.Lock()
	defer lru.mu.Unlock()

	el, ok := lru.cache[k]
	if ok {
		lru.list.MoveToFront(el)
		return lru.list.Front().Value.(*item).v
	}

	return nil
}

// debug iterates over all entries in the list and prints them.
func (lru *LRU) debug() {
	lru.mu.Lock()
	defer lru.mu.Unlock()

	for e := lru.list.Front(); e != nil; e = e.Next() {
		fmt.Println(e.Value.(*item).k, string(e.Value.(*item).v))
	}
}

// newLRU initializes a LRU cache.
func newLRU(size int) *LRU {
	return &LRU{
		cap:   size,
		cache: map[int]*list.Element{},
		list:  list.New(),
	}
}
