// Surface tests assert that every public client method hits the canonical API
// path. They guard against future regressions like the GET /v1/events/{id}
// bug that shipped in earlier versions before the route was mounted under
// /v1/calendars/{cal_id}/events/{id}.

package chronary_test

import (
	"context"
	"net/http"
	"testing"

	chronary "github.com/Chronary/chronary-go"
	"github.com/Chronary/chronary-go/internal/testutil"
)

// --- Event path bug fix + confirm/release ---

func TestEventGet_UsesCalendarScopedPath(t *testing.T) {
	client, _ := testutil.Setup(t, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		testutil.AssertMethod(t, r, "GET")
		testutil.AssertPath(t, r, "/v1/calendars/cal_1/events/evt_1")
		testutil.RespondJSON(w, 200, map[string]interface{}{
			"id": "evt_1", "calendar_id": "cal_1", "title": "Standup",
			"start_time": "2026-04-14T09:00:00Z", "end_time": "2026-04-14T09:30:00Z",
			"all_day": false, "status": "confirmed", "source": "internal",
			"created_at": "2026-04-14T00:00:00Z", "updated_at": "2026-04-14T00:00:00Z",
		})
	}))

	event, err := client.Events.Get(context.Background(), "cal_1", "evt_1")
	if err != nil {
		t.Fatal(err)
	}
	if event.ID != "evt_1" {
		t.Errorf("expected evt_1, got %s", event.ID)
	}
}

func TestEventUpdate_UsesCalendarScopedPath(t *testing.T) {
	client, _ := testutil.Setup(t, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		testutil.AssertMethod(t, r, "PATCH")
		testutil.AssertPath(t, r, "/v1/calendars/cal_1/events/evt_1")
		testutil.RespondJSON(w, 200, map[string]interface{}{
			"id": "evt_1", "calendar_id": "cal_1", "title": "Renamed",
			"start_time": "2026-04-14T09:00:00Z", "end_time": "2026-04-14T09:30:00Z",
			"all_day": false, "status": "confirmed", "source": "internal",
			"created_at": "2026-04-14T00:00:00Z", "updated_at": "2026-04-14T00:00:00Z",
		})
	}))

	event, err := client.Events.Update(context.Background(), "cal_1", "evt_1", &chronary.UpdateEventParams{
		Title: chronary.String("Renamed"),
	})
	if err != nil {
		t.Fatal(err)
	}
	if event.Title != "Renamed" {
		t.Errorf("expected Renamed, got %s", event.Title)
	}
}

func TestEventDelete_UsesCalendarScopedPath(t *testing.T) {
	client, _ := testutil.Setup(t, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		testutil.AssertMethod(t, r, "DELETE")
		testutil.AssertPath(t, r, "/v1/calendars/cal_1/events/evt_1")
		w.WriteHeader(204)
	}))

	if err := client.Events.Delete(context.Background(), "cal_1", "evt_1"); err != nil {
		t.Fatal(err)
	}
}

func TestEventGetByID_UsesEventOnlyPath(t *testing.T) {
	client, _ := testutil.Setup(t, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		testutil.AssertMethod(t, r, "GET")
		testutil.AssertPath(t, r, "/v1/events/evt_1")
		testutil.RespondJSON(w, 200, map[string]interface{}{
			"id": "evt_1", "calendar_id": "cal_1", "title": "Standup",
			"start_time": "2026-04-14T09:00:00Z", "end_time": "2026-04-14T09:30:00Z",
			"all_day": false, "status": "confirmed", "source": "internal",
			"created_at": "2026-04-14T00:00:00Z", "updated_at": "2026-04-14T00:00:00Z",
		})
	}))

	event, err := client.Events.GetByID(context.Background(), "evt_1")
	if err != nil {
		t.Fatal(err)
	}
	if event.ID != "evt_1" {
		t.Errorf("expected evt_1, got %s", event.ID)
	}
}

