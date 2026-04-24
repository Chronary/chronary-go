package chronary

import (
	"net/http"
	"strconv"
	"testing"
	"time"
)

func TestVerifySignatureValid(t *testing.T) {
	payload := []byte(`{"type":"event.created","data":{}}`)
	secret := "whsec_test_secret"
	now := time.Now()
	ts := now.Unix()

	sig := ComputeSignature(payload, secret, ts)
	headers := http.Header{}
	headers.Set("X-Signature", sig)
	headers.Set("X-Timestamp", strconv.FormatInt(ts, 10))

	err := VerifySignatureWithOptions(payload, headers, secret, withNow(func() time.Time { return now }))
	if err != nil {
		t.Fatalf("expected valid signature, got: %v", err)
	}
}

func TestVerifySignatureInvalid(t *testing.T) {
	payload := []byte(`{"type":"event.created","data":{}}`)
	secret := "whsec_test_secret"
	now := time.Now()
	ts := now.Unix()

	headers := http.Header{}
	headers.Set("X-Signature", "sha256=invalidsignature")
	headers.Set("X-Timestamp", strconv.FormatInt(ts, 10))

	err := VerifySignatureWithOptions(payload, headers, secret, withNow(func() time.Time { return now }))
	if err == nil {
		t.Fatal("expected signature mismatch error")
	}
}

func TestVerifySignatureExpired(t *testing.T) {
	payload := []byte(`{"type":"event.created","data":{}}`)
	secret := "whsec_test_secret"
	oldTime := time.Now().Add(-10 * time.Minute)
	ts := oldTime.Unix()

	sig := ComputeSignature(payload, secret, ts)
	headers := http.Header{}
	headers.Set("X-Signature", sig)
	headers.Set("X-Timestamp", strconv.FormatInt(ts, 10))

	err := VerifySignature(payload, headers, secret)
	if err == nil {
		t.Fatal("expected timestamp too old error")
	}
}

func TestVerifySignatureCustomTolerance(t *testing.T) {
	payload := []byte(`{"type":"event.created","data":{}}`)
	secret := "whsec_test_secret"
	oldTime := time.Now().Add(-10 * time.Minute)
	ts := oldTime.Unix()

	sig := ComputeSignature(payload, secret, ts)
	headers := http.Header{}
	headers.Set("X-Signature", sig)
	headers.Set("X-Timestamp", strconv.FormatInt(ts, 10))

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
}

func TestVerifySignatureTamperedBody(t *testing.T) {
	payload := []byte(`{"type":"event.created","data":{}}`)
	secret := "whsec_test_secret"
	now := time.Now()
	ts := now.Unix()

	sig := ComputeSignature(payload, secret, ts)
	headers := http.Header{}
	headers.Set("X-Signature", sig)
	headers.Set("X-Timestamp", strconv.FormatInt(ts, 10))

	tampered := []byte(`{"type":"event.created","data":{"injected":true}}`)
	err := VerifySignatureWithOptions(tampered, headers, secret, withNow(func() time.Time { return now }))
	if err == nil {
		t.Fatal("expected signature mismatch for tampered body")
	}
}

func TestConstructEvent(t *testing.T) {
	payload := []byte(`{"type":"event.created","data":{"id":"evt_1"}}`)
	secret := "whsec_test_secret"
	now := time.Now()
	ts := now.Unix()

	sig := ComputeSignature(payload, secret, ts)
	headers := http.Header{}
	headers.Set("X-Signature", sig)
	headers.Set("X-Timestamp", strconv.FormatInt(ts, 10))

	event, err := ConstructEvent(payload, headers, secret, withNow(func() time.Time { return now }))
	if err != nil {
		t.Fatal(err)
	}
	if event.Type != "event.created" {
		t.Errorf("expected event.created, got %s", event.Type)
	}
	if event.Data["id"] != "evt_1" {
		t.Errorf("expected evt_1, got %v", event.Data["id"])
	}
}

func TestSignedHeaders(t *testing.T) {
	payload := []byte(`{"test":true}`)
	secret := "test_secret"
	headers := SignedHeaders(payload, secret)

	if headers.Get("X-Signature") == "" {
		t.Error("expected X-Signature header")
	}
	if headers.Get("X-Timestamp") == "" {
		t.Error("expected X-Timestamp header")
	}

	err := VerifySignature(payload, headers, secret)
	if err != nil {
		t.Fatalf("SignedHeaders should produce valid signatures: %v", err)
	}
}
