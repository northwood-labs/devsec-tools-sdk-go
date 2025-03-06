package devsectools

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"sync"
	"time"
)

// Endpoint represents an API endpoint with a base URL.
type Endpoint struct {
	BaseURL string // The base URL of the API.
}

// Predefined API Endpoints.
var (
	PRODUCTION = Endpoint{BaseURL: "https://api.devsec.tools"}
	LOCALDEV   = Endpoint{BaseURL: "http://api.devsec.local"}
)

// Default values.
const (
	DefaultTimeout = 5 * time.Second // Default network timeout (5 seconds)
)

// Config holds configuration settings for the API client.
type Config struct {
	Endpoint *Endpoint   // API endpoint (PRODUCTION, LOCALDEV, or custom)
	Timeout  time.Duration // Network timeout duration
}

// Client represents the DevSecTools API client.
type Client struct {
	httpClient *http.Client
	config     *Config
	once       sync.Once
}

// NewClient initializes a new API client with default settings (PRODUCTION API, 5s timeout).
//
// Returns:
//   - A pointer to the newly created Client.
func NewClient() *Client {
	return NewClientWithConfig(&Config{
		Endpoint: &PRODUCTION,
		Timeout:  DefaultTimeout,
	})
}

// NewClientWithConfig initializes a new API client with custom configuration settings.
//
// Parameters:
//   - config: A pointer to a `Config` struct containing API endpoint and timeout settings.
//
// Returns:
//   - A pointer to the newly created Client.
func NewClientWithConfig(config *Config) *Client {
	client := &Client{
		config: config,
	}
	client.once.Do(func() {
		client.httpClient = &http.Client{Timeout: config.Timeout}
	})
	return client
}

// SetEndpoint updates the API endpoint for the client.
//
// Parameters:
//   - endpoint: A pointer to an `Endpoint` struct (e.g., `&PRODUCTION`, `&LOCALDEV`).
func (c *Client) SetEndpoint(endpoint *Endpoint) {
	c.config.Endpoint = endpoint
}

// SetBaseURL allows setting a custom API base URL.
//
// Parameters:
//   - url: A string representing the new API base URL.
func (c *Client) SetBaseURL(url string) {
	c.config.Endpoint = &Endpoint{BaseURL: url}
}

// SetTimeout updates the network timeout duration for API requests.
//
// Parameters:
//   - timeout: The new timeout duration, specified as a `time.Duration` value (e.g., `10*time.Second`).
func (c *Client) SetTimeout(timeout time.Duration) {
	c.config.Timeout = timeout
	c.httpClient.Timeout = timeout
}

// makeRequest performs an HTTP request with context-based timeout handling.
//
// Parameters:
//   - ctx: A context to allow request cancellation or custom timeouts.
//   - method: The HTTP method (e.g., "GET").
//   - endpoint: The API endpoint path (e.g., "/domain").
//   - payload: The request body (set to `nil` for GET requests).
//   - result: A pointer to a struct where the response will be unmarshaled.
//
// Returns:
//   - An error if the request fails or an API error occurs.
func (c *Client) makeRequest(ctx context.Context, method, endpoint string, payload any, result any) error {
	url := fmt.Sprintf("%s%s", c.config.Endpoint.BaseURL, endpoint)

	ctx, cancel := context.WithTimeout(ctx, c.config.Timeout)
	defer cancel()

	var reqBody io.Reader
	if payload != nil {
		data, err := json.Marshal(payload)
		if err != nil {
			return err
		}
		reqBody = bytes.NewReader(data)
	}

	req, err := http.NewRequestWithContext(ctx, method, url, reqBody)
	if err != nil {
		return err
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		var errResp ErrorResponse
		json.NewDecoder(resp.Body).Decode(&errResp)
		return errors.New(errResp.Error)
	}

	return json.NewDecoder(resp.Body).Decode(result)
}

// BatchRequest represents a single request within a batch operation.
type BatchRequest struct {
	Method string      // The API method to call: "domain", "http", or "tls".
	URL    string      // The URL to scan.
	Result interface{} // A pointer to store the result.
	Err    error       // Stores any error encountered.
}

// Batch executes multiple API requests concurrently using Goroutines.
//
// This method improves performance by utilizing concurrency in Go.
//
// Parameters:
//   - ctx: A context to manage request timeouts and cancellations.
//   - requests: A slice of `BatchRequest` structs defining the API calls.
//
// Example Usage:
//
//   batchRequests := []devsectools.BatchRequest{
//       {Method: "domain", URL: "example.com", Result: &devsectools.DomainResponse{}},
//       {Method: "http", URL: "example.com", Result: &devsectools.HttpResponse{}},
//       {Method: "tls", URL: "example.com", Result: &devsectools.TlsResponse{}},
//   }
//
//   client.Batch(context.Background(), batchRequests)
//
//   for _, req := range batchRequests {
//       if req.Err != nil {
//           log.Printf("Error fetching %s: %v\n", req.Method, req.Err)
//           continue
//       }
//       fmt.Printf("Result for %s: %+v\n", req.Method, req.Result)
//   }
func (c *Client) Batch(ctx context.Context, requests []BatchRequest) {
	var wg sync.WaitGroup
	for i := range requests {
		wg.Add(1)
		go func(req *BatchRequest) {
			defer wg.Done()
			var err error
			switch req.Method {
			case "domain":
				req.Result, err = c.Domain(ctx, req.URL)
			case "http":
				req.Result, err = c.HTTP(ctx, req.URL)
			case "tls":
				req.Result, err = c.TLS(ctx, req.URL)
			default:
				err = errors.New("invalid batch request method: " + req.Method)
			}
			if err != nil {
				req.Err = err
			}
		}(&requests[i])
	}
	wg.Wait()
}
