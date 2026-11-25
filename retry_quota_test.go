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

// bughunt-BH-057 (#2558): a refusal whose Retry-After is beyond the 60 s retry
// cap (the Free iCal manual-sync cooldown, a daily sign-up cap) is returned at
// once instead of sleeping 60 s + 60 s and then returning the same error.
func TestRetryAfterBeyondCapIsNotRetried(t *testing.T) {
	cases := []struct {
		name   string
		status int
		ra     string
		body   string
	}{
		{"ical sync cooldown 429", 429, "900", `{"error":{"type":"rate_limit_error","message":"A sync for this feed started recently."}}`},
		{"daily cap 503", 503, "36000", `{"error":{"type":"anonymous_signup_capacity","message":"Sign-up capacity reached for today."}}`},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			var attempts int32
			srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				atomic.AddInt32(&attempts, 1)
				w.Header().Set("Retry-After", tc.ra)
				w.WriteHeader(tc.status)
				w.Write([]byte(tc.body))
			}))
			defer srv.Close()

			client, _ := NewClient(WithAPIKey("test_key"), WithBaseURL(srv.URL), WithMaxRetries(2))
			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()
			err := client.ICalSubscriptions.Sync(ctx, "ics_1")
			var apiErr *Error
			if !errors.As(err, &apiErr) || apiErr.StatusCode != tc.status {
				t.Fatalf("expected the %d API error, got %v", tc.status, err)
			}
			if got := atomic.LoadInt32(&attempts); got != 1 {
				t.Errorf("expected 1 attempt (Retry-After %s is beyond the cap), got %d", tc.ra, got)
			}
		})
	}
}

func TestRetryAfterWithinCapIsRetried(t *testing.T) {
	var attempts int32
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if atomic.AddInt32(&attempts, 1) == 1 {
			w.Header().Set("Retry-After", "1")
			w.WriteHeader(503)
			w.Write([]byte(`{"error":{"type":"anonymous_signup_capacity","message":"busy"}}`))
			return
		}
		w.WriteHeader(202)
		w.Write([]byte(`{"status":"syncing"}`))
	}))
	defer srv.Close()

	client, _ := NewClient(WithAPIKey("test_key"), WithBaseURL(srv.URL), WithMaxRetries(2))
	if err := client.ICalSubscriptions.Sync(context.Background(), "ics_1"); err != nil {
		t.Fatal(err)
	}
	if got := atomic.LoadInt32(&attempts); got != 2 {
		t.Errorf("expected 2 attempts, got %d", got)
	}
}

// Review of #2748: an HTTP-date Retry-After and a huge delta-seconds value
// (which would overflow time.Duration) are both beyond the cap.
func TestRetryAfterDateAndHugeValueAreNotRetried(t *testing.T) {
	for _, ra := range []string{
		time.Now().Add(2 * time.Hour).UTC().Format(http.TimeFormat),
		"99999999999",
		"9300000000", // past time.Duration's range: wraps negative when multiplied
	} {
		t.Run(ra, func(t *testing.T) {
			var attempts int32
			srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				atomic.AddInt32(&attempts, 1)
				w.Header().Set("Retry-After", ra)
				w.WriteHeader(429)
				w.Write([]byte(`{"error":{"type":"rate_limit_error","message":"slow down"}}`))
			}))
			defer srv.Close()

			client, _ := NewClient(WithAPIKey("test_key"), WithBaseURL(srv.URL), WithMaxRetries(2))
			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()
			_, err := client.Agents.Get(ctx, "agt_1")
			var apiErr *Error
			if !errors.As(err, &apiErr) || apiErr.StatusCode != 429 {
				t.Fatalf("expected the 429 API error, got %v", err)
			}
			if got := atomic.LoadInt32(&attempts); got != 1 {
				t.Errorf("expected 1 attempt, got %d", got)
			}
		})
	}
}

func TestParseErrorResponseReadsDateRetryAfter(t *testing.T) {
	h := http.Header{}
	h.Set("Retry-After", time.Now().Add(2*time.Hour).UTC().Format(http.TimeFormat))
	e := parseErrorResponse(429, []byte(`{"error":{"type":"rate_limit_error","message":"x"}}`), h)
	if e.RetryAfter < 7100 || e.RetryAfter > 7200 {
		t.Errorf("expected about 7200 s, got %d", e.RetryAfter)
	}
	if got := backoff(1, &Error{RetryAfter: 9300000000}); got != 60*time.Second {
		t.Errorf("huge Retry-After: expected the 60s cap, got %v", got)
	}
}
