package chronary

import (
	"bytes"
	"context"
	"crypto/rand"
	"encoding/json"
	"fmt"
	"io"
	"math"
	"math/big"
	"net/http"
	"net/url"
	"reflect"
	"strconv"
	"strings"
	"time"
)

// service is the base for all resource services.
type service struct {
	client *Client
}

// do performs an HTTP request, decoding the JSON response into v.
func (c *Client) do(ctx context.Context, method, path string, body interface{}, v interface{}, opts ...RequestOption) error {
	rc := resolveRequestConfig(c, opts)

	var bodyReader io.Reader
	if body != nil {
		b, err := json.Marshal(body)
		if err != nil {
			return fmt.Errorf("chronary: encoding request body: %w", err)
		}
		bodyReader = bytes.NewReader(b)
	}

	var lastErr error
	for attempt := 0; attempt <= c.maxRetries; attempt++ {
		if attempt > 0 {
			delay := backoff(attempt, lastErr)
			select {
			case <-ctx.Done():
				return &Error{Type: "timeout", Message: ctx.Err().Error(), StatusCode: 0}
			case <-time.After(delay):
			}
			// Reset body reader for retry
			if body != nil {
				b, _ := json.Marshal(body)
				bodyReader = bytes.NewReader(b)
			}
		}

		reqCtx := ctx
		if rc.timeout > 0 {
			var cancel context.CancelFunc
			reqCtx, cancel = context.WithTimeout(ctx, rc.timeout)
			defer cancel()
		}

		req, err := http.NewRequestWithContext(reqCtx, method, c.baseURL+path, bodyReader)
		if err != nil {
			return fmt.Errorf("chronary: creating request: %w", err)
		}

		// Headers
		if c.apiKey != "" {
			req.Header.Set("Authorization", "Bearer "+c.apiKey)
		}
		req.Header.Set("User-Agent", "chronary-go/"+Version)
		req.Header.Set("X-Chronary-SDK-Version", Version)
		if body != nil {
			req.Header.Set("Content-Type", "application/json")
		}
		req.Header.Set("Accept", "application/json")

		// Idempotency key for mutating methods
		if method != http.MethodGet {
			key := rc.idempotencyKey
			if key == "" {
				key = newUUID()
			}
			req.Header.Set("Idempotency-Key", key)
		}

		// Per-request custom headers
		for k, vals := range rc.headers {
			for _, val := range vals {
				req.Header.Set(k, val)
			}
		}

		resp, err := c.httpClient.Do(req)
		if err != nil {
			if ctx.Err() != nil {
				return &Error{Type: "timeout", Message: ctx.Err().Error(), StatusCode: 0}
			}
			lastErr = &Error{Type: "connection_error", Message: err.Error(), StatusCode: 0}
			if attempt < c.maxRetries {
				continue
			}
			return lastErr
		}
		defer resp.Body.Close()

		respBody, err := io.ReadAll(resp.Body)
		if err != nil {
			return fmt.Errorf("chronary: reading response: %w", err)
		}

		if resp.StatusCode >= 200 && resp.StatusCode < 300 {
			if v != nil && len(respBody) > 0 {
				if err := json.Unmarshal(respBody, v); err != nil {
					return fmt.Errorf("chronary: decoding response: %w", err)
				}
			}
			return nil
		}

		apiErr := parseErrorResponse(resp.StatusCode, respBody, resp.Header)
		lastErr = apiErr

		if isRetryable(resp.StatusCode) && attempt < c.maxRetries {
			continue
		}
		return apiErr
	}

	return lastErr
}

// doNoContent performs a request expecting 204 No Content.
func (c *Client) doNoContent(ctx context.Context, method, path string, body interface{}, opts ...RequestOption) error {
	return c.do(ctx, method, path, body, nil, opts...)
}

func resolveRequestConfig(c *Client, opts []RequestOption) requestConfig {
	rc := requestConfig{timeout: c.timeout}
	for _, opt := range opts {
		opt(&rc)
	}
	return rc
}

func isRetryable(statusCode int) bool {
	switch statusCode {
	case 408, 429, 500, 502, 503, 504:
		return true
	}
	return false
}

func backoff(attempt int, lastErr error) time.Duration {
	base := 500 * time.Millisecond

	// Respect Retry-After header from rate limit errors
	if e, ok := lastErr.(*Error); ok && e.RetryAfter > 0 {
		ra := time.Duration(e.RetryAfter) * time.Second
		if ra > 60*time.Second {
			ra = 60 * time.Second
		}
		return ra
	}

	delay := time.Duration(float64(base) * math.Pow(2, float64(attempt-1)))
	if delay > 60*time.Second {
		delay = 60 * time.Second
	}

	// Add jitter: 0..base
	jitter, _ := rand.Int(rand.Reader, big.NewInt(int64(base)))
	if jitter != nil {
		delay += time.Duration(jitter.Int64())
	}
	return delay
}

// newUUID generates a UUID v4 using crypto/rand.
func newUUID() string {
	var uuid [16]byte
	_, _ = rand.Read(uuid[:])
	uuid[6] = (uuid[6] & 0x0f) | 0x40 // version 4
	uuid[8] = (uuid[8] & 0x3f) | 0x80 // variant 2
	return fmt.Sprintf("%x-%x-%x-%x-%x",
		uuid[0:4], uuid[4:6], uuid[6:8], uuid[8:10], uuid[10:16])
}

// addQueryParams encodes a struct with `url` tags into url.Values.
func addQueryParams(params interface{}) url.Values {
	v := url.Values{}
	if params == nil {
		return v
	}

	val := reflect.ValueOf(params)
	if val.Kind() == reflect.Ptr {
		if val.IsNil() {
			return v
		}
		val = val.Elem()
	}
	if val.Kind() != reflect.Struct {
		return v
	}

	t := val.Type()
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		tag := field.Tag.Get("url")
		if tag == "" || tag == "-" {
			continue
		}

		parts := strings.Split(tag, ",")
		name := parts[0]
		omitempty := len(parts) > 1 && parts[1] == "omitempty"

		fv := val.Field(i)
		str := encodeFieldValue(fv)
		if omitempty && str == "" {
			continue
		}
		if str != "" {
			v.Set(name, str)
		}
	}
	return v
}

func encodeFieldValue(fv reflect.Value) string {
	if fv.Kind() == reflect.Ptr {
		if fv.IsNil() {
			return ""
		}
		fv = fv.Elem()
	}

	switch fv.Kind() {
	case reflect.String:
		return fv.String()
	case reflect.Int, reflect.Int64:
		n := fv.Int()
		if n == 0 {
			return ""
		}
		return strconv.FormatInt(n, 10)
	case reflect.Bool:
		if !fv.Bool() {
			return ""
		}
		return "true"
	case reflect.Slice:
		if fv.Len() == 0 {
			return ""
		}
		parts := make([]string, fv.Len())
		for i := 0; i < fv.Len(); i++ {
			parts[i] = fmt.Sprintf("%v", fv.Index(i).Interface())
		}
		return strings.Join(parts, ",")
	}
	return ""
}
