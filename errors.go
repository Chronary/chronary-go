package chronary

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
)

// Error is the base error type returned by the Chronary API.
type Error struct {
	Type       string `json:"type"`
	Message    string `json:"message"`
	StatusCode int    `json:"-"`
	RequestID  string `json:"request_id,omitempty"`
	RetryAfter int    `json:"-"` // seconds; 0 if not present
}

// Sentinel errors for use with errors.Is.
var (
	ErrAuthentication = &Error{Type: "authentication_error"}
	ErrNotFound       = &Error{Type: "not_found"}
	ErrValidation     = &Error{Type: "validation_error"}
	ErrRateLimit      = &Error{Type: "rate_limited"}
	ErrQuotaExceeded  = &Error{Type: "quota_exceeded"}
	ErrTimeout        = &Error{Type: "timeout"}
	ErrConnection     = &Error{Type: "connection_error"}
)

func (e *Error) Error() string {
	if e.RequestID != "" {
		return fmt.Sprintf("chronary: %s: %s (request_id: %s)", e.Type, e.Message, e.RequestID)
	}
	return fmt.Sprintf("chronary: %s: %s", e.Type, e.Message)
}

// Is supports errors.Is matching on the Type field.
func (e *Error) Is(target error) bool {
	t, ok := target.(*Error)
	if !ok {
		return false
	}
	if t.Type == "" {
		return true // bare *Error matches any *Error
	}
	return e.Type == t.Type
}

// errorEnvelope is the JSON envelope: {"error": {...}}
type errorEnvelope struct {
	Error Error `json:"error"`
}

// parseErrorResponse parses an API error response.
func parseErrorResponse(statusCode int, body []byte, headers http.Header) *Error {
	var envelope errorEnvelope
	if err := json.Unmarshal(body, &envelope); err == nil && envelope.Error.Type != "" {
		e := &envelope.Error
		e.StatusCode = statusCode
		if ra := headers.Get("Retry-After"); ra != "" {
			if secs, err := strconv.Atoi(ra); err == nil {
				e.RetryAfter = secs
			}
		}
		return e
	}

	// Fallback: construct error from status code
	e := &Error{
		StatusCode: statusCode,
		Message:    fmt.Sprintf("HTTP %d", statusCode),
	}
	switch statusCode {
	case 401:
		e.Type = "authentication_error"
	case 404:
		e.Type = "not_found"
	case 422:
		e.Type = "validation_error"
	case 429:
		e.Type = "rate_limited"
	case 402:
		e.Type = "quota_exceeded"
	default:
		e.Type = "api_error"
	}
	if ra := headers.Get("Retry-After"); ra != "" {
		if secs, err := strconv.Atoi(ra); err == nil {
			e.RetryAfter = secs
		}
	}
	return e
}
