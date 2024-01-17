package lib

import (
	"sync"
	"time"
)

// PerformanceMetrics holds the metrics for performance evaluation.
type PerformanceMetrics struct {
	Mu               sync.Mutex // Protects the metrics
	TotalRequests    int32
	TotalResponses   int32
	TotalLatency     time.Duration
	MaxLatency       time.Duration
	MinLatency       time.Duration
	ResponseCounters map[int]int32
}

type RequestError struct {
	Verb string
	URL  string
	Err  error
}
