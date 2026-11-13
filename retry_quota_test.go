package chronary

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"sync/atomic"
	"testing"
	"time"
)

// #1936 §1: a monthly plan limit (429 quota_exceeded) is never retried, and a
// server Retry-After never makes the client sleep longer than 60 seconds.

func TestQuotaExceededIsNotRetried(t *testing.T) {
	var attempts int32
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		atomic.AddInt32(&attempts, 1)
		w.Header().Set("Retry-After", "1700000")
		w.WriteHeader(429)
		w.Write([]byte(`{"error":{"type":"quota_exceeded","message":"Monthly API call limit of 1,000 reached.","metric":"api_calls"}}`))
	}))
	defer srv.Close()

	client, _ := NewClient(WithAPIKey("test_key"), WithBaseURL(srv.URL), WithMaxRetries(2))
	start := time.Now()
	_, err := client.Agents.Get(context.Background(), "agt_1")
	if !errors.Is(err, ErrQuotaExceeded) {
		t.Fatalf("expected ErrQuotaExceeded, got %v", err)
	}
	if got := atomic.LoadInt32(&attempts); got != 1 {
		t.Errorf("expected 1 attempt (no retry on quota_exceeded), got %d", got)
	}
	if elapsed := time.Since(start); elapsed > time.Second {
		t.Errorf("expected an immediate return, took %v", elapsed)
	}
}

func TestRateLimit429StillRetries(t *testing.T) {
	var attempts int32
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if atomic.AddInt32(&attempts, 1) == 1 {
			w.Header().Set("Retry-After", "1")
			w.WriteHeader(429)
			w.Write([]byte(`{"error":{"type":"rate_limit_error","message":"Too many requests"}}`))
			return
		}
		w.WriteHeader(200)
		w.Write([]byte(`{"id":"agt_1","name":"a","type":"ai","status":"active"}`))
	}))
	defer srv.Close()

	client, _ := NewClient(WithAPIKey("test_key"), WithBaseURL(srv.URL), WithMaxRetries(2))
	if _, err := client.Agents.Get(context.Background(), "agt_1"); err != nil {
		t.Fatal(err)
	}
	if got := atomic.LoadInt32(&attempts); got != 2 {
		t.Errorf("expected 2 attempts, got %d", got)
	}
}

func TestBackoffHonoursAndClampsRetryAfter(t *testing.T) {
	if got := backoff(1, &Error{Type: "rate_limit_error", RetryAfter: 2}); got != 2*time.Second {
		t.Errorf("Retry-After 2: expected 2s, got %v", got)
	}
	if got := backoff(1, &Error{Type: "rate_limit_error", RetryAfter: 1700000}); got != 60*time.Second {
		t.Errorf("Retry-After 1700000: expected the 60s cap, got %v", got)
	}
}
