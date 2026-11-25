package chronary

import (
	"net/http"
	"strconv"
	"strings"
	"testing"
	"time"
)

var liveDelivery = SignedDelivery{DeliveryID: "whd_1", EventType: "event.created"}

// deliveryHeaders returns the headers Chronary sends for one delivery signed
// at ts: X-Signature-V2 over the v2 message, the deprecated v1 X-Signature,
// and the metadata headers the v2 signature covers.
func deliveryHeaders(payload []byte, secret string, ts int64, d SignedDelivery) http.Header {
	h := http.Header{}
	h.Set("X-Signature-V2", ComputeSignatureV2(payload, secret, ts, d))
	h.Set("X-Signature", ComputeSignature(payload, secret, ts))
	h.Set("X-Timestamp", strconv.FormatInt(ts, 10))
	h.Set("X-Delivery-Id", d.DeliveryID)
	h.Set("X-Chronary-Event-Type", d.EventType)
	if d.Test {
		h.Set("X-Chronary-Test", "1")
	}
	if d.ReplayOf != "" {
		h.Set("X-Chronary-Replay-Of", d.ReplayOf)
	}
	return h
}

func TestVerifySignatureValid(t *testing.T) {
	payload := []byte(`{"type":"event.created","data":{}}`)
	secret := "whsec_test_secret"
	now := time.Now()
	headers := deliveryHeaders(payload, secret, now.Unix(), liveDelivery)

	err := VerifySignatureWithOptions(payload, headers, secret, withNow(func() time.Time { return now }))
	if err != nil {
		t.Fatalf("expected valid signature, got: %v", err)
	}
}

func TestVerifySignatureInvalid(t *testing.T) {
	payload := []byte(`{"type":"event.created","data":{}}`)
	secret := "whsec_test_secret"
	now := time.Now()
	headers := deliveryHeaders(payload, secret, now.Unix(), liveDelivery)
	headers.Set("X-Signature-V2", "sha256=invalidsignature")

	err := VerifySignatureWithOptions(payload, headers, secret, withNow(func() time.Time { return now }))
	if err == nil {
		t.Fatal("expected signature mismatch error")
	}
}

func TestVerifySignatureExpired(t *testing.T) {
	payload := []byte(`{"type":"event.created","data":{}}`)
	secret := "whsec_test_secret"
	headers := deliveryHeaders(payload, secret, time.Now().Add(-10*time.Minute).Unix(), liveDelivery)

	err := VerifySignature(payload, headers, secret)
	if err == nil {
		t.Fatal("expected timestamp too old error")
	}
}

func TestVerifySignatureCustomTolerance(t *testing.T) {
	payload := []byte(`{"type":"event.created","data":{}}`)
	secret := "whsec_test_secret"
	headers := deliveryHeaders(payload, secret, time.Now().Add(-10*time.Minute).Unix(), liveDelivery)

	err := VerifySignatureWithOptions(payload, headers, secret, WithTolerance(15*time.Minute))
	if err != nil {
		t.Fatalf("expected valid with 15m tolerance, got: %v", err)
	}
}

func TestVerifySignatureMissingHeaders(t *testing.T) {
	payload := []byte(`{}`)
	secret := "whsec_test_secret"

	err := VerifySignature(payload, http.Header{}, secret)
	if err == nil {
		t.Fatal("expected error for missing headers")
	}

	now := time.Now()
	for _, name := range []string{"X-Signature-V2", "X-Timestamp", "X-Delivery-Id", "X-Chronary-Event-Type"} {
		headers := deliveryHeaders(payload, secret, now.Unix(), liveDelivery)
		headers.Del(name)
		err := VerifySignatureWithOptions(payload, headers, secret, withNow(func() time.Time { return now }))
		if err == nil || !strings.Contains(err.Error(), "missing "+name) {
			t.Errorf("without %s: expected a missing-header error, got %v", name, err)
		}
	}
}

func TestVerifySignatureInvalidTimestamp(t *testing.T) {
	payload := []byte(`{}`)
	headers := deliveryHeaders(payload, "whsec_test_secret", time.Now().Unix(), liveDelivery)
	headers.Set("X-Timestamp", "nope")

	err := VerifySignature(payload, headers, "whsec_test_secret")
	if err == nil || !strings.Contains(err.Error(), "invalid X-Timestamp") {
		t.Fatalf("expected invalid X-Timestamp error, got %v", err)
	}
}

func TestVerifySignatureTamperedBody(t *testing.T) {
	payload := []byte(`{"type":"event.created","data":{}}`)
	secret := "whsec_test_secret"
	now := time.Now()
	headers := deliveryHeaders(payload, secret, now.Unix(), liveDelivery)

	tampered := []byte(`{"type":"event.created","data":{"injected":true}}`)
	err := VerifySignatureWithOptions(tampered, headers, secret, withNow(func() time.Time { return now }))
	if err == nil {
		t.Fatal("expected signature mismatch for tampered body")
	}
}

