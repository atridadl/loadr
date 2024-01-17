package lib

import (
	"fmt"
	"math"
	"net/http"
	"os"
	"path/filepath"
	"time"
)

// Initialize the metrics with default values.
var metrics = PerformanceMetrics{
	MinLatency:       time.Duration(math.MaxInt64),
	ResponseCounters: make(map[int]int32),
}

// updateMetrics updates the performance metrics.
func UpdateMetrics(duration time.Duration, resp *http.Response, second int) {
	metrics.Mu.Lock()
	defer metrics.Mu.Unlock()

	metrics.TotalRequests++
	metrics.TotalLatency += duration
	if duration > metrics.MaxLatency {
		metrics.MaxLatency = duration
	}
	if duration < metrics.MinLatency {
		metrics.MinLatency = duration
	}
	if resp.StatusCode == http.StatusOK {
		metrics.TotalResponses++
		metrics.ResponseCounters[second]++
	}
}

// calculateAndPrintMetrics calculates and prints the performance metrics.
func CalculateAndPrintMetrics(startTime time.Time, requestsPerSecond float64, endpoint string, verb string) {
	averageLatency := time.Duration(0)
	if metrics.TotalRequests > 0 {
		averageLatency = metrics.TotalLatency / time.Duration(metrics.TotalRequests)
	}

	totalDuration := time.Since(startTime).Seconds()
	totalResponses := int32(0)
	for _, count := range metrics.ResponseCounters {
		totalResponses += count
	}

	// Format the results
	results := fmt.Sprintf("Endpoint: %s\n", endpoint)
	results += fmt.Sprintf("HTTP Verb: %s\n", verb)
	results += fmt.Sprintln("--------------------")
	results += fmt.Sprintln("Performance Metrics:")
	results += fmt.Sprintf("Total Requests Sent: %d\n", metrics.TotalRequests)
	results += fmt.Sprintf("Total Responses Received: %d\n", totalResponses)
	results += fmt.Sprintf("Average Latency: %s\n", averageLatency)
	results += fmt.Sprintf("Max Latency: %s\n", metrics.MaxLatency)
	results += fmt.Sprintf("Min Latency: %s\n", metrics.MinLatency)
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
