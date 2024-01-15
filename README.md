# Loadr

A lightweight REST load testing tool with rubust support for different verbs, token auth, and performance reports.

## Example using source:

1. `go run main.go -rate=20 -max=100 -url=https://api.example.com/resource -type=POST -json=./data.json -token=YourBearerTokenHere`

## Example using binary:

1. go build
2. `./loadr -rate=20 -max=100 -url=https://api.example.com/resource -type=POST -json=./data.json -token=YourBearerTokenHere`

## Flags:

- `-rate`: Number of requests per second. Default is 10.
- `-max`: Maximum number of requests to send. Must be a non-zero integer. Default is 50.
- `-url`: The URL to make requests to. Default is "https://example.com".
- `-type`: Type of HTTP request. Can be GET, POST, PUT, DELETE, etc. Default is "GET".
- `-json`: Path to the JSON file with request data. If not provided, no data is sent with the requests.
- `-token`: Bearer token for authorization. If not provided, no Authorization header is sent with the requests.

## Reports

Reports are logged at the end of a test run. They are also saved in a directory called `.reports`. All reports are saved as text files with `YYYYMMdd-HHmmss` time format names.