func TestConstructEvent(t *testing.T) {
	payload := []byte(`{"event":{"id":"evt_1"}}`)
	secret := "whsec_test_secret"
	now := time.Now()
	headers := deliveryHeaders(payload, secret, now.Unix(), SignedDelivery{
		DeliveryID: "whd_2", EventType: "event.created", Test: true, ReplayOf: "whd_1",
	})

	event, err := ConstructEvent(payload, headers, secret, withNow(func() time.Time { return now }))
	if err != nil {
		t.Fatal(err)
	}
	if event.Type != "event.created" {
		t.Errorf("expected event.created, got %s", event.Type)
	}
	if event.DeliveryID != "whd_2" || !event.Test || event.ReplayOf != "whd_1" {
		t.Errorf("unexpected signed metadata: %+v", event)
	}
	eventData, ok := event.Data["event"].(map[string]interface{})
	if !ok {
		t.Fatalf("expected event payload object, got %T", event.Data["event"])
	}
	if eventData["id"] != "evt_1" {
		t.Errorf("expected evt_1, got %v", eventData["id"])
	}
}

func TestConstructEventMissingEventType(t *testing.T) {
	payload := []byte(`{"event":{"id":"evt_1"}}`)
	secret := "whsec_test_secret"
	now := time.Now()
	headers := deliveryHeaders(payload, secret, now.Unix(), liveDelivery)
	headers.Del("X-Chronary-Event-Type")

	_, err := ConstructEvent(payload, headers, secret, withNow(func() time.Time { return now }))
	if err == nil {
		t.Fatal("expected missing X-Chronary-Event-Type error")
	}
}

func TestConstructEventInvalidJSON(t *testing.T) {
	payload := []byte(`not json`)
	headers := SignedHeaders(payload, "whsec_test_secret")

	_, err := ConstructEvent(payload, headers, "whsec_test_secret")
	if err == nil || !strings.Contains(err.Error(), "parsing webhook payload") {
		t.Fatalf("expected a JSON parse error, got %v", err)
	}
}

func TestSignedHeaders(t *testing.T) {
	payload := []byte(`{"test":true}`)
	secret := "test_secret"
	headers := SignedHeaders(payload, secret)

	for _, name := range []string{"X-Signature-V2", "X-Signature", "X-Timestamp"} {
		if headers.Get(name) == "" {
			t.Errorf("expected %s header", name)
		}
	}
	if !strings.HasPrefix(headers.Get("X-Delivery-Id"), "whd_test_") {
		t.Errorf("expected a generated delivery id, got %q", headers.Get("X-Delivery-Id"))
	}
	if headers.Get("X-Chronary-Event-Type") != "event.created" {
		t.Errorf("expected the default event type, got %q", headers.Get("X-Chronary-Event-Type"))
	}

	err := VerifySignature(payload, headers, secret)
	if err != nil {
		t.Fatalf("SignedHeaders should produce valid signatures: %v", err)
	}
}

func TestSignedHeadersWithEventType(t *testing.T) {
	payload := []byte(`{"event":{"id":"evt_1"}}`)
	secret := "test_secret"
	headers := SignedHeadersWithEventType(payload, secret, "event.updated")

	if headers.Get("X-Chronary-Event-Type") != "event.updated" {
		t.Error("expected X-Chronary-Event-Type header")
	}

	event, err := ConstructEvent(payload, headers, secret)
	if err != nil {
		t.Fatalf("ConstructEvent should accept signed headers: %v", err)
	}
	if event.Type != "event.updated" {
		t.Errorf("expected event.updated, got %s", event.Type)
	}
}

func TestSignedHeadersFor(t *testing.T) {
	payload := []byte(`{"event":{"id":"evt_1"}}`)
	headers := SignedHeadersFor(payload, "test_secret", SignedDelivery{
		DeliveryID: "whd_9", EventType: "event.deleted", Test: true, ReplayOf: "whd_8",
	})

	event, err := ConstructEvent(payload, headers, "test_secret")
	if err != nil {
		t.Fatalf("ConstructEvent should accept SignedHeadersFor: %v", err)
	}
	if event.Type != "event.deleted" || event.DeliveryID != "whd_9" || !event.Test || event.ReplayOf != "whd_8" {
		t.Errorf("unexpected event: %+v", event)
	}
}

// --- BH-069 / #2569: the delivery metadata is covered by the signature ---

