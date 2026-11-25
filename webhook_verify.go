package chronary

import (
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"math"
	"net/http"
	"strconv"
	"strings"
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

// SignedDelivery is the delivery metadata the v2 signature (X-Signature-V2)
// covers, next to X-Timestamp and the raw body.
type SignedDelivery struct {
	// DeliveryID is X-Delivery-Id (whd_...).
	DeliveryID string
	// EventType is X-Chronary-Event-Type, for example "event.created".
	EventType string
	// Test is true for a test send (X-Chronary-Test: 1).
	Test bool
	// ReplayOf is X-Chronary-Replay-Of on a replay; "" otherwise.
	ReplayOf string
}

// v2Message is the v2 signed message: "v2", the timestamp, delivery id,
// event type, test flag ("1" or "0"), replay-of id ("" when none) and the
// raw body, joined by LF. No header field can hold an LF, so none can shift.
func v2Message(timestamp string, d SignedDelivery, payload []byte) []byte {
	test := "0"
	if d.Test {
		test = "1"
	}
	head := strings.Join([]string{"v2", timestamp, d.DeliveryID, d.EventType, test, d.ReplayOf, ""}, "\n")
	return append([]byte(head), payload...)
}

func hmacSHA256(secret string, message []byte) string {
	mac := hmac.New(sha256.New, []byte(secret))
	mac.Write(message)
	return "sha256=" + hex.EncodeToString(mac.Sum(nil))
}

// VerifySignature verifies a webhook delivery's X-Signature-V2. The signature
// covers the raw body and the X-Timestamp, X-Delivery-Id,
// X-Chronary-Event-Type, X-Chronary-Test and X-Chronary-Replay-Of headers,
// so once it returns nil those headers can be trusted. Returns nil if valid,
// an error if invalid. During a secret rotation's grace window, deliveries
// also carry X-Signature-V2-Previous (signed with the replaced secret); the
// payload is valid when secret matches either header. The deprecated v1
// X-Signature is never accepted.
func VerifySignature(payload []byte, headers http.Header, secret string) error {
	return VerifySignatureWithOptions(payload, headers, secret)
}

// VerifySignatureWithOptions verifies a webhook signature with custom options.
func VerifySignatureWithOptions(payload []byte, headers http.Header, secret string, opts ...VerifyOption) error {
	_, err := verifyV2(payload, headers, secret, opts...)
	return err
}

func verifyV2(payload []byte, headers http.Header, secret string, opts ...VerifyOption) (SignedDelivery, error) {
	cfg := verifyConfig{
		tolerance: 5 * time.Minute,
		now:       time.Now,
	}
	for _, opt := range opts {
		opt(&cfg)
	}

	// X-Signature (v1) is never read: it does not cover the delivery id,
	// event type, test flag or replay id, so accepting it would let a
	// stripped X-Signature-V2 downgrade the check (#2569).
	sig := headers.Get("X-Signature-V2")
	if sig == "" {
		return SignedDelivery{}, fmt.Errorf("chronary: missing X-Signature-V2 header (this SDK verifies only v2-signed deliveries; X-Signature (v1) is not accepted)")
	}
	// Sent during a secret rotation's grace window, signed with the replaced
	// secret over the same message.
	previousSig := headers.Get("X-Signature-V2-Previous")

	tsStr := headers.Get("X-Timestamp")
	if tsStr == "" {
		return SignedDelivery{}, fmt.Errorf("chronary: missing X-Timestamp header")
	}
	d := SignedDelivery{
		DeliveryID: headers.Get("X-Delivery-Id"),
		EventType:  headers.Get("X-Chronary-Event-Type"),
		ReplayOf:   headers.Get("X-Chronary-Replay-Of"),
	}
	if d.DeliveryID == "" {
		return SignedDelivery{}, fmt.Errorf("chronary: missing X-Delivery-Id header")
	}
	if d.EventType == "" {
		return SignedDelivery{}, fmt.Errorf("chronary: missing X-Chronary-Event-Type header")
	}
	// Absent (or empty) means a live delivery; Chronary sends exactly "1" on a test.
	switch headers.Get("X-Chronary-Test") {
	case "":
	case "1":
		d.Test = true
	default:
		return SignedDelivery{}, fmt.Errorf("chronary: invalid X-Chronary-Test header: expected 1 or no header")
	}
	for _, v := range []string{tsStr, d.DeliveryID, d.EventType, d.ReplayOf} {
		if strings.ContainsAny(v, "\r\n") {
			return SignedDelivery{}, fmt.Errorf("chronary: webhook header contains CR or LF")
		}
	}

	ts, err := strconv.ParseInt(tsStr, 10, 64)
	if err != nil {
		return SignedDelivery{}, fmt.Errorf("chronary: invalid X-Timestamp header: %w", err)
	}

	// Check timestamp tolerance
	now := cfg.now()
	diff := math.Abs(float64(now.Unix() - ts))
	if diff > cfg.tolerance.Seconds() {
		return SignedDelivery{}, fmt.Errorf("chronary: webhook timestamp too old (age: %.0fs, tolerance: %.0fs)", diff, cfg.tolerance.Seconds())
	}

	expected := hmacSHA256(secret, v2Message(tsStr, d, payload))

	// Accept the secret matching either header, so a receiver can switch
	// secrets at any point in a rotation's window. Both comparisons always
	// run, in constant time.
	matchesCurrent := hmac.Equal([]byte(expected), []byte(sig))
	matchesPrevious := previousSig != "" && hmac.Equal([]byte(expected), []byte(previousSig))
	if !matchesCurrent && !matchesPrevious {
		return SignedDelivery{}, fmt.Errorf("chronary: webhook signature mismatch")
	}

	return d, nil
}

// ConstructEvent verifies the signature and parses a webhook payload into a
// WebhookEvent carrying its signed metadata.
func ConstructEvent(payload []byte, headers http.Header, secret string, opts ...VerifyOption) (*WebhookEvent, error) {
	d, err := verifyV2(payload, headers, secret, opts...)
	if err != nil {
		return nil, err
	}

	var data map[string]interface{}
	if err := json.Unmarshal(payload, &data); err != nil {
		return nil, fmt.Errorf("chronary: parsing webhook payload: %w", err)
	}
	return &WebhookEvent{
		Type:       WebhookEventType(d.EventType),
		Data:       data,
		DeliveryID: d.DeliveryID,
		Test:       d.Test,
		ReplayOf:   d.ReplayOf,
	}, nil
}

// ComputeSignature computes the deprecated v1 X-Signature for a webhook
// payload: HMAC-SHA256 over "<timestamp>.<payload>". It does not cover the
// delivery metadata and VerifySignature does not accept it; use
// ComputeSignatureV2. Kept for testing hand-written v1 verifiers.
func ComputeSignature(payload []byte, secret string, timestamp int64) string {
	return hmacSHA256(secret, []byte(fmt.Sprintf("%d.%s", timestamp, string(payload))))
}

// ComputeSignatureV2 computes the X-Signature-V2 value for a webhook payload
// and its delivery metadata. This is useful for testing webhook handlers.
func ComputeSignatureV2(payload []byte, secret string, timestamp int64, d SignedDelivery) string {
	return hmacSHA256(secret, v2Message(strconv.FormatInt(timestamp, 10), d, payload))
}

// SignedHeaders returns HTTP headers with a valid signature for a webhook
// payload, as an event.created delivery with a generated delivery id.
// This is useful for testing webhook handlers.
func SignedHeaders(payload []byte, secret string) http.Header {
	return SignedHeadersWithEventType(payload, secret, "")
}

// SignedHeadersWithEventType returns HTTP headers with a valid signature and
// event type (event.created when eventType is empty) and a generated delivery
// id. This is useful for testing webhook handlers that call ConstructEvent.
func SignedHeadersWithEventType(payload []byte, secret string, eventType string) http.Header {
	return SignedHeadersFor(payload, secret, SignedDelivery{EventType: eventType})
}

// SignedHeadersFor returns the headers Chronary sends for a delivery: the
// v2 and v1 signatures, X-Timestamp (now), X-Delivery-Id,
// X-Chronary-Event-Type, and X-Chronary-Test / X-Chronary-Replay-Of when set.
// An empty DeliveryID is generated (whd_test_<hex>) and an empty EventType
// is event.created. This is useful for testing webhook handlers.
func SignedHeadersFor(payload []byte, secret string, d SignedDelivery) http.Header {
	if d.DeliveryID == "" {
		d.DeliveryID = "whd_test_" + randomHex(8)
	}
	if d.EventType == "" {
		d.EventType = string(WebhookEventEventCreated)
	}
	ts := time.Now().Unix()
	h := http.Header{}
	h.Set("X-Timestamp", strconv.FormatInt(ts, 10))
	h.Set("X-Signature-V2", ComputeSignatureV2(payload, secret, ts, d))
	h.Set("X-Signature", ComputeSignature(payload, secret, ts))
	h.Set("X-Delivery-Id", d.DeliveryID)
	h.Set("X-Chronary-Event-Type", d.EventType)
	if d.Test {
		h.Set("X-Chronary-Test", "1")
	}
	if d.ReplayOf != "" {
		h.Set("X-Chronary-Replay-Of", d.ReplayOf)
	}
	h.Set("Content-Type", "application/json")
	return h
}

func randomHex(n int) string {
	b := make([]byte, n)
	if _, err := rand.Read(b); err != nil {
		return strconv.FormatInt(time.Now().UnixNano(), 16)
	}
	return hex.EncodeToString(b)
}
