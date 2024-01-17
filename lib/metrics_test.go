package lib_test

import (
	"loadr/lib"
	"reflect"
	"testing"
	"time"
)

func comparePerformanceMetrics(a, b *lib.PerformanceMetrics) bool {
	return a.TotalRequests == b.TotalRequests &&
		a.TotalResponses == b.TotalResponses &&
		a.TotalLatency == b.TotalLatency &&
		a.MaxLatency == b.MaxLatency &&
		a.MinLatency == b.MinLatency &&
		reflect.DeepEqual(a.ResponseCounters, b.ResponseCounters)
}

func TestCalculateAndPrintMetrics(t *testing.T) {
	// Define test cases.
	tests := []struct {
		name              string
		startTime         time.Time
		requestsPerSecond float64
		endpoint          string
		verb              string
		expectedMetrics   lib.PerformanceMetrics
	}{
		{
			name:              "Test 1",
			startTime:         time.Now().Add(-1 * time.Second), // 1 second ago
			requestsPerSecond: 1.0,
			endpoint:          "http://localhost",
			verb:              "GET",
			expectedMetrics: lib.PerformanceMetrics{
				TotalRequests:    1,
				TotalResponses:   1,
				TotalLatency:     1 * time.Second,
				MaxLatency:       1 * time.Second,
				MinLatency:       1 * time.Second,
				ResponseCounters: map[int]int32{1: 1},
			},
		},
		// Add more test cases as needed.
	}

	for i := range tests {
		tt := &tests[i]
		t.Run(tt.name, func(t *testing.T) {
			// Reset the metrics before each test
			metrics := lib.PerformanceMetrics{}

			// Mock the system behavior
			metrics.TotalRequests++
			metrics.ResponseCounters = make(map[int]int32)
			metrics.TotalResponses++
			metrics.TotalLatency += 1 * time.Second
			metrics.MaxLatency = 1 * time.Second
			metrics.MinLatency = 1 * time.Second
			metrics.ResponseCounters[1]++

			// Call the function
			lib.CalculateAndPrintMetrics(tt.startTime, tt.requestsPerSecond, tt.endpoint, tt.verb)

			// Check if the metrics are correct
			if !comparePerformanceMetrics(&metrics, &tt.expectedMetrics) {
				t.Errorf("CalculateAndPrintMetrics() = TotalRequests: %v, TotalResponses: %v, TotalLatency: %v, MaxLatency: %v, MinLatency: %v, ResponseCounters: %v, want TotalRequests: %v, TotalResponses: %v, TotalLatency: %v, MaxLatency: %v, MinLatency: %v, ResponseCounters: %v",
					metrics.TotalRequests, metrics.TotalResponses, metrics.TotalLatency, metrics.MaxLatency, metrics.MinLatency, metrics.ResponseCounters,
					tt.expectedMetrics.TotalRequests, tt.expectedMetrics.TotalResponses, tt.expectedMetrics.TotalLatency, tt.expectedMetrics.MaxLatency, tt.expectedMetrics.MinLatency, tt.expectedMetrics.ResponseCounters)
			}
		})
	}
}
