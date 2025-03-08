package main

import (
	"net/http/httputil"
	"net/url"
	"sync"
	"time"
)

// Simple backend struct
// Reverse proxy allows for requests to pass directly to the backend and allows balancer to return the result
type Backend struct {
	URL           *url.URL
	Alive         bool
	mux           sync.RWMutex
	ReverseProxy  *httputil.ReverseProxy
	ResponseTimes []time.Duration // Stores last N response times
}

// Race condition STUFF
// SetAlive for this backend
func (b *Backend) SetAlive(alive bool) {
	b.mux.Lock()
	// Lock so only one goroutine at a time can access b.Alive
	// i.e. only one can read/write at a time
	b.Alive = alive
	b.mux.Unlock()
}

// IsAlive returns true when backend is alive
func (b *Backend) IsAlive() (alive bool) {
	b.mux.RLock()
	// RLock to multiple goroutines can read (not write) at a time
	// Makes sense for checking if a backend IsAlive
	alive = b.Alive
	b.mux.RUnlock()
	return
}

// Race condition stuff ^^

// Mutate the ResponseList by adding new duration and insuring only 10 responses
func (b *Backend) AddResponseTime(duration time.Duration) {
	b.mux.Lock()
	defer b.mux.Unlock()
	if len(b.ResponseTimes) >= 10 {
		b.ResponseTimes = b.ResponseTimes[1:] // Only keep the most recent 10
	}
	b.ResponseTimes = append(b.ResponseTimes, duration)
}

// Returns average response time of a backend
func (b *Backend) GetAverageResponseTime() time.Duration {
	b.mux.RLock() // Since we only want reading use RLock
	defer b.mux.RUnlock()
	if len(b.ResponseTimes) == 0 {
		return time.Hour // Edge case protection
	}
	var total time.Duration
	for _, t := range b.ResponseTimes {
		total += t // Add each time to total
	}
	return total / time.Duration(len(b.ResponseTimes)) // Turn len to time duration for math to work
}
