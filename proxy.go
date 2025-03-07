package main

import (
	"context"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"time"
)

// lb load balances the incoming request
func lb(w http.ResponseWriter, r *http.Request) {
	// Returns error if MaxAttempts reached
	// Otherwise, get next available sever, reverse proxy it with the request
	attempts := GetAttemptsFromContext(r)
	if attempts > MaxAttempts {
		log.Printf("%s(%s) Max attempts reached, terminating\n", r.RemoteAddr, r.URL.Path)
		http.Error(w, "Service not available", http.StatusServiceUnavailable)
		return
	}

	peer := serverPool.GetNextPeer() // Gets the next available server
	if peer != nil {
		peer.ReverseProxy.ServeHTTP(w, r) // ReverseProxy allows for req/res to pass
		return
	}
	http.Error(w, "Service not available", http.StatusServiceUnavailable)
}

func createReverseProxy(serverUrl *url.URL) *httputil.ReverseProxy {
	// Set up ReverseProxy
	proxy := httputil.NewSingleHostReverseProxy(serverUrl)
	// Custom error handling (overriding basic "Bad Gateway")
	proxy.ErrorHandler = func(writer http.ResponseWriter, request *http.Request, e error) {
		log.Printf("[%s] %s\n", serverUrl.Host, e.Error())
		// Get number of retries of reqs
		retries := GetRetryFromContext(request)
		if retries < 3 {
			select {
			case <-time.After(10 * time.Millisecond): // Only allow every 10ms (reduce load of backend)
				ctx := context.WithValue(request.Context(), Retry, retries+1) // Increment retries context
				proxy.ServeHTTP(writer, request.WithContext(ctx))             // Try req again
			}
			return
		}

		// after 3 retries, mark this backend as down
		serverPool.MarkBackendStatus(serverUrl, false)

		// if the same request routing for few attempts with different backends, increase the count
		attempts := GetAttemptsFromContext(request)
		log.Printf("%s(%s) Attempting retry %d\n", request.RemoteAddr, request.URL.Path, attempts)
		ctx := context.WithValue(request.Context(), Attempts, attempts+1)
		lb(writer, request.WithContext(ctx))
	}
	return proxy
}