func TestEventUpdateByID_UsesEventOnlyPath(t *testing.T) {
	client, _ := testutil.Setup(t, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		testutil.AssertMethod(t, r, "PATCH")
		testutil.AssertPath(t, r, "/v1/events/evt_1")
		testutil.RespondJSON(w, 200, map[string]interface{}{
			"id": "evt_1", "calendar_id": "cal_1", "title": "Renamed",
			"start_time": "2026-04-14T09:00:00Z", "end_time": "2026-04-14T09:30:00Z",
			"all_day": false, "status": "confirmed", "source": "internal",
			"created_at": "2026-04-14T00:00:00Z", "updated_at": "2026-04-14T00:00:00Z",
		})
	}))

	event, err := client.Events.UpdateByID(context.Background(), "evt_1", &chronary.UpdateEventParams{
		Title: chronary.String("Renamed"),
	})
	if err != nil {
		t.Fatal(err)
	}
	if event.Title != "Renamed" {
		t.Errorf("expected Renamed, got %s", event.Title)
	}
}

func TestEventDeleteByID_UsesEventOnlyPath(t *testing.T) {
	client, _ := testutil.Setup(t, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		testutil.AssertMethod(t, r, "DELETE")
		testutil.AssertPath(t, r, "/v1/events/evt_1")
		w.WriteHeader(204)
	}))

	if err := client.Events.DeleteByID(context.Background(), "evt_1"); err != nil {
		t.Fatal(err)
	}
}

func TestEventConfirm(t *testing.T) {
	client, _ := testutil.Setup(t, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		testutil.AssertMethod(t, r, "PUT")
		testutil.AssertPath(t, r, "/v1/events/evt_1/confirm")
		testutil.RespondJSON(w, 200, map[string]interface{}{
			"id": "evt_1", "calendar_id": "cal_1", "title": "Held",
			"start_time": "2026-04-14T09:00:00Z", "end_time": "2026-04-14T09:30:00Z",
			"all_day": false, "status": "confirmed", "source": "internal",
			"created_at": "2026-04-14T00:00:00Z", "updated_at": "2026-04-14T00:00:00Z",
		})
	}))

	event, err := client.Events.Confirm(context.Background(), "evt_1")
	if err != nil {
		t.Fatal(err)
	}
	if event.Status != chronary.EventStatusConfirmed {
		t.Errorf("expected confirmed, got %s", event.Status)
	}
}

func TestEventRelease(t *testing.T) {
	client, _ := testutil.Setup(t, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		testutil.AssertMethod(t, r, "PUT")
		testutil.AssertPath(t, r, "/v1/events/evt_1/release")
		testutil.RespondJSON(w, 200, map[string]interface{}{
			"id": "evt_1", "calendar_id": "cal_1", "title": "Held",
			"start_time": "2026-04-14T09:00:00Z", "end_time": "2026-04-14T09:30:00Z",
			"all_day": false, "status": "cancelled", "source": "internal",
			"created_at": "2026-04-14T00:00:00Z", "updated_at": "2026-04-14T00:00:00Z",
		})
	}))

	if _, err := client.Events.Release(context.Background(), "evt_1"); err != nil {
		t.Fatal(err)
	}
}

// --- Calendar context + availability rules ---

func TestCalendarGetContext(t *testing.T) {
	client, _ := testutil.Setup(t, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		testutil.AssertMethod(t, r, "GET")
		testutil.AssertPath(t, r, "/v1/calendars/cal_1/context")
		testutil.RespondJSON(w, 200, map[string]interface{}{
			"calendar_id":   "cal_1",
			"now":           "2026-04-14T09:15:00Z",
			"agent_status":  "working",
			"current_event": nil,
			"next_event":    nil,
			"recent_events": []interface{}{},
			"upcoming":      []interface{}{},
		})
	}))

	ctxResp, err := client.Calendars.GetContext(context.Background(), "cal_1")
	if err != nil {
		t.Fatal(err)
	}
	if ctxResp.AgentStatus != chronary.AgentLiveStatusWorking {
		t.Errorf("expected working, got %s", ctxResp.AgentStatus)
	}
}

