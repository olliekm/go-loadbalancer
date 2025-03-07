package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"strings"
)

var serverPool ServerPool

func main() {
	var serverList string
	var port int
	// Takes input of servers (urls)
	flag.StringVar(&serverList, "backends", "", "Load balanced backends, use commas to separate")
	// Takes port to serve lb to
	flag.IntVar(&port, "port", 3030, "Port to serve")
	flag.Parse()

	// Must have atleast 1 backend
	if len(serverList) == 0 {
		log.Fatal("Please provide one or more backends to load balance")
	}

	// parse servers
	tokens := strings.Split(serverList, ",")
	for _, tok := range tokens {
		serverUrl, err := url.Parse(tok) // Parse the raw token to valid URL
		if err != nil {
			log.Fatal(err)
		}

		proxy := createReverseProxy(serverUrl)
		// Passes errors we can just add the backend
		serverPool.AddBackend(&Backend{
			URL:          serverUrl,
			Alive:        true,
			ReverseProxy: proxy,
		})
		log.Printf("Configured server: %s\n", serverUrl)
	}

	// create http server
	server := http.Server{
		Addr:    fmt.Sprintf(":%d", port),
		Handler: http.HandlerFunc(lb),
	}

	// start health checking
	go healthCheck() // Concurrent goroutine

	log.Printf("Load Balancer started at :%d\n", port)
	if err := server.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}
