package main

import (
	"flag"
	"fmt"
	"loadr/lib"
	"os"
)

var version string = "1.0.1"

func parseCommandLine() (float64, int, string, string, string, string) {
	requestsPerSecond := flag.Float64("rate", 10, "Number of requests per second")
	maxRequests := flag.Int("max", 50, "Maximum number of requests to send (0 for unlimited)")
	url := flag.String("url", "https://example.com", "The URL to make requests to")
	requestType := flag.String("type", "GET", "Type of HTTP request (GET, POST, PUT, DELETE, etc.)")
	jsonFilePath := flag.String("json", "", "Path to the JSON file with request data")
	bearerToken := flag.String("token", "", "Bearer token for authorization")
	versionFlag := flag.Bool("version", false, "Print the version and exit")
	versionFlagShort := flag.Bool("v", false, "Print the version and exit")

	// Parse the command-line flags.
	flag.Parse()

	// If the version flag is present, print the version number and exit.
	if *versionFlag || *versionFlagShort {
		fmt.Println("Version:", version)
		os.Exit(0)
	}

	return *requestsPerSecond, *maxRequests, *url, *requestType, *jsonFilePath, *bearerToken
}

// readJSONFile reads the contents of the JSON file at the given path and returns the bytes.
func readJSONFile(filePath string) ([]byte, error) {
	if filePath == "" {
		return nil, nil
	}
	return os.ReadFile(filePath)
}

func main() {
	// Parse the command-line flags.
	requestsPerSecond, maxRequests, url, requestType, jsonFilePath, bearerToken := parseCommandLine()

	// Ensure maxRequests is greater than 0.
	if maxRequests <= 0 {
		fmt.Println("Error: max must be an integer greater than 0")
		return
	}

	// Read the JSON file if the path is provided.
	jsonData, err := readJSONFile(jsonFilePath)
	if err != nil {
		fmt.Println("Error reading JSON file:", err)
		return
	}

	lib.SendRequests(url, bearerToken, requestType, jsonData, maxRequests, requestsPerSecond)
}
