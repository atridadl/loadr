package lib_test

import (
	"loadr/lib"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestSendRequests(t *testing.T) {
	// Create a test server that responds with a 200 status.
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	defer ts.Close()

	tests := []struct {
		name              string
		url               string
		bearerToken       string
		requestType       string
		jsonData          []byte
		maxRequests       int
		requestsPerSecond float64
	}{
		{
			name:              "Test 1",
			url:               ts.URL,
			bearerToken:       "testToken",
			requestType:       "GET",
			jsonData:          []byte(`{"key":"value"}`),
			maxRequests:       5,
			requestsPerSecond: 1.0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			start := time.Now()
			lib.SendRequests(tt.url, tt.bearerToken, tt.requestType, tt.jsonData, tt.maxRequests, tt.requestsPerSecond)
			elapsed := time.Since(start)

			// Check if the requests were sent within the expected time frame.
			if elapsed > time.Duration(float64(tt.maxRequests)*1.5)*time.Second {
				t.Errorf("SendRequests() took too long, got: %v, want: less than %v", elapsed, time.Duration(float64(tt.maxRequests)*1.5)*time.Second)
			}
		})
	}
}
