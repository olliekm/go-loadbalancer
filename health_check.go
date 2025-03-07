package main

import (
	"log"
	"net"
	"net/url"
	"time"
)

// isAlive checks whether a backend is Alive by establishing a TCP connection
func isBackendAlive(u *url.URL) bool {
	// Pings backend for response
	timeout := 2 * time.Second
	conn, err := net.DialTimeout("tcp", u.Host, timeout) // The ping
	if err != nil {
		log.Println("Site unreachable, error: ", err)
		return false
	}
	defer conn.Close() // if connection is made, close it and return true
	return true
}

// HealthCheck pings the backends and update the status
func (s *ServerPool) HealthCheck() {
	// Loops through all backends in s, prints status of each and sets alive state accordingly
	for _, b := range s.backends {
		status := "up"
		alive := isBackendAlive(b.URL) // Returns bool, alive or not?
		b.SetAlive(alive)              // Sets to alive state
		if !alive {
			status = "down"
		}
		log.Printf("%s [%s]\n", b.URL, status)
	}
}

// healthCheck runs a routine for check status of the backends every 2 mins
func healthCheck() {
	// Is a concurrent routine
	// Every two minutes, run a health check. Repeats indefinitely
	t := time.NewTicker(time.Minute * 2)
	for {
		select {
		case <-t.C:
			log.Println("Starting health check...")
			serverPool.HealthCheck()
			log.Println("Health check completed")
		}
	}
}
