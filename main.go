package main

import (
	"flag"
	"fmt"
	"net/http"
	"sync/atomic"
	"time"
)

var client = &http.Client{}

func makeGetRequest(url string) {
	_, err := client.Get(url)
	if err != nil {
		fmt.Println("Error making GET request:", err)
		return
	}
	fmt.Println("Request Sent")
}

func main() {
	// Define command-line flags
	requestsPerSecond := flag.Float64("rate", 10, "Number of requests per second")
	url := flag.String("url", "https://example.com", "The URL to make requests to")

	// Parse the flags
	flag.Parse()

	rateLimit := time.Second / time.Duration(*requestsPerSecond)
	ticker := time.NewTicker(rateLimit)
	defer ticker.Stop()

	var requestCount int32 = 0

	for range ticker.C {
		go func(u string) {
			makeGetRequest(u)
			count := atomic.AddInt32(&requestCount, 1)
			fmt.Printf("Number of requests: %d\n", count)
		}(*url)
	}
}
