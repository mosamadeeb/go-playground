package pokecache

import (
	"sync"
	"time"
)

type Cache struct {
	cacheMap map[string]cacheEntry
	mu       *sync.Mutex

	reapStopChan chan struct{}
}

type cacheEntry struct {
	createdAt time.Time
	val       []byte
}

func (c Cache) Add(key string, val []byte) {
	c.mu.Lock()
	c.cacheMap[key] = cacheEntry{time.Now(), val}
	c.mu.Unlock()
}

func (c Cache) Get(key string) (val []byte, ok bool) {
	c.mu.Lock()
	defer c.mu.Unlock()

	entry, ok := c.cacheMap[key]
	if !ok {
		return nil, false
	}

	return entry.val, true
}

func (c Cache) reapLoop(interval time.Duration) {
	// No need to care about stopping the timer when the cache is closed
	// Starting from Go 1.23, unreferenced timers will be recovered even if not stopped
	ticker := time.NewTicker(interval)

	for {
		select {
		case <-ticker.C:
			c.mu.Lock()
			for k, v := range c.cacheMap {
				if time.Now().Sub(v.createdAt) > interval {
					// Removing keys in a range loop is safe
					delete(c.cacheMap, k)
				}
			}
			c.mu.Unlock()
		case <-c.reapStopChan:
			return
		}
	}
}

func (c Cache) Close() error {
	// This channel is expected to be always open
	c.reapStopChan <- struct{}{}
	return nil
}

func NewCache(interval time.Duration) Cache {
	c := Cache{
		cacheMap:     make(map[string]cacheEntry),
		mu:           &sync.Mutex{},
		reapStopChan: make(chan struct{}),
	}

	go c.reapLoop(interval)

	return c
}
