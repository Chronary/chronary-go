package chronary_test

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"sync/atomic"
	"testing"

	chronary "github.com/Chronary/chronary-go"
)

func TestDoSetsHeaders(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("Authorization") != "Bearer test_key" {
			t.Error("missing or wrong Authorization header")
		}
		if r.Header.Get("User-Agent") == "" {
			t.Error("missing User-Agent header")
		}
		if r.Header.Get("X-Chronary-SDK-Version") != chronary.Version {
			t.Error("missing or wrong SDK version header")
		}
		w.WriteHeader(200)
		json.NewEncoder(w).Encode(map[string]string{"id": "agt_1"})
	}))
	defer srv.Close()

	client, _ := chronary.NewClient(
		chronary.WithAPIKey("test_key"),
		chronary.WithBaseURL(srv.URL),
	)
	_, err := client.Agents.Get(context.Background(), "agt_1")
	if err != nil {
		t.Fatal(err)
	}
}

func TestDoIdempotencyKeyOnPost(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost && r.Header.Get("Idempotency-Key") == "" {
			t.Error("POST should have Idempotency-Key header")
		}
		w.WriteHeader(200)
		json.NewEncoder(w).Encode(map[string]string{"id": "agt_1", "name": "test", "type": "ai", "status": "active"})
	}))
	defer srv.Close()

	client, _ := chronary.NewClient(
		chronary.WithAPIKey("test_key"),
		chronary.WithBaseURL(srv.URL),
		chronary.WithMaxRetries(0),
	)
	_, err := client.Agents.Create(context.Background(), &chronary.CreateAgentParams{
		Name: "test",
		Type: chronary.AgentTypeAI,
	})
	if err != nil {
		t.Fatal(err)
	}
}

func TestDoReusesAutoIdempotencyKeyOnRetry(t *testing.T) {
	var attempts int32
	var firstKey string
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		key := r.Header.Get("Idempotency-Key")
		if key == "" {
			t.Error("POST should have Idempotency-Key header")
		}
		n := atomic.AddInt32(&attempts, 1)
		if n == 1 {
			firstKey = key
			w.WriteHeader(500)
			w.Write([]byte(`{"error":{"type":"server_error","message":"oops"}}`))
			return
		}
		if key != firstKey {
			t.Errorf("expected retry to reuse idempotency key %q, got %q", firstKey, key)
		}
		w.WriteHeader(200)
		json.NewEncoder(w).Encode(map[string]string{"id": "agt_1", "name": "test", "type": "ai", "status": "active"})
	}))
	defer srv.Close()

	client, _ := chronary.NewClient(
		chronary.WithAPIKey("test_key"),
		chronary.WithBaseURL(srv.URL),
		chronary.WithMaxRetries(1),
	)
	_, err := client.Agents.Create(context.Background(), &chronary.CreateAgentParams{
		Name: "test",
		Type: chronary.AgentTypeAI,
	})
	if err != nil {
		t.Fatal(err)
	}
	if atomic.LoadInt32(&attempts) != 2 {
		t.Errorf("expected 2 attempts, got %d", atomic.LoadInt32(&attempts))
	}
}

func TestDoNoIdempotencyKeyOnGet(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet && r.Header.Get("Idempotency-Key") != "" {
			t.Error("GET should not have Idempotency-Key header")
		}
		w.WriteHeader(200)
		json.NewEncoder(w).Encode(map[string]string{"id": "agt_1"})
	}))
	defer srv.Close()

	client, _ := chronary.NewClient(
		chronary.WithAPIKey("test_key"),
		chronary.WithBaseURL(srv.URL),
	)
	_, err := client.Agents.Get(context.Background(), "agt_1")
	if err != nil {
		t.Fatal(err)
	}
}

func TestDoRetries(t *testing.T) {
	var attempts int32
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		n := atomic.AddInt32(&attempts, 1)
		if n <= 2 {
			w.WriteHeader(500)
			w.Write([]byte(`{"error":{"type":"server_error","message":"oops"}}`))
			return
		}
		w.WriteHeader(200)
		json.NewEncoder(w).Encode(map[string]string{"id": "agt_1"})
	}))
	defer srv.Close()

	client, _ := chronary.NewClient(
		chronary.WithAPIKey("test_key"),
		chronary.WithBaseURL(srv.URL),
		chronary.WithMaxRetries(2),
	)
	_, err := client.Agents.Get(context.Background(), "agt_1")
	if err != nil {
		t.Fatal(err)
	}
	if atomic.LoadInt32(&attempts) != 3 {
		t.Errorf("expected 3 attempts, got %d", atomic.LoadInt32(&attempts))
	}
}

func TestDoNoRetryOn4xx(t *testing.T) {
	var attempts int32
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		atomic.AddInt32(&attempts, 1)
		w.WriteHeader(404)
		w.Write([]byte(`{"error":{"type":"not_found","message":"not found"}}`))
	}))
	defer srv.Close()

	client, _ := chronary.NewClient(
		chronary.WithAPIKey("test_key"),
		chronary.WithBaseURL(srv.URL),
		chronary.WithMaxRetries(2),
	)
	_, err := client.Agents.Get(context.Background(), "agt_1")
	if err == nil {
		t.Fatal("expected error")
	}
	if !errors.Is(err, chronary.ErrNotFound) {
		t.Errorf("expected ErrNotFound, got %v", err)
	}
	if atomic.LoadInt32(&attempts) != 1 {
		t.Errorf("expected 1 attempt (no retry on 404), got %d", atomic.LoadInt32(&attempts))
	}
}

func TestDoErrorParsing(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(422)
		w.Write([]byte(`{"error":{"type":"validation_error","message":"name is required","request_id":"req_abc"}}`))
	}))
	defer srv.Close()

	client, _ := chronary.NewClient(
		chronary.WithAPIKey("test_key"),
		chronary.WithBaseURL(srv.URL),
		chronary.WithMaxRetries(0),
	)
	_, err := client.Agents.Get(context.Background(), "agt_1")
	if !errors.Is(err, chronary.ErrValidation) {
		t.Errorf("expected ErrValidation, got %v", err)
	}
	var chronErr *chronary.Error
	if errors.As(err, &chronErr) {
		if chronErr.RequestID != "req_abc" {
			t.Errorf("expected request_id req_abc, got %s", chronErr.RequestID)
		}
	}
}