func TestCalendarSetAvailabilityRules(t *testing.T) {
	client, _ := testutil.Setup(t, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		testutil.AssertMethod(t, r, "PUT")
		testutil.AssertPath(t, r, "/v1/calendars/cal_1/availability-rules")
		testutil.RespondJSON(w, 200, map[string]interface{}{
			"id": "avr_1", "calendar_id": "cal_1",
			"buffer_before_minutes": 5, "buffer_after_minutes": 10,
			"working_hours": nil, "timezone": "UTC",
			"created_at": "2026-04-14T00:00:00Z", "updated_at": "2026-04-14T00:00:00Z",
		})
	}))

	five := 5
	ten := 10
	rules, err := client.Calendars.SetAvailabilityRules(context.Background(), "cal_1", &chronary.SetAvailabilityRulesParams{
		BufferBeforeMinutes: &five,
		BufferAfterMinutes:  &ten,
	})
	if err != nil {
		t.Fatal(err)
	}
	if rules.BufferBeforeMinutes != 5 {
		t.Errorf("expected 5, got %d", rules.BufferBeforeMinutes)
	}
}

func TestCalendarGetAvailabilityRules(t *testing.T) {
	client, _ := testutil.Setup(t, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		testutil.AssertMethod(t, r, "GET")
		testutil.AssertPath(t, r, "/v1/calendars/cal_1/availability-rules")
		testutil.RespondJSON(w, 200, map[string]interface{}{
			"id": "avr_1", "calendar_id": "cal_1",
			"buffer_before_minutes": 0, "buffer_after_minutes": 0,
			"working_hours": nil, "timezone": "UTC",
			"created_at": "2026-04-14T00:00:00Z", "updated_at": "2026-04-14T00:00:00Z",
		})
	}))

	if _, err := client.Calendars.GetAvailabilityRules(context.Background(), "cal_1"); err != nil {
		t.Fatal(err)
	}
}

func TestCalendarDeleteAvailabilityRules(t *testing.T) {
	client, _ := testutil.Setup(t, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		testutil.AssertMethod(t, r, "DELETE")
		testutil.AssertPath(t, r, "/v1/calendars/cal_1/availability-rules")
		w.WriteHeader(204)
	}))

	if err := client.Calendars.DeleteAvailabilityRules(context.Background(), "cal_1"); err != nil {
		t.Fatal(err)
	}
}

// --- Webhook deliveries ---

func TestWebhookListDeliveries(t *testing.T) {
	client, _ := testutil.Setup(t, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		testutil.AssertMethod(t, r, "GET")
		testutil.AssertPath(t, r, "/v1/webhooks/whk_1/deliveries")
		testutil.RespondJSON(w, 200, map[string]interface{}{
			"data": []map[string]interface{}{
				{
					"id": "whd_1", "subscription_id": "whk_1",
					"event_type": "event.created", "status": "delivered",
					"attempts": 1, "last_attempt_at": nil, "next_retry_at": nil,
					"created_at": "2026-04-14T00:00:00Z",
				},
			},
			"total": 1, "limit": 20, "offset": 0,
			"stats": map[string]int{"pending": 0, "delivered": 1, "failed": 0},
		})
	}))

	resp, err := client.Webhooks.ListDeliveries(context.Background(), "whk_1", nil)
	if err != nil {
		t.Fatal(err)
	}
	if resp.Total != 1 || resp.Stats.Delivered != 1 {
		t.Errorf("unexpected envelope: %+v", resp)
	}
}

// --- Scheduling proposals ---

func TestProposalCreate(t *testing.T) {
	client, _ := testutil.Setup(t, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		testutil.AssertMethod(t, r, "POST")
		testutil.AssertPath(t, r, "/v1/scheduling/proposals")
		testutil.RespondJSON(w, 201, map[string]interface{}{
			"id": "prp_1", "title": "Sync", "description": nil,
			"organizer_agent_id":     "agt_1",
			"participant_agent_ids":  []string{"agt_2"},
			"calendar_id":            "cal_1",
			"status":                 "pending",
			"expires_at":             nil,
			"resolved_slot":          nil,
			"created_event_id":       nil,
			"metadata":               map[string]interface{}{},
			"created_at":             "2026-04-14T00:00:00Z",
			"updated_at":             "2026-04-14T00:00:00Z",
		})
	}))

	p, err := client.Scheduling.Create(context.Background(), &chronary.CreateProposalParams{
		Title: "Sync", OrganizerAgentID: "agt_1",
		ParticipantAgentIDs: []string{"agt_2"},
		CalendarID:          "cal_1",
		Slots:               []chronary.ProposalSlot{{StartTime: "2026-04-14T10:00:00Z", EndTime: "2026-04-14T10:30:00Z"}},
	})
	if err != nil {
		t.Fatal(err)
	}
	if p.Status != chronary.ProposalStatusPending {
		t.Errorf("expected pending, got %s", p.Status)
	}
}

