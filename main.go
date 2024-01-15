package main

import (
	"bytes"
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"sync/atomic"
	"time"
)

// Global HTTP client used for making requests.
var client = &http.Client{}

// makeRequest sends an HTTP request with the specified verb, URL, bearer token, and JSON data.
// It prints out the status code, response time, and response size.
func makeRequest(verb, url, token string, jsonData []byte) {
	startTime := time.Now()

	// Create a new request with the provided verb, URL, and JSON data if provided.
	var req *http.Request
	var err error
	if jsonData != nil {
		req, err = http.NewRequest(verb, url, bytes.NewBuffer(jsonData))
		req.Header.Set("Content-Type", "application/json")
	} else {
		req, err = http.NewRequest(verb, url, nil)
	}
	if err != nil {
		fmt.Println("Error creating request:", err)
		return
	}

	// Add the bearer token to the request's Authorization header if provided.
	if token != "" {
		req.Header.Set("Authorization", "Bearer "+token)
	}

	// Send the request.
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error making request:", err)
		return
	}
	defer resp.Body.Close()

	// Calculate the duration of the request.
	duration := time.Since(startTime)

	// Read the response body to determine its size.
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error reading response body:", err)
		return
	}

	// Print out the request details.
	fmt.Printf("Request Sent. Status Code: %d, Duration: %s, Response Size: %d bytes\n", resp.StatusCode, duration, len(body))
}

// readJSONFile reads the contents of the JSON file at the given path and returns the bytes.
func readJSONFile(filePath string) ([]byte, error) {
	if filePath == "" {
		return nil, nil
	}
	return os.ReadFile(filePath)
}

func main() {
	// Define command-line flags for configuring the load test.
	requestsPerSecond := flag.Float64("rate", 10, "Number of requests per second")
	maxRequests := flag.Int("max", 0, "Maximum number of requests to send (0 for unlimited)")
	url := flag.String("url", "https://example.com", "The URL to make requests to")
	requestType := flag.String("type", "GET", "Type of HTTP request (GET, POST, PUT, DELETE, etc.)")
	jsonFilePath := flag.String("json", "", "Path to the JSON file with request data")
	bearerToken := flag.String("token", "", "Bearer token for authorization")

	// Parse the command-line flags.
	flag.Parse()

	// Read the JSON file if the path is provided.
	jsonData, err := readJSONFile(*jsonFilePath)
	if err != nil {
		fmt.Println("Error reading JSON file:", err)
		return
	}

	// Calculate the rate limit based on the requests per second.
	rateLimit := time.Second / time.Duration(*requestsPerSecond)
	ticker := time.NewTicker(rateLimit)
	defer ticker.Stop()

	// Initialize the request count.
	var requestCount int32 = 0

	// Start sending requests at the specified rate.
	for range ticker.C {
		// Stop if the maximum number of requests is reached.
		if *maxRequests > 0 && int(requestCount) >= *maxRequests {
			break
		}
		// Send the request in a new goroutine.
		go func(u, t, verb string, data []byte) {
			makeRequest(verb, u, t, data)
			atomic.AddInt32(&requestCount, 1)
		}(*url, *bearerToken, strings.ToUpper(*requestType), jsonData)
	}

	// Print out the total number of requests sent after the load test is finished.
	fmt.Printf("Finished sending requests. Total requests: %d\n", requestCount)
}
