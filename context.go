package main

import "net/http"

// Holds context for number of attempts made to balancer
const (
	Attempts int = iota
	Retry
)

const MaxAttempts = 3

// GetAttemptsFromContext returns the attempts for request
func GetAttemptsFromContext(r *http.Request) int {
	// returns number of attempts made to this load balancer
	if attempts, ok := r.Context().Value(Attempts).(int); ok {
		return attempts
	}
	return 1
}

// GetRetryFromContext returns the retries for request
func GetRetryFromContext(r *http.Request) int {
	// Number of retries for request
	if retry, ok := r.Context().Value(Retry).(int); ok {
		return retry
	}
	return 0
}