func expectRejected(t *testing.T, name string, payload []byte, headers http.Header, secret string, now time.Time) {
	t.Helper()
	at := withNow(func() time.Time { return now })
	if err := VerifySignatureWithOptions(payload, headers, secret, at); err == nil {
		t.Errorf("%s: VerifySignature accepted a forged delivery", name)
	}
	if event, err := ConstructEvent(payload, headers, secret, at); err == nil {
		t.Errorf("%s: ConstructEvent accepted a forged delivery: %+v", name, event)
	}
}

func TestSignatureCoversDeliveryMetadata(t *testing.T) {
	secret := "whsec_bughunt_fixture"
	payload := []byte(`{"calendar_id":"cal_1","event":{"id":"evt_1","title":"x"}}`)
	now := time.Now()
	ts := now.Unix()
	original := SignedDelivery{DeliveryID: "whd_original", EventType: "event.created"}

	if _, err := ConstructEvent(payload, deliveryHeaders(payload, secret, ts, original), secret, withNow(func() time.Time { return now })); err != nil {
		t.Fatalf("genuine delivery rejected: %v", err)
	}

	relabelled := deliveryHeaders(payload, secret, ts, original)
	relabelled.Set("X-Chronary-Event-Type", "event.updated")
	expectRejected(t, "relabelled event type", payload, relabelled, secret, now)

	holdBody := []byte(`{"calendar_id":"cal_1","event_id":"evt_1"}`)
	hold := deliveryHeaders(holdBody, secret, ts, SignedDelivery{DeliveryID: "whd_h", EventType: "event.hold_released"})
	hold.Set("X-Chronary-Event-Type", "event.deleted")
	expectRejected(t, "hold_released relabelled as deleted", holdBody, hold, secret, now)

	freshID := deliveryHeaders(payload, secret, ts, original)
	freshID.Set("X-Delivery-Id", "whd_attacker_fresh")
	expectRejected(t, "swapped delivery id", payload, freshID, secret, now)

	testSend := deliveryHeaders(payload, secret, ts, SignedDelivery{DeliveryID: "whd_t", EventType: "event.created", Test: true})
	testSend.Del("X-Chronary-Test")
	expectRejected(t, "stripped test flag", payload, testSend, secret, now)

	addedTest := deliveryHeaders(payload, secret, ts, original)
	addedTest.Set("X-Chronary-Test", "1")
	expectRejected(t, "added test flag", payload, addedTest, secret, now)
	addedTest.Set("X-Chronary-Test", "true")
	expectRejected(t, "non-1 test flag", payload, addedTest, secret, now)

	replay := SignedDelivery{DeliveryID: "whd_r", EventType: "event.created", ReplayOf: "whd_original"}
	changedReplay := deliveryHeaders(payload, secret, ts, replay)
	changedReplay.Set("X-Chronary-Replay-Of", "whd_other")
	expectRejected(t, "changed replay-of", payload, changedReplay, secret, now)
	strippedReplay := deliveryHeaders(payload, secret, ts, replay)
	strippedReplay.Del("X-Chronary-Replay-Of")
	expectRejected(t, "stripped replay-of", payload, strippedReplay, secret, now)
	addedReplay := deliveryHeaders(payload, secret, ts, original)
	addedReplay.Set("X-Chronary-Replay-Of", "whd_original")
	expectRejected(t, "added replay-of", payload, addedReplay, secret, now)

	v1Only := deliveryHeaders(payload, secret, ts, original)
	v1Only.Del("X-Signature-V2")
	expectRejected(t, "v1 only (no fallback)", payload, v1Only, secret, now)
	v1Only.Set("X-Chronary-Event-Type", "event.updated")
	expectRejected(t, "v1 only, relabelled", payload, v1Only, secret, now)

	v1InV2 := deliveryHeaders(payload, secret, ts, original)
	v1InV2.Set("X-Signature-V2", v1InV2.Get("X-Signature"))
	expectRejected(t, "v1 signature in the v2 header", payload, v1InV2, secret, now)

	crlf := deliveryHeaders(payload, secret, ts, original)
	crlf["X-Delivery-Id"] = []string{"whd_original\nevent.created"}
	expectRejected(t, "LF in a header value", payload, crlf, secret, now)
	// Rejected by the CR/LF guard itself, not only by the signature mismatch.
	for _, name := range []string{"X-Timestamp", "X-Delivery-Id", "X-Chronary-Event-Type", "X-Chronary-Replay-Of"} {
		h := deliveryHeaders(payload, secret, ts, original)
		h[name] = []string{h.Get(name) + "\r"}
		err := VerifySignatureWithOptions(payload, h, secret, withNow(func() time.Time { return now }))
		if err == nil || !strings.Contains(err.Error(), "CR or LF") {
			t.Errorf("CR in %s: expected the CR or LF error, got %v", name, err)
		}
	}
}

