package main

import (
	"bytes"
	"flag"
	"fmt"
	"io"
	"math"
	"net/http"
	"os"
	"strings"
	"sync"
	"sync/atomic"
	"time"
)

// Global HTTP client used for making requests.
var client = &http.Client{}

// PerformanceMetrics holds the metrics for performance evaluation.
type PerformanceMetrics struct {
	mu sync.Mutex // Protects the metrics

	totalRequests    int32
	totalResponses   int32
	totalLatency     time.Duration
	maxLatency       time.Duration
	minLatency       time.Duration
	responseCounters map[int]int32
}

// Initialize the metrics with default values.
var metrics = PerformanceMetrics{
	minLatency:       time.Duration(math.MaxInt64),
	responseCounters: make(map[int]int32),
}

// makeRequest sends an HTTP request and updates performance metrics.
func makeRequest(verb, url, token string, jsonData []byte, second int) {
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

	// Update the performance metrics.
	metrics.mu.Lock()
	metrics.totalRequests++
	metrics.totalLatency += duration
	if duration > metrics.maxLatency {
		metrics.maxLatency = duration
	}
	if duration < metrics.minLatency {
		metrics.minLatency = duration
	}
	if resp.StatusCode == http.StatusOK {
		metrics.totalResponses++
		metrics.responseCounters[second]++
	}
	metrics.mu.Unlock()

	// Read the response body to determine its size (not shown in the output).
	_, err = io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error reading response body:", err)
		return
	}
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
	maxRequests := flag.Int("max", 50, "Maximum number of requests to send (0 for unlimited)")
	url := flag.String("url", "https://example.com", "The URL to make requests to")
	requestType := flag.String("type", "GET", "Type of HTTP request (GET, POST, PUT, DELETE, etc.)")
	jsonFilePath := flag.String("json", "", "Path to the JSON file with request data")
	bearerToken := flag.String("token", "", "Bearer token for authorization")

	// Parse the command-line flags.
	flag.Parse()

	// Ensure maxRequests is greater than 0.
	if *maxRequests <= 0 {
		fmt.Println("Error: max must be an integer greater than 0")
		return
	}

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

	// Wait for all goroutines to finish.
	var wg sync.WaitGroup

	// Log beginning of requests
	fmt.Println("Starting Loadr Requests...")

	// Start sending requests at the specified rate.
	startTime := time.Now()
	for range ticker.C {
		second := int(time.Since(startTime).Seconds())
		if int(requestCount) >= *maxRequests {
			break
		}
		wg.Add(1)
		go func(u, t, verb string, data []byte, sec int) {
			defer wg.Done()
			makeRequest(verb, u, t, data, sec)
			atomic.AddInt32(&requestCount, 1)
		}(*url, *bearerToken, strings.ToUpper(*requestType), jsonData, second)
	}

	wg.Wait() // Wait for all requests to finish.

	// Calculate and print performance metrics.
	averageLatency := time.Duration(0)
	if metrics.totalRequests > 0 {
		averageLatency = metrics.totalLatency / time.Duration(metrics.totalRequests)
	}

	totalDuration := time.Since(startTime).Seconds()
	totalResponses := int32(0)
	for _, count := range metrics.responseCounters {
		totalResponses += count
	}

	fmt.Printf("Performance Metrics:\n")
	fmt.Printf("Total Requests Sent: %d\n", metrics.totalRequests)
	fmt.Printf("Total Responses Received: %d\n", totalResponses)
	fmt.Printf("Average Latency: %s\n", averageLatency)
	fmt.Printf("Max Latency: %s\n", metrics.maxLatency)
	fmt.Printf("Min Latency: %s\n", metrics.minLatency)
	fmt.Printf("Requests Per Second (Sent): %.2f\n", float64(*requestsPerSecond))
	fmt.Printf("Responses Per Second (Received): %.2f\n", float64(totalResponses)/totalDuration)
}
