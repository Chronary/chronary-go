package chronary_test

import (
	"context"
	"net/http"
	"testing"

	"github.com/Chronary/chronary-go/internal/testutil"
)

func icalSubscriptionBody(status string) map[string]interface{} {
	return map[string]interface{}{
		"id": "ics_1", "agent_id": "agt_1", "calendar_id": "cal_1",
		"url": "https://example.com/feed.ics", "label": nil, "status": status,
		"last_synced_at": nil, "last_error": nil, "created_at": "2026-11-01T00:00:00Z",
	}
}

func TestICalSubscriptionPause(t *testing.T) {
	client, _ := testutil.Setup(t, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		testutil.AssertMethod(t, r, "POST")
		testutil.AssertPath(t, r, "/v1/ical-subscriptions/ics_1/pause")
		testutil.RespondJSON(w, 200, icalSubscriptionBody("paused"))
	}))

	sub, err := client.ICalSubscriptions.Pause(context.Background(), "ics_1")
	if err != nil {
		t.Fatal(err)
	}
	if sub.Status != "paused" {
		t.Errorf("expected status paused, got %s", sub.Status)
	}
}

func TestICalSubscriptionResume(t *testing.T) {
	client, _ := testutil.Setup(t, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		testutil.AssertMethod(t, r, "POST")
		testutil.AssertPath(t, r, "/v1/ical-subscriptions/ics_1/resume")
		testutil.RespondJSON(w, 200, icalSubscriptionBody("active"))
	}))

	sub, err := client.ICalSubscriptions.Resume(context.Background(), "ics_1")
	if err != nil {
		t.Fatal(err)
	}
	if sub.Status != "active" {
		t.Errorf("expected status active, got %s", sub.Status)
	}
}

func TestICalSubscriptionPauseError(t *testing.T) {
	client, _ := testutil.Setup(t, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		testutil.RespondError(w, 404, "not_found", "Subscription not found")
	}))

	if _, err := client.ICalSubscriptions.Pause(context.Background(), "ics_missing"); err == nil {
		t.Fatal("expected an error")
	}
	if _, err := client.ICalSubscriptions.Resume(context.Background(), "ics_missing"); err == nil {
		t.Fatal("expected an error")
	}
}