func TestSignatureV2KnownAnswerVectors(t *testing.T) {
	// The same vectors as the API and the TS and Python SDK tests.
	payload := []byte(`{"calendar_id":"cal_1","event":{"id":"evt_1"}}`)
	const ts = int64(1745784205)
	cases := []struct {
		d    SignedDelivery
		want string
	}{
		{SignedDelivery{DeliveryID: "whd_vector1", EventType: "event.created"}, "sha256=5ee321bd30e58aef4e26fb8cb10850982d8d11531e767983737f346908a88560"},
		{SignedDelivery{DeliveryID: "whd_vector1", EventType: "event.created", Test: true}, "sha256=d81393f416bf1894fb8ddd17eb2527fd74104cf57f9e5b1d5b7eeb6d2c00048d"},
		{SignedDelivery{DeliveryID: "whd_vector2", EventType: "event.created", ReplayOf: "whd_vector1"}, "sha256=cc34c8c295b86e898ce7675a34c9ca16ab1e843459ebcaa74d56dae5b8e0df88"},
	}
	for _, c := range cases {
		if got := ComputeSignatureV2(payload, "whsec_v2_test_vector", ts, c.d); got != c.want {
			t.Errorf("%+v: got %s, want %s", c.d, got, c.want)
		}
		headers := deliveryHeaders(payload, "whsec_v2_test_vector", ts, c.d)
		if err := VerifySignatureWithOptions(payload, headers, "whsec_v2_test_vector", withNow(func() time.Time { return time.Unix(ts, 0) })); err != nil {
			t.Errorf("%+v: %v", c.d, err)
		}
	}
}

// --- Secret rotation (#594) ---

func rotatedHeaders(payload []byte, oldSecret, newSecret string, ts int64) http.Header {
	h := deliveryHeaders(payload, newSecret, ts, liveDelivery)
	h.Set("X-Signature-Previous", ComputeSignature(payload, oldSecret, ts))
	h.Set("X-Signature-V2-Previous", ComputeSignatureV2(payload, oldSecret, ts, liveDelivery))
	return h
}

func TestVerifySignatureRotationAcceptsEitherSecret(t *testing.T) {
	payload := []byte(`{"type":"event.created","data":{}}`)
	now := time.Now()
	headers := rotatedHeaders(payload, "whsec_old", "whsec_new", now.Unix())
	at := withNow(func() time.Time { return now })

	if err := VerifySignatureWithOptions(payload, headers, "whsec_new", at); err != nil {
		t.Fatalf("new secret via X-Signature-V2: %v", err)
	}
	if err := VerifySignatureWithOptions(payload, headers, "whsec_old", at); err != nil {
		t.Fatalf("old secret via X-Signature-V2-Previous: %v", err)
	}
	if err := VerifySignatureWithOptions(payload, headers, "whsec_other", at); err == nil {
		t.Fatal("expected a mismatch when neither header matches")
	}
	relabelled := rotatedHeaders(payload, "whsec_old", "whsec_new", now.Unix())
	relabelled.Set("X-Chronary-Event-Type", "event.updated")
	expectRejected(t, "rotation twin, relabelled", payload, relabelled, "whsec_old", now)
}

func TestVerifySignatureRotationAfterWindow(t *testing.T) {
	payload := []byte(`{"type":"event.created","data":{}}`)
	now := time.Now()
	headers := rotatedHeaders(payload, "whsec_old", "whsec_new", now.Unix())
	headers.Del("X-Signature-V2-Previous")

	// X-Signature-Previous (v1) is still there and valid, and is ignored.
	err := VerifySignatureWithOptions(payload, headers, "whsec_old", withNow(func() time.Time { return now }))
	if err == nil {
		t.Fatal("expected the old secret to fail once X-Signature-V2-Previous is gone")
	}
}

func TestVerifySignatureRotationStillRequiresXSignatureV2(t *testing.T) {
	payload := []byte(`{"type":"event.created","data":{}}`)
	now := time.Now()
	headers := rotatedHeaders(payload, "whsec_old", "whsec_new", now.Unix())
	headers.Del("X-Signature-V2")

	err := VerifySignatureWithOptions(payload, headers, "whsec_old", withNow(func() time.Time { return now }))
	if err == nil || !strings.HasPrefix(err.Error(), "chronary: missing X-Signature-V2 header") {
		t.Fatalf("expected missing X-Signature-V2 error, got %v", err)
	}
}

func TestVerifySignatureRotationKeepsTimestampCheck(t *testing.T) {
	payload := []byte(`{"type":"event.created","data":{}}`)
	now := time.Now()
	headers := rotatedHeaders(payload, "whsec_old", "whsec_new", now.Add(-time.Hour).Unix())

	err := VerifySignatureWithOptions(payload, headers, "whsec_old", withNow(func() time.Time { return now }))
	if err == nil {
		t.Fatal("expected a stale timestamp to fail even with a valid previous signature")
	}
}
