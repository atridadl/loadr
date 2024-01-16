package lib

import (
	"fmt"
	"math"
	"net/http"
	"os"
	"path/filepath"
	"sync"
	"time"
)

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

// updateMetrics updates the performance metrics.
func UpdateMetrics(duration time.Duration, resp *http.Response, second int) {
	metrics.mu.Lock()
	defer metrics.mu.Unlock()

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
}

// calculateAndPrintMetrics calculates and prints the performance metrics.
func CalculateAndPrintMetrics(startTime time.Time, requestsPerSecond float64, endpoint string, verb string) {
	averageLatency := time.Duration(0)
	if metrics.totalRequests > 0 {
		averageLatency = metrics.totalLatency / time.Duration(metrics.totalRequests)
	}

	totalDuration := time.Since(startTime).Seconds()
	totalResponses := int32(0)
	for _, count := range metrics.responseCounters {
		totalResponses += count
	}

	// Format the results
	results := fmt.Sprintf("Endpoint: %s\n", endpoint)
	results += fmt.Sprintf("HTTP Verb: %s\n", verb)
	results += fmt.Sprintln("--------------------")
	results += fmt.Sprintln("Performance Metrics:")
	results += fmt.Sprintf("Total Requests Sent: %d\n", metrics.totalRequests)
	results += fmt.Sprintf("Total Responses Received: %d\n", totalResponses)
	results += fmt.Sprintf("Average Latency: %s\n", averageLatency)
	results += fmt.Sprintf("Max Latency: %s\n", metrics.maxLatency)
	results += fmt.Sprintf("Min Latency: %s\n", metrics.minLatency)
	results += fmt.Sprintf("Requests Per Second (Sent): %.2f\n", float64(requestsPerSecond))
	results += fmt.Sprintf("Responses Per Second (Received): %.2f\n", float64(totalResponses)/totalDuration)

	// Print the results to the console
	fmt.Println(results)

	// Save the results to a file
	resultsDir := ".reports"
	os.MkdirAll(resultsDir, os.ModePerm) // Ensure the directory exists

	// Use the current epoch timestamp as the filename
	resultsFile := filepath.Join(resultsDir, fmt.Sprintf("%d.txt", time.Now().Unix()))
	f, err := os.Create(resultsFile)
	if err != nil {
		fmt.Println("Error creating file: ", err)
		return
	}
	defer f.Close()

	_, err = f.WriteString(results)
	if err != nil {
		fmt.Println("Error writing to file: ", err)
		return
	}

	fmt.Println("Results saved to: ", resultsFile)
}
