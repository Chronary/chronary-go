// Package testutil provides shared test helpers for the Chronary Go SDK.
package testutil

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	chronary "github.com/Chronary/chronary-go"
)

// Setup creates an httptest.Server and a configured Client pointing at it.
func Setup(t *testing.T, handler http.Handler) (*chronary.Client, *httptest.Server) {
	t.Helper()
	srv := httptest.NewServer(handler)
	t.Cleanup(srv.Close)

	client, err := chronary.NewClient(
		chronary.WithAPIKey("test_key"),
		chronary.WithBaseURL(srv.URL),
		chronary.WithMaxRetries(0),
	)
	if err != nil {
		t.Fatal(err)
	}
	return client, srv
}

// RespondJSON writes a JSON response with the given status code.
func RespondJSON(w http.ResponseWriter, status int, body interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(body)
}

// RespondError writes a Chronary API error response.
func RespondError(w http.ResponseWriter, status int, errType, message string) {
	RespondJSON(w, status, map[string]interface{}{
		"error": map[string]interface{}{
			"type":    errType,
			"message": message,
		},
	})
}

// AssertMethod checks that the request used the expected HTTP method.
func AssertMethod(t *testing.T, r *http.Request, expected string) {
	t.Helper()
	if r.Method != expected {
		t.Errorf("expected method %s, got %s", expected, r.Method)
	}
}

// AssertPath checks that the request path matches.
func AssertPath(t *testing.T, r *http.Request, expected string) {
	t.Helper()
	if r.URL.Path != expected {
		t.Errorf("expected path %s, got %s", expected, r.URL.Path)
	}
}

// AssertHasHeader checks a header is present with a specific value.
func AssertHasHeader(t *testing.T, r *http.Request, key, value string) {
	t.Helper()
	if got := r.Header.Get(key); got != value {
		t.Errorf("expected header %s=%s, got %s", key, value, got)
	}
}

// AssertHasAuth checks the Authorization header is present.
func AssertHasAuth(t *testing.T, r *http.Request) {
	t.Helper()
	if got := r.Header.Get("Authorization"); got != "Bearer test_key" {
		t.Errorf("expected Authorization: Bearer test_key, got %s", got)
	}
}

// AssertHasIdempotencyKey checks the Idempotency-Key header is present and non-empty.
func AssertHasIdempotencyKey(t *testing.T, r *http.Request) {
	t.Helper()
	if got := r.Header.Get("Idempotency-Key"); got == "" {
		t.Error("expected Idempotency-Key header to be present")
	}
}
