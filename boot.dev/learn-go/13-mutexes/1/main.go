package main

import (
	"sync"
	"time"
)

type safeCounter struct {
	counts map[string]int

	// This uses a pointer because "A Mutex must not be copied after first use. a"
	mu *sync.Mutex
}

func (sc safeCounter) inc(key string) {
	// Simple lock before, unlock after, because the function is small anyway
	sc.mu.Lock()
	sc.slowIncrement(key)
	sc.mu.Unlock()
}

func (sc safeCounter) val(key string) int {
	// Using defer so unlocking happens after returning
	// This way we can call unlock without having to store the return value before returining it
	sc.mu.Lock()
	defer sc.mu.Unlock()
	return sc.slowVal(key)
}

// don't touch below this line

func (sc safeCounter) slowIncrement(key string) {
	tempCounter := sc.counts[key]
	time.Sleep(time.Microsecond)
	tempCounter++
	sc.counts[key] = tempCounter
}

func (sc safeCounter) slowVal(key string) int {
	time.Sleep(time.Microsecond)
	return sc.counts[key]
}
