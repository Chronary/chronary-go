package chronary

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"math"
	"net/http"
	"strconv"
	"time"
)

// VerifyOption configures webhook signature verification.
type VerifyOption func(*verifyConfig)

type verifyConfig struct {
	tolerance time.Duration
	now       func() time.Time
}

// WithTolerance sets the maximum age for a webhook timestamp (default: 5 minutes).
func WithTolerance(d time.Duration) VerifyOption {
	return func(c *verifyConfig) { c.tolerance = d }
}

// withNow overrides the current time for testing.
func withNow(fn func() time.Time) VerifyOption {
	return func(c *verifyConfig) { c.now = fn }
}

// VerifySignature verifies the HMAC-SHA256 signature of a webhook payload.
// Returns nil if valid, an error if invalid.
func VerifySignature(payload []byte, headers http.Header, secret string) error {
	return VerifySignatureWithOptions(payload, headers, secret)
}

// VerifySignatureWithOptions verifies a webhook signature with custom options.
func VerifySignatureWithOptions(payload []byte, headers http.Header, secret string, opts ...VerifyOption) error {
	cfg := verifyConfig{
		tolerance: 5 * time.Minute,
		now:       time.Now,
	}
	for _, opt := range opts {
		opt(&cfg)
	}

	sig := headers.Get("X-Signature")
	if sig == "" {
		return fmt.Errorf("chronary: missing X-Signature header")
	}

	tsStr := headers.Get("X-Timestamp")
	if tsStr == "" {
		return fmt.Errorf("chronary: missing X-Timestamp header")
	}

	ts, err := strconv.ParseInt(tsStr, 10, 64)
	if err != nil {
		return fmt.Errorf("chronary: invalid X-Timestamp header: %w", err)
	}

	// Check timestamp tolerance
	now := cfg.now()
	diff := math.Abs(float64(now.Unix() - ts))
	if diff > cfg.tolerance.Seconds() {
		return fmt.Errorf("chronary: webhook timestamp too old (age: %.0fs, tolerance: %.0fs)", diff, cfg.tolerance.Seconds())
	}

	// Compute expected signature: HMAC-SHA256(secret, "<timestamp>.<payload>")
	mac := hmac.New(sha256.New, []byte(secret))
	mac.Write([]byte(fmt.Sprintf("%s.%s", tsStr, string(payload))))
	expected := "sha256=" + hex.EncodeToString(mac.Sum(nil))

	// Constant-time comparison
	if !hmac.Equal([]byte(expected), []byte(sig)) {
		return fmt.Errorf("chronary: webhook signature mismatch")
	}

	return nil
}

// ConstructEvent verifies the signature and parses a webhook payload into a WebhookEvent.
func ConstructEvent(payload []byte, headers http.Header, secret string, opts ...VerifyOption) (*WebhookEvent, error) {
	if err := VerifySignatureWithOptions(payload, headers, secret, opts...); err != nil {
		return nil, err
	}

	eventType := headers.Get("X-Chronary-Event-Type")
	if eventType == "" {
		return nil, fmt.Errorf("chronary: missing X-Chronary-Event-Type header")
	}

	var data map[string]interface{}
	if err := json.Unmarshal(payload, &data); err != nil {
		return nil, fmt.Errorf("chronary: parsing webhook payload: %w", err)
	}
	return &WebhookEvent{Type: eventType, Data: data}, nil
}

// ComputeSignature computes the HMAC-SHA256 signature for a webhook payload.
// This is useful for testing webhook handlers.
func ComputeSignature(payload []byte, secret string, timestamp int64) string {
	mac := hmac.New(sha256.New, []byte(secret))
	mac.Write([]byte(fmt.Sprintf("%d.%s", timestamp, string(payload))))
	return "sha256=" + hex.EncodeToString(mac.Sum(nil))
}

// SignedHeaders returns HTTP headers with a valid signature for a webhook payload.
// This is useful for testing webhook handlers.
func SignedHeaders(payload []byte, secret string) http.Header {
	return SignedHeadersWithEventType(payload, secret, "")
}

// SignedHeadersWithEventType returns HTTP headers with a valid signature and event type.
// This is useful for testing webhook handlers that call ConstructEvent.
func SignedHeadersWithEventType(payload []byte, secret string, eventType string) http.Header {
	ts := time.Now().Unix()
	h := http.Header{}
	h.Set("X-Timestamp", strconv.FormatInt(ts, 10))
	h.Set("X-Signature", ComputeSignature(payload, secret, ts))
	if eventType != "" {
		h.Set("X-Chronary-Event-Type", eventType)
	}
	h.Set("Content-Type", "application/json")
	return h
}
