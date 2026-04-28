package chronary

import (
	"net/http"
	"time"
)

// ClientOption configures the Client at construction time.
type ClientOption func(*clientConfig)

type clientConfig struct {
	apiKey     string
	anonymous  bool
	baseURL    string
	httpClient *http.Client
	timeout    time.Duration
	maxRetries int
}

// WithAPIKey sets the API key for authentication.
func WithAPIKey(key string) ClientOption {
	return func(c *clientConfig) { c.apiKey = key }
}

// WithAnonymous constructs a client with no API key. Only the unauthenticated
// endpoints work on such a client — currently AgentAuth.SignUp and Plans.List.
// Authenticated calls return a 401 from the API.
func WithAnonymous() ClientOption {
	return func(c *clientConfig) { c.anonymous = true }
}

// WithBaseURL sets the base URL for the API (default: https://api.chronary.ai).
func WithBaseURL(url string) ClientOption {
	return func(c *clientConfig) { c.baseURL = url }
}

// WithHTTPClient sets a custom http.Client.
func WithHTTPClient(hc *http.Client) ClientOption {
	return func(c *clientConfig) { c.httpClient = hc }
}

// WithTimeout sets the default request timeout (default: 30s).
func WithTimeout(d time.Duration) ClientOption {
	return func(c *clientConfig) { c.timeout = d }
}

// WithMaxRetries sets the maximum number of retries for failed requests (default: 2).
func WithMaxRetries(n int) ClientOption {
	return func(c *clientConfig) { c.maxRetries = n }
}

// RequestOption overrides behavior for a single API request.
type RequestOption func(*requestConfig)

type requestConfig struct {
	idempotencyKey string
	timeout        time.Duration
	headers        http.Header
}

// WithIdempotencyKey overrides the auto-generated idempotency key for a request.
func WithIdempotencyKey(key string) RequestOption {
	return func(c *requestConfig) { c.idempotencyKey = key }
}

// WithRequestTimeout overrides the client timeout for a single request.
func WithRequestTimeout(d time.Duration) RequestOption {
	return func(c *requestConfig) { c.timeout = d }
}

// WithHeader adds a custom header to a single request.
func WithHeader(key, value string) RequestOption {
	return func(c *requestConfig) {
		if c.headers == nil {
			c.headers = make(http.Header)
		}
		c.headers.Set(key, value)
	}
}