func TestProposalRespond(t *testing.T) {
	client, _ := testutil.Setup(t, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		testutil.AssertMethod(t, r, "POST")
		testutil.AssertPath(t, r, "/v1/scheduling/proposals/prp_1/respond")
		testutil.RespondJSON(w, 200, map[string]interface{}{
			"id": "prr_1", "agent_id": "agt_2", "response": "accept",
			"selected_slot_id": "pst_1", "counter_slots": nil,
			"message": nil, "created_at": "2026-04-14T00:00:00Z",
		})
	}))

	if _, err := client.Scheduling.Respond(context.Background(), "prp_1", &chronary.RespondToProposalParams{
		AgentID:  "agt_2",
		Response: chronary.ProposalResponseAccept,
	}); err != nil {
		t.Fatal(err)
	}
}

func TestProposalResolve(t *testing.T) {
	client, _ := testutil.Setup(t, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		testutil.AssertMethod(t, r, "POST")
		testutil.AssertPath(t, r, "/v1/scheduling/proposals/prp_1/resolve")
		testutil.RespondJSON(w, 200, map[string]interface{}{
			"status":        "confirmed",
			"resolved_slot": map[string]string{"start_time": "2026-04-14T10:00:00Z", "end_time": "2026-04-14T10:30:00Z"},
		})
	}))

	if _, err := client.Scheduling.Resolve(context.Background(), "prp_1"); err != nil {
		t.Fatal(err)
	}
}

func TestProposalCancel(t *testing.T) {
	client, _ := testutil.Setup(t, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		testutil.AssertMethod(t, r, "POST")
		testutil.AssertPath(t, r, "/v1/scheduling/proposals/prp_1/cancel")
		testutil.RespondJSON(w, 200, map[string]string{"status": "cancelled"})
	}))

	if _, err := client.Scheduling.Cancel(context.Background(), "prp_1"); err != nil {
		t.Fatal(err)
	}
}

// --- Scoped API keys ---

func TestKeysCreate(t *testing.T) {
	client, _ := testutil.Setup(t, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		testutil.AssertMethod(t, r, "POST")
		testutil.AssertPath(t, r, "/v1/keys")
		testutil.RespondJSON(w, 201, map[string]interface{}{
			"id": "skey_1", "key_prefix": "chr_ak_",
			"agent_id": "agt_1", "label": nil,
			"created_at": "2026-04-14T00:00:00Z",
			"key":        "chr_ak_xxxxxxxx",
		})
	}))

	key, err := client.Keys.Create(context.Background(), &chronary.CreateScopedAPIKeyParams{
		AgentID: "agt_1",
	})
	if err != nil {
		t.Fatal(err)
	}
	if key.Key == "" {
		t.Errorf("expected key string in response")
	}
}

func TestKeysList(t *testing.T) {
	client, _ := testutil.Setup(t, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		testutil.AssertMethod(t, r, "GET")
		testutil.AssertPath(t, r, "/v1/keys")
		testutil.RespondJSON(w, 200, map[string]interface{}{
			"keys": []map[string]interface{}{
				{"id": "skey_1", "key_prefix": "chr_ak_",
					"agent_id": "agt_1", "label": nil,
					"created_at": "2026-04-14T00:00:00Z"},
			},
		})
	}))

	keys, err := client.Keys.List(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	if len(keys) != 1 {
		t.Errorf("expected 1 key, got %d", len(keys))
	}
}

