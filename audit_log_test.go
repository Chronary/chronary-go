package chronary_test

import (
	"context"
	"net/http"
	"testing"

	chronary "github.com/Chronary/chronary-go"
	"github.com/Chronary/chronary-go/internal/testutil"
)

func TestAuditLogListActorFilters(t *testing.T) {
	want := map[string]string{
		"action":       "agent_key.create",
		"actor_type":   "session",
		"actor_id":     "org_1",
		"request_id":   "req_1",
		"resource":     "cal_",
		"status_class": "2xx",
		"category":     "keys",
		"limit":        "10",
	}
	client, _ := testutil.Setup(t, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		testutil.AssertMethod(t, r, http.MethodGet)
		testutil.AssertPath(t, r, "/v1/audit-log")
		q := r.URL.Query()
		if len(q) != len(want) {
			t.Errorf("query = %v, want %v", q, want)
		}
		for k, v := range want {
			if got := q.Get(k); got != v {
				t.Errorf("query %s = %q, want %q", k, got, v)
			}
		}
		testutil.RespondJSON(w, 200, map[string]interface{}{
			"data": []interface{}{
				map[string]interface{}{
					"id": "aud_1", "action": "agent_key.create", "actor_key_prefix": nil,
					"actor_type": "session", "actor_id": "org_1", "actor_label": "owner@example.com",
					"agent_id": nil, "resource": "", "ip": nil, "status": 201, "method": "POST",
					"path": "/v1/auth/keys/scoped", "duration_ms": 12, "request_id": "req_1",
					"created_at": "2026-10-30T00:00:00.000Z",
				},
			},
			"pagination":     map[string]interface{}{"next_cursor": nil},
			"retention_days": 90,
			"range_clamped":  false,
		})
	}))

	resp, err := client.AuditLog.List(context.Background(), &chronary.ListAuditLogParams{
		Action:      "agent_key.create",
		ActorType:   "session",
		ActorID:     "org_1",
		RequestID:   "req_1",
		Resource:    "cal_",
		StatusClass: "2xx",
		Category:    "keys",
		Limit:       10,
	})
	if err != nil {
		t.Fatal(err)
	}
	e := resp.Data[0]
	if e.ActorType == nil || *e.ActorType != "session" || e.ActorID == nil || *e.ActorID != "org_1" ||
		e.ActorLabel == nil || *e.ActorLabel != "owner@example.com" || e.ActorKeyPrefix != nil {
		t.Fatalf("unexpected actor fields: %#v", e)
	}
}
