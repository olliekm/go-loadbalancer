package main

import (
	"net/http/httputil"
	"net/url"
	"sync"
)

// Simple backend struct
// Reverse proxy allows for requests to pass directly to the backend and allows balancer to return the result
type Backend struct {
	URL          *url.URL
	Alive        bool
	mux          sync.RWMutex
	ReverseProxy *httputil.ReverseProxy
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