func TestKeysDelete(t *testing.T) {
	client, _ := testutil.Setup(t, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		testutil.AssertMethod(t, r, "DELETE")
		testutil.AssertPath(t, r, "/v1/keys/skey_1")
		w.WriteHeader(204)
	}))

	if err := client.Keys.Delete(context.Background(), "skey_1"); err != nil {
		t.Fatal(err)
	}
}

// --- Agent auth ---

func TestAgentAuthSignUp_Anonymous(t *testing.T) {
	client, srv := testutil.Setup(t, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		testutil.AssertMethod(t, r, "POST")
		testutil.AssertPath(t, r, "/v1/agent/sign-up")
		testutil.RespondJSON(w, 200, map[string]interface{}{
			"org_id": "org_1", "agent_id": "agt_1",
			"api_key": "chr_sk_xxx",
			"message": "Verification code sent to email",
		})
	}))
	_ = client // keep authenticated client for compatibility
	_ = srv

	// Reconstruct as anonymous client to verify the WithAnonymous() path.
	anon, err := chronary.NewClient(chronary.WithAnonymous(), chronary.WithBaseURL(srv.URL), chronary.WithMaxRetries(0))
	if err != nil {
		t.Fatal(err)
	}

	resp, err := anon.AgentAuth.SignUp(context.Background(), &chronary.AgentSignUpParams{
		Email: "agent@example.com", AgentName: "Bot", TosVersion: "2026-04-01",
	})
	if err != nil {
		t.Fatal(err)
	}
	if !resp.IsNewOrg() {
		t.Errorf("expected new-org branch, got %+v", resp)
	}
}

func TestAgentAuthVerify(t *testing.T) {
	client, _ := testutil.Setup(t, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		testutil.AssertMethod(t, r, "POST")
		testutil.AssertPath(t, r, "/v1/agent/verify")
		testutil.RespondJSON(w, 200, map[string]interface{}{"verified": true, "message": "Full access unlocked"})
	}))

	resp, err := client.AgentAuth.Verify(context.Background(), &chronary.AgentVerifyParams{OTP: "123456"})
	if err != nil {
		t.Fatal(err)
	}
	if !resp.Verified {
		t.Errorf("expected verified=true")
	}
}

// --- Feedback + Plans ---

func TestFeedbackSubmit(t *testing.T) {
	client, _ := testutil.Setup(t, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		testutil.AssertMethod(t, r, "POST")
		testutil.AssertPath(t, r, "/v1/feedback")
		testutil.RespondJSON(w, 202, map[string]string{"status": "accepted"})
	}))

	if _, err := client.Feedback.Submit(context.Background(), &chronary.SubmitFeedbackParams{
		Type: chronary.FeedbackTypeBug, Message: "Reproducible 500 on PATCH /v1/calendars/{id}",
	}); err != nil {
		t.Fatal(err)
	}
}

func TestPlansList(t *testing.T) {
	client, _ := testutil.Setup(t, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		testutil.AssertMethod(t, r, "GET")
		testutil.AssertPath(t, r, "/v1/plans")
		testutil.RespondJSON(w, 200, map[string]interface{}{
			"plans": []map[string]interface{}{
				{"id": "free", "name": "Free", "tagline": "Get started",
					"price": 0, "currency": "usd", "limits": nil,
					"display_features": []string{}, "recommended": false},
			},
		})
	}))

	resp, err := client.Plans.List(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	if len(resp.Plans) != 1 {
		t.Errorf("expected 1 plan, got %d", len(resp.Plans))
	}
}

// --- Anonymous client semantics ---

func TestAnonymousClient_HasNoAuthHeader(t *testing.T) {
	gotAuth := ""
	srv := testutil.SetupRaw(t, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotAuth = r.Header.Get("Authorization")
		testutil.RespondJSON(w, 200, map[string]interface{}{"plans": []interface{}{}})
	}))

	anon, err := chronary.NewClient(chronary.WithAnonymous(), chronary.WithBaseURL(srv.URL), chronary.WithMaxRetries(0))
	if err != nil {
		t.Fatal(err)
	}
	if _, err := anon.Plans.List(context.Background()); err != nil {
		t.Fatal(err)
	}
	if gotAuth != "" {
		t.Errorf("expected no Authorization header on anonymous client, got %q", gotAuth)
	}
}
