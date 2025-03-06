package devsectools

import "context"

// Domain retrieves the parsed domain information from the API.
//
// Parameters:
//   - ctx: Context for handling timeouts and cancellations.
//   - url: The domain to scan (e.g., "example.com").
//
// Returns:
//   - A pointer to a `DomainResponse` struct containing the parsed hostname.
//   - An error if the request fails.
func (c *Client) Domain(ctx context.Context, url string) (*DomainResponse, error) {
	var response DomainResponse
	err := c.makeRequest(ctx, "GET", "/domain?url="+url, nil, &response)
	return &response, err
}

// HTTP retrieves HTTP protocol support information from the API.
//
// Parameters:
//   - ctx: Context for handling timeouts and cancellations.
//   - url: The domain to scan (e.g., "example.com").
//
// Returns:
//   - A pointer to a `HttpResponse` struct containing HTTP version support details.
//   - An error if the request fails.
func (c *Client) HTTP(ctx context.Context, url string) (*HttpResponse, error) {
	var response HttpResponse
	err := c.makeRequest(ctx, "GET", "/http?url="+url, nil, &response)
	return &response, err
}

// TLS retrieves TLS protocol support information from the API.
//
// Parameters:
//   - ctx: Context for handling timeouts and cancellations.
//   - url: The domain to scan (e.g., "example.com").
//
// Returns:
//   - A pointer to a `TlsResponse` struct containing TLS version support details and cipher suites.
//   - An error if the request fails.
func (c *Client) TLS(ctx context.Context, url string) (*TlsResponse, error) {
	var response TlsResponse
	err := c.makeRequest(ctx, "GET", "/tls?url="+url, nil, &response)
	return &response, err
}
