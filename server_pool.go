package main

import (
	"net/url"
	"sync/atomic"
)

// A list of our backends
// current keeps track of the backend we are one
//
//	cur
//
// [s1, s2, s3, s4, s5] so s3 is our current backend routing to
type ServerPool struct {
	backends []*Backend
	current  uint64
}

// AddBackend to the server pool
// Does what it says, adds a backend to our server list
func (s *ServerPool) AddBackend(backend *Backend) {
	s.backends = append(s.backends, backend)
}

func (s *ServerPool) NextIndex() int {
	// Return the next backend index
	// AddUint64 is concurrent safe way to increment current within server pool
	// modulo allows NextIndex to wrap around, (Round Robin)
	return int(atomic.AddUint64(&s.current, uint64(1)) % uint64(len(s.backends)))
}

// MarkBackendStatus changes a status of a backend
func (s *ServerPool) MarkBackendStatus(backendUrl *url.URL, alive bool) {
	// Given the backendUrl and alive, mutate the alive state of this backend
	// Literally loops through and finds the backend, sets to given state and breaks loop
	for _, b := range s.backends {
		if b.URL.String() == backendUrl.String() {
			b.SetAlive(alive)
			break
		}
	}
}

// GetNextPeer returns next active peer to take a connection
func (s *ServerPool) GetNextPeer() *Backend {
	// loop entire backends to find out an Alive backend
	// Find the next available backend
	next := s.NextIndex()
	l := len(s.backends) + next // Starting from NextIndex, loop through (modulo takes care of wrap)
	for i := next; i < l; i++ {
		idx := i % len(s.backends)     // take an index by modding
		if s.backends[idx].IsAlive() { // if we have an alive backend, use it and store if its not the original one
			if i != next {
				atomic.StoreUint64(&s.current, uint64(idx))
			}
			return s.backends[idx]
		}
	}
	return nil // No alive backends
}
