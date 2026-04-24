package chronary

import (
	"errors"
	"net/http"
	"testing"
)

func TestErrorIs(t *testing.T) {
	tests := []struct {
		name   string
		err    *Error
		target *Error
		want   bool
	}{
		{"not_found matches ErrNotFound", &Error{Type: "not_found", Message: "gone"}, ErrNotFound, true},
		{"authentication_error matches ErrAuthentication", &Error{Type: "authentication_error"}, ErrAuthentication, true},
		{"rate_limited matches ErrRateLimit", &Error{Type: "rate_limited"}, ErrRateLimit, true},
		{"validation_error matches ErrValidation", &Error{Type: "validation_error"}, ErrValidation, true},
		{"quota_exceeded matches ErrQuotaExceeded", &Error{Type: "quota_exceeded"}, ErrQuotaExceeded, true},
		{"timeout matches ErrTimeout", &Error{Type: "timeout"}, ErrTimeout, true},
		{"connection_error matches ErrConnection", &Error{Type: "connection_error"}, ErrConnection, true},
		{"not_found does not match ErrAuthentication", &Error{Type: "not_found"}, ErrAuthentication, false},
		{"any error matches bare Error", &Error{Type: "not_found"}, &Error{}, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := errors.Is(tt.err, tt.target); got != tt.want {
				t.Errorf("errors.Is() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestErrorAs(t *testing.T) {
	err := &Error{Type: "not_found", Message: "agent not found", StatusCode: 404, RequestID: "req_123"}
	var chronErr *Error
	if !errors.As(err, &chronErr) {
		t.Fatal("errors.As should match *Error")
	}
	if chronErr.StatusCode != 404 {
		t.Errorf("expected status 404, got %d", chronErr.StatusCode)
	}
	if chronErr.RequestID != "req_123" {
		t.Errorf("expected request_id req_123, got %s", chronErr.RequestID)
	}
}

func TestErrorString(t *testing.T) {
	e := &Error{Type: "not_found", Message: "agent not found"}
	if got := e.Error(); got != "chronary: not_found: agent not found" {
		t.Errorf("unexpected: %s", got)
	}

	e.RequestID = "req_abc"
	if got := e.Error(); got != "chronary: not_found: agent not found (request_id: req_abc)" {
		t.Errorf("unexpected: %s", got)
	}
}

func TestParseErrorResponse(t *testing.T) {
	body := []byte(`{"error":{"type":"not_found","message":"Agent not found","request_id":"req_123"}}`)
	h := http.Header{}
	e := parseErrorResponse(404, body, h)
	if e.Type != "not_found" {
		t.Errorf("expected not_found, got %s", e.Type)
	}
	if e.StatusCode != 404 {
		t.Errorf("expected 404, got %d", e.StatusCode)
	}
	if e.RequestID != "req_123" {
		t.Errorf("expected req_123, got %s", e.RequestID)
	}
}

func TestParseErrorResponseRetryAfter(t *testing.T) {
	body := []byte(`{"error":{"type":"rate_limited","message":"Too many requests"}}`)
	h := http.Header{}
	h.Set("Retry-After", "30")
	e := parseErrorResponse(429, body, h)
	if e.RetryAfter != 30 {
		t.Errorf("expected RetryAfter 30, got %d", e.RetryAfter)
	}
}

func TestParseErrorResponseFallback(t *testing.T) {
	body := []byte(`not json`)
	h := http.Header{}
	e := parseErrorResponse(401, body, h)
	if e.Type != "authentication_error" {
		t.Errorf("expected authentication_error, got %s", e.Type)
	}
	if e.StatusCode != 401 {
		t.Errorf("expected 401, got %d", e.StatusCode)
	}
}
