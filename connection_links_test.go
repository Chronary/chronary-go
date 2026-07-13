package chronary_test

import (
	"context"
	"encoding/json"
	"net/http"
	"testing"

	chronary "github.com/Chronary/chronary-go"
	"github.com/Chronary/chronary-go/internal/testutil"
)

func TestConnectionLinksLifecycle(t *testing.T) {
	step := 0
	client, _ := testutil.Setup(t, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		step++
		switch step {
		case 1:
			testutil.AssertMethod(t, r, http.MethodPost)
			testutil.AssertPath(t, r, "/v1/calendars/cal_1/connection-links")
			var body chronary.CreateConnectionLinkParams
			if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
				t.Fatal(err)
			}
			if len(body.Capabilities) != 1 || body.Capabilities[0] != "availability" {
				t.Fatalf("unexpected capabilities: %#v", body.Capabilities)
			}
			testutil.RespondJSON(w, 201, map[string]interface{}{"id": "csl_1", "calendar_id": "cal_1", "setup_url": "https://setup.test", "status": "awaiting_human", "expires_at": "2026-07-12T20:00:00Z"})
		case 2:
			testutil.AssertMethod(t, r, http.MethodGet)
			testutil.AssertPath(t, r, "/v1/connection-links/csl_1")
			testutil.RespondJSON(w, 200, map[string]interface{}{"id": "csl_1", "calendar_id": "cal_1", "status": "completed", "expires_at": "2026-07-12T20:00:00Z", "connection_id": "con_1"})
		case 3:
			testutil.AssertMethod(t, r, http.MethodDelete)
			testutil.AssertPath(t, r, "/v1/connection-links/csl_1")
			w.WriteHeader(http.StatusNoContent)
		}
	}))

	link, err := client.ConnectionLinks.Create(context.Background(), "cal_1", &chronary.CreateConnectionLinkParams{Capabilities: []string{"availability"}})
	if err != nil || link.ID != "csl_1" || link.SetupURL == nil {
		t.Fatalf("create: %#v, %v", link, err)
	}
	status, err := client.ConnectionLinks.Get(context.Background(), link.ID)
	if err != nil || status.ConnectionID == nil || *status.ConnectionID != "con_1" {
		t.Fatalf("get: %#v, %v", status, err)
	}
	if err := client.ConnectionLinks.Cancel(context.Background(), link.ID); err != nil {
		t.Fatal(err)
	}
}

func TestAvailabilityAllowStaleAndHealth(t *testing.T) {
	client, _ := testutil.Setup(t, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Query().Get("allow_stale") != "true" {
			t.Errorf("allow_stale not serialized")
		}
		testutil.RespondJSON(w, 200, map[string]interface{}{
			"agents": []string{"agt_1"}, "slots": []interface{}{}, "availability_state": "stale",
			"sources":  map[string]interface{}{"chronary": "current", "external": "stale"},
			"warnings": []map[string]string{{"code": "availability_incomplete", "message": "stale"}},
		})
	}))
	resp, err := client.Availability.ForAgent(context.Background(), "agt_1", &chronary.AvailabilityParams{Start: "2026-07-12T00:00:00Z", End: "2026-07-13T00:00:00Z", AllowStale: chronary.Bool(true)})
	if err != nil || resp.AvailabilityState != "stale" || len(resp.Warnings) != 1 {
		t.Fatalf("response: %#v, %v", resp, err)
	}
}
