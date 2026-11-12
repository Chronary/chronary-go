// Surface tests assert that every public client method hits the canonical API
// path. They guard against future regressions like the GET /v1/events/{id}
// bug that shipped in earlier versions before the route was mounted under
// /v1/calendars/{cal_id}/events/{id}.

package chronary_test

import (
	"context"
	"encoding/json"
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
			"id": "evt_1", "calendarId": "cal_1", "orgId": "org_1", "title": "Standup",
			"startTime": "2026-04-14T09:00:00Z", "endTime": "2026-04-14T09:30:00Z",
			"allDay": false, "status": "confirmed", "source": "internal",
			"createdAt": "2026-04-14T00:00:00Z", "updatedAt": "2026-04-14T00:00:00Z",
		})
	}))

	event, err := client.Events.Get(context.Background(), "cal_1", "evt_1")
	if err != nil {
		t.Fatal(err)
	}
	if event.ID != "evt_1" {
		t.Errorf("expected evt_1, got %s", event.ID)
	}
	// Guard the live (camelCase) response shape: these deserialized as empty /
	// zero values when the struct tags were snake_case.
	if event.CalendarID != "cal_1" {
		t.Errorf("expected calendarId cal_1, got %q", event.CalendarID)
	}
	if event.OrgID != "org_1" {
		t.Errorf("expected orgId org_1, got %q", event.OrgID)
	}
	if event.StartTime.IsZero() || event.EndTime.IsZero() {
		t.Error("expected StartTime/EndTime to deserialize from camelCase, got zero time")
	}
	if event.CreatedAt.IsZero() || event.UpdatedAt.IsZero() {
		t.Error("expected CreatedAt/UpdatedAt to deserialize from camelCase, got zero time")
	}
}

func TestEventUpdate_UsesCalendarScopedPath(t *testing.T) {
	client, _ := testutil.Setup(t, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		testutil.AssertMethod(t, r, "PATCH")
		testutil.AssertPath(t, r, "/v1/calendars/cal_1/events/evt_1")
		testutil.RespondJSON(w, 200, map[string]interface{}{
			"id": "evt_1", "calendarId": "cal_1", "title": "Renamed",
			"startTime": "2026-04-14T09:00:00Z", "endTime": "2026-04-14T09:30:00Z",
			"allDay": false, "status": "confirmed", "source": "internal",
			"createdAt": "2026-04-14T00:00:00Z", "updatedAt": "2026-04-14T00:00:00Z",
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
			"id": "evt_1", "calendarId": "cal_1", "orgId": "org_1", "title": "Standup",
			"startTime": "2026-04-14T09:00:00Z", "endTime": "2026-04-14T09:30:00Z",
			"allDay": false, "status": "confirmed", "source": "internal",
			"createdAt": "2026-04-14T00:00:00Z", "updatedAt": "2026-04-14T00:00:00Z",
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
			"id": "evt_1", "calendarId": "cal_1", "title": "Renamed",
			"startTime": "2026-04-14T09:00:00Z", "endTime": "2026-04-14T09:30:00Z",
			"allDay": false, "status": "confirmed", "source": "internal",
			"createdAt": "2026-04-14T00:00:00Z", "updatedAt": "2026-04-14T00:00:00Z",
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
			"id": "evt_1", "calendarId": "cal_1", "title": "Held",
			"startTime": "2026-04-14T09:00:00Z", "endTime": "2026-04-14T09:30:00Z",
			"allDay": false, "status": "confirmed", "source": "internal",
			"createdAt": "2026-04-14T00:00:00Z", "updatedAt": "2026-04-14T00:00:00Z",
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
			"id": "evt_1", "calendarId": "cal_1", "title": "Held",
			"startTime": "2026-04-14T09:00:00Z", "endTime": "2026-04-14T09:30:00Z",
			"allDay": false, "status": "cancelled", "source": "internal",
			"createdAt": "2026-04-14T00:00:00Z", "updatedAt": "2026-04-14T00:00:00Z",
		})
	}))

	if _, err := client.Events.Release(context.Background(), "evt_1"); err != nil {
		t.Fatal(err)
	}
}

// --- Recurring events (#996) ---

func TestEventCreate_SendsRecurrenceRule(t *testing.T) {
	client, _ := testutil.Setup(t, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		testutil.AssertMethod(t, r, "POST")
		testutil.AssertPath(t, r, "/v1/calendars/cal_1/events")
		var body map[string]interface{}
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			t.Fatal(err)
		}
		if body["recurrence_rule"] != "FREQ=WEEKLY;BYDAY=MO,WE;COUNT=12" {
			t.Errorf("expected recurrence_rule in request body, got %v", body["recurrence_rule"])
		}
		testutil.RespondJSON(w, 201, map[string]interface{}{
			"id": "evt_1", "calendarId": "cal_1", "orgId": "org_1", "title": "Standup",
			"startTime": "2026-04-13T09:00:00Z", "endTime": "2026-04-13T09:30:00Z",
			"allDay": false, "status": "confirmed", "source": "internal",
			"recurrenceRule": "FREQ=WEEKLY;BYDAY=MO,WE;COUNT=12", "recurrenceExdates": []string{},
			"createdAt": "2026-04-13T00:00:00Z", "updatedAt": "2026-04-13T00:00:00Z",
		})
	}))

	event, err := client.Events.Create(context.Background(), "cal_1", &chronary.CreateEventParams{
		Title:          "Standup",
		StartTime:      "2026-04-13T09:00:00Z",
		EndTime:        "2026-04-13T09:30:00Z",
		RecurrenceRule: chronary.String("FREQ=WEEKLY;BYDAY=MO,WE;COUNT=12"),
	})
	if err != nil {
		t.Fatal(err)
	}
	if event.RecurrenceRule == nil || *event.RecurrenceRule != "FREQ=WEEKLY;BYDAY=MO,WE;COUNT=12" {
		t.Errorf("expected recurrenceRule to deserialize, got %v", event.RecurrenceRule)
	}
	if event.RecurrenceExdates == nil || len(event.RecurrenceExdates) != 0 {
		t.Errorf("expected empty recurrenceExdates, got %v", event.RecurrenceExdates)
	}
}

func TestEventUpdate_SendsRecurrenceRule(t *testing.T) {
	client, _ := testutil.Setup(t, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		testutil.AssertMethod(t, r, "PATCH")
		testutil.AssertPath(t, r, "/v1/calendars/cal_1/events/evt_1")
		var body map[string]interface{}
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			t.Fatal(err)
		}
		if body["recurrence_rule"] != "FREQ=DAILY;COUNT=5" {
			t.Errorf("expected recurrence_rule in request body, got %v", body["recurrence_rule"])
		}
		testutil.RespondJSON(w, 200, map[string]interface{}{
			"id": "evt_1", "calendarId": "cal_1", "title": "Standup",
			"startTime": "2026-04-13T09:00:00Z", "endTime": "2026-04-13T09:30:00Z",
			"allDay": false, "status": "confirmed", "source": "internal",
			"recurrenceRule": "FREQ=DAILY;COUNT=5", "recurrenceExdates": []string{},
			"createdAt": "2026-04-13T00:00:00Z", "updatedAt": "2026-04-13T00:00:00Z",
		})
	}))

	event, err := client.Events.Update(context.Background(), "cal_1", "evt_1", &chronary.UpdateEventParams{
		RecurrenceRule: chronary.String("FREQ=DAILY;COUNT=5"),
	})
	if err != nil {
		t.Fatal(err)
	}
	if event.RecurrenceRule == nil || *event.RecurrenceRule != "FREQ=DAILY;COUNT=5" {
		t.Errorf("expected recurrenceRule to deserialize, got %v", event.RecurrenceRule)
	}
}

func TestEventList_ExpandQuery(t *testing.T) {
	client, _ := testutil.Setup(t, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		testutil.AssertMethod(t, r, "GET")
		testutil.AssertPath(t, r, "/v1/calendars/cal_1/events")
		q := r.URL.Query()
		if q.Get("expand") != "true" {
			t.Errorf("expected expand=true, got %q", q.Get("expand"))
		}
		if q.Get("start_after") == "" || q.Get("start_before") == "" {
			t.Error("expected start_after and start_before to be sent with expand")
		}
		testutil.RespondJSON(w, 200, map[string]interface{}{
			"data": []map[string]interface{}{
				{
					"id": "evt_1", "calendarId": "cal_1", "title": "Standup",
					"startTime": "2026-04-15T09:00:00Z", "endTime": "2026-04-15T09:30:00Z",
					"allDay": false, "status": "confirmed", "source": "internal",
					"recurringEventId":  "evt_1",
					"originalStartTime": "2026-04-15T09:00:00Z",
					"createdAt":         "2026-04-13T00:00:00Z", "updatedAt": "2026-04-13T00:00:00Z",
				},
			},
			"total": 1, "limit": 50, "offset": 0,
		})
	}))

	events, err := client.Events.List(context.Background(), &chronary.ListEventsParams{
		CalendarID:  chronary.String("cal_1"),
		StartAfter:  chronary.String("2026-04-13T00:00:00Z"),
		StartBefore: chronary.String("2026-05-13T00:00:00Z"),
		Expand:      chronary.Bool(true),
	}).Collect(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	if len(events) != 1 {
		t.Fatalf("expected 1 event, got %d", len(events))
	}
	// Guard the expanded-instance fields.
	if events[0].RecurringEventID != "evt_1" {
		t.Errorf("expected recurringEventId evt_1, got %q", events[0].RecurringEventID)
	}
	if events[0].OriginalStartTime != "2026-04-15T09:00:00Z" {
		t.Errorf("expected originalStartTime to deserialize, got %q", events[0].OriginalStartTime)
	}
}

func TestEventCancelOccurrence_UsesCalendarScopedPath(t *testing.T) {
	client, _ := testutil.Setup(t, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		testutil.AssertMethod(t, r, "DELETE")
		testutil.AssertPath(t, r, "/v1/calendars/cal_1/events/evt_1")
		if got := r.URL.Query().Get("occurrence_start"); got != "2026-04-15T09:00:00Z" {
			t.Errorf("expected occurrence_start query param, got %q", got)
		}
		// occurrence-cancel returns 200 with the updated series master (not 204).
		testutil.RespondJSON(w, 200, map[string]interface{}{
			"id": "evt_1", "calendarId": "cal_1", "title": "Standup",
			"startTime": "2026-04-13T09:00:00Z", "endTime": "2026-04-13T09:30:00Z",
			"allDay": false, "status": "confirmed", "source": "internal",
			"recurrenceRule":    "FREQ=WEEKLY;BYDAY=MO,WE;COUNT=12",
			"recurrenceExdates": []string{"2026-04-15T09:00:00Z"},
			"createdAt":         "2026-04-13T00:00:00Z", "updatedAt": "2026-04-13T00:00:00Z",
		})
	}))

	event, err := client.Events.CancelOccurrence(context.Background(), "cal_1", "evt_1", "2026-04-15T09:00:00Z")
	if err != nil {
		t.Fatal(err)
	}
	if len(event.RecurrenceExdates) != 1 || event.RecurrenceExdates[0] != "2026-04-15T09:00:00Z" {
		t.Errorf("expected cancelled occurrence in recurrenceExdates, got %v", event.RecurrenceExdates)
	}
}

func TestEventCancelOccurrenceByID_UsesEventOnlyPath(t *testing.T) {
	client, _ := testutil.Setup(t, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		testutil.AssertMethod(t, r, "DELETE")
		testutil.AssertPath(t, r, "/v1/events/evt_1")
		if got := r.URL.Query().Get("occurrence_start"); got != "2026-04-15T09:00:00Z" {
			t.Errorf("expected occurrence_start query param, got %q", got)
		}
		testutil.RespondJSON(w, 200, map[string]interface{}{
			"id": "evt_1", "calendarId": "cal_1", "title": "Standup",
			"startTime": "2026-04-13T09:00:00Z", "endTime": "2026-04-13T09:30:00Z",
			"allDay": false, "status": "confirmed", "source": "internal",
			"recurrenceRule":    "FREQ=WEEKLY;BYDAY=MO,WE;COUNT=12",
			"recurrenceExdates": []string{"2026-04-15T09:00:00Z"},
			"createdAt":         "2026-04-13T00:00:00Z", "updatedAt": "2026-04-13T00:00:00Z",
		})
	}))

	event, err := client.Events.CancelOccurrenceByID(context.Background(), "evt_1", "2026-04-15T09:00:00Z")
	if err != nil {
		t.Fatal(err)
	}
	if len(event.RecurrenceExdates) != 1 {
		t.Errorf("expected cancelled occurrence in recurrenceExdates, got %v", event.RecurrenceExdates)
	}
}

func TestUsageGet_IncludesRecurringEvents(t *testing.T) {
	client, _ := testutil.Setup(t, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		testutil.AssertMethod(t, r, "GET")
		testutil.AssertPath(t, r, "/v1/usage")
		testutil.RespondJSON(w, 200, map[string]interface{}{
			"period_start": "2026-07-01T00:00:00Z", "period_end": "2026-08-01T00:00:00Z",
			"plan":             "free",
			"agents":           map[string]interface{}{"used": 1, "limit": 3},
			"calendars":        map[string]interface{}{"used": 1, "limit": 10},
			"events":           map[string]interface{}{"used": 5, "limit": 2500},
			"api_calls":        map[string]interface{}{"used": 100, "limit": 50000},
			"webhooks":         map[string]interface{}{"used": 0, "limit": 5000},
			"recurring_events": map[string]interface{}{"used": 2, "limit": 10},
		})
	}))

	usage, err := client.Usage.Get(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	if usage.RecurringEvents.Used != 2 {
		t.Errorf("expected recurring_events.used 2, got %d", usage.RecurringEvents.Used)
	}
	if usage.RecurringEvents.Limit == nil || *usage.RecurringEvents.Limit != 10 {
		t.Errorf("expected recurring_events.limit 10, got %v", usage.RecurringEvents.Limit)
	}
}

func TestUsageHistory(t *testing.T) {
	var gotQuery string
	client, _ := testutil.Setup(t, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		testutil.AssertMethod(t, r, "GET")
		testutil.AssertPath(t, r, "/v1/usage/history")
		gotQuery = r.URL.RawQuery
		testutil.RespondJSON(w, 200, map[string]interface{}{
			"plan":   "free",
			"months": 1,
			"limits": map[string]interface{}{
				"api_calls": 50000, "events": 2500, "webhook_deliveries": 5000,
				"availability_queries": 10000, "proposals": nil,
			},
			"data": []map[string]interface{}{{
				"period": "2026-10", "period_start": "2026-10-01T00:00:00.000Z",
				"period_end": "2026-10-31T23:59:59.000Z", "is_current": true,
				"metrics": map[string]interface{}{"api_calls": 12480, "events": 3},
			}},
		})
	}))

	history, err := client.Usage.History(context.Background(), 1)
	if err != nil {
		t.Fatal(err)
	}
	if gotQuery != "months=1" {
		t.Errorf("expected months=1 query, got %q", gotQuery)
	}
	if len(history.Data) != 1 || !history.Data[0].IsCurrent || history.Data[0].Metrics.APICalls != 12480 {
		t.Errorf("unexpected history data: %+v", history.Data)
	}
	if history.Limits.APICalls == nil || *history.Limits.APICalls != 50000 {
		t.Errorf("expected api_calls limit 50000, got %v", history.Limits.APICalls)
	}
	if history.Limits.Proposals != nil {
		t.Errorf("expected nil proposals limit, got %v", *history.Limits.Proposals)
	}

	// months=0 omits the query so the server default applies.
	if _, err := client.Usage.History(context.Background(), 0); err != nil {
		t.Fatal(err)
	}
	if gotQuery != "" {
		t.Errorf("expected no query for months=0, got %q", gotQuery)
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

func TestWebhookListDeliveriesFilters(t *testing.T) {
	client, _ := testutil.Setup(t, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		testutil.AssertPath(t, r, "/v1/webhooks/whk_1/deliveries")
		q := r.URL.Query()
		want := map[string]string{
			"status": "dropped", "event_type": "event.", "created_after": "2026-04-18T00:00:00Z",
			"delivery_id": "whd_1", "is_test": "false", "cursor": "abc", "limit": "5",
		}
		for k, v := range want {
			if got := q.Get(k); got != v {
				t.Errorf("query %s = %q, want %q", k, got, v)
			}
		}
		testutil.RespondJSON(w, 200, map[string]interface{}{
			"data": []map[string]interface{}{}, "total": 3, "limit": 5, "offset": 0,
			"next_cursor": "next", "has_more": true,
			"stats": map[string]int{"pending": 0, "delivered": 1, "failed": 0, "dropped": 2},
		})
	}))

	status := chronary.WebhookDeliveryDropped
	eventType, after, id, cursor := "event.", "2026-04-18T00:00:00Z", "whd_1", "abc"
	isTest := false
	resp, err := client.Webhooks.ListDeliveries(context.Background(), "whk_1", &chronary.ListWebhookDeliveriesParams{
		Status: &status, EventType: &eventType, CreatedAfter: &after, DeliveryID: &id,
		IsTest: &isTest, Cursor: &cursor, Limit: 5,
	})
	if err != nil {
		t.Fatal(err)
	}
	if !resp.HasMore || resp.NextCursor == nil || *resp.NextCursor != "next" || resp.Stats.Dropped != 2 {
		t.Errorf("unexpected envelope: %+v", resp)
	}
}

func TestWebhookGetDelivery(t *testing.T) {
	client, _ := testutil.Setup(t, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		testutil.AssertMethod(t, r, "GET")
		testutil.AssertPath(t, r, "/v1/webhooks/whk_1/deliveries/whd_1")
		testutil.RespondJSON(w, 200, map[string]interface{}{
			"id": "whd_1", "subscription_id": "whk_1",
			"event_type": "event.created", "status": "failed",
			"attempts": 1, "max_attempts": 4,
			"last_attempt_at": "2026-04-14T00:00:00Z", "next_retry_at": nil,
			"created_at":           "2026-04-14T00:00:00Z",
			"last_response_status": 503, "last_error_code": "non_2xx", "last_duration_ms": 120,
			"payload": map[string]interface{}{"id": "evt_1"}, "payload_redacted": false,
			"request_headers": map[string]string{
				"content-type": "application/json", "x-delivery-id": "whd_1",
				"x-chronary-event-type": "event.created",
			},
			"attempt_log": []map[string]interface{}{
				{
					"id": "wha_1", "attempt_number": 1, "started_at": "2026-04-14T00:00:00Z",
					"duration_ms": 120, "response_status": 503, "response_body": "unavailable",
					"error_code": "non_2xx",
				},
			},
		})
	}))

	resp, err := client.Webhooks.GetDelivery(context.Background(), "whk_1", "whd_1")
	if err != nil {
		t.Fatal(err)
	}
	if resp.MaxAttempts != 4 || len(resp.AttemptLog) != 1 {
		t.Fatalf("unexpected detail: %+v", resp)
	}
	a := resp.AttemptLog[0]
	if a.ErrorCode == nil || *a.ErrorCode != chronary.WebhookAttemptErrorNon2xx {
		t.Errorf("expected error_code non_2xx, got %v", a.ErrorCode)
	}
	if resp.LastResponseStatus == nil || *resp.LastResponseStatus != 503 {
		t.Errorf("expected last_response_status 503, got %v", resp.LastResponseStatus)
	}
}

func TestWebhookSendTest(t *testing.T) {
	client, _ := testutil.Setup(t, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		testutil.AssertMethod(t, r, "POST")
		testutil.AssertPath(t, r, "/v1/webhooks/whk_1/test")
		var body map[string]interface{}
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			t.Fatal(err)
		}
		if body["event_type"] != "event.created" {
			t.Errorf("expected event_type event.created, got %v", body["event_type"])
		}
		testutil.RespondJSON(w, 202, map[string]interface{}{
			"delivery_id": "whd_t", "webhook_id": "whk_1",
			"event_type": "event.created", "status": "pending", "is_test": true,
		})
	}))

	resp, err := client.Webhooks.SendTest(context.Background(), "whk_1", &chronary.SendTestWebhookParams{EventType: "event.created"})
	if err != nil {
		t.Fatal(err)
	}
	if resp.DeliveryID != "whd_t" || !resp.IsTest || resp.Status != "pending" {
		t.Fatalf("unexpected response: %+v", resp)
	}
}

func TestWebhookReplayDelivery(t *testing.T) {
	client, _ := testutil.Setup(t, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		testutil.AssertMethod(t, r, "POST")
		testutil.AssertPath(t, r, "/v1/webhooks/whk_1/deliveries/whd_orig/replay")
		if got := r.Header.Get("Idempotency-Key"); got != "replay-1" {
			t.Errorf("expected Idempotency-Key replay-1, got %q", got)
		}
		testutil.RespondJSON(w, 202, map[string]interface{}{
			"delivery_id": "whd_new", "webhook_id": "whk_1", "event_type": "event.created",
			"status": "pending", "replay_of_delivery_id": "whd_orig",
		})
	}))

	resp, err := client.Webhooks.ReplayDelivery(context.Background(), "whk_1", "whd_orig", chronary.WithIdempotencyKey("replay-1"))
	if err != nil {
		t.Fatal(err)
	}
	if resp.DeliveryID != "whd_new" || resp.ReplayOfDeliveryID != "whd_orig" || resp.Status != "pending" {
		t.Fatalf("unexpected response: %+v", resp)
	}
}

func TestWebhookDeliveryReplayLink(t *testing.T) {
	client, _ := testutil.Setup(t, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		testutil.RespondJSON(w, 200, map[string]interface{}{
			"id": "whd_new", "subscription_id": "whk_1", "event_type": "event.created", "status": "pending",
			"attempts": 0, "max_attempts": 8, "created_at": "2026-11-04T00:00:00Z", "is_test": false,
			"replay_of_delivery_id": "whd_orig", "payload": map[string]interface{}{}, "payload_redacted": false,
			"request_headers": map[string]string{"x-chronary-replay-of": "whd_orig"}, "attempt_log": []interface{}{},
		})
	}))

	d, err := client.Webhooks.GetDelivery(context.Background(), "whk_1", "whd_new")
	if err != nil {
		t.Fatal(err)
	}
	if d.ReplayOfDeliveryID == nil || *d.ReplayOfDeliveryID != "whd_orig" || d.RequestHeaders["x-chronary-replay-of"] != "whd_orig" {
		t.Fatalf("unexpected delivery: %+v", d)
	}
}

func TestWebhookRotateSecret(t *testing.T) {
	client, _ := testutil.Setup(t, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		testutil.AssertMethod(t, r, "POST")
		testutil.AssertPath(t, r, "/v1/webhooks/whk_1/secret/rotate")
		var body map[string]interface{}
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			t.Fatal(err)
		}
		if body["expire_old_in"] != "1h" {
			t.Errorf("expected expire_old_in 1h, got %v", body["expire_old_in"])
		}
		testutil.RespondJSON(w, 200, map[string]interface{}{
			"webhook_id": "whk_1", "secret": "whsec_new",
			"previous_secret_expires_at": "2026-11-05T10:00:00.000Z",
		})
	}))

	resp, err := client.Webhooks.RotateSecret(context.Background(), "whk_1", &chronary.RotateWebhookSecretParams{ExpireOldIn: chronary.WebhookSecretExpire1h})
	if err != nil {
		t.Fatal(err)
	}
	if resp.Secret != "whsec_new" || resp.PreviousSecretExpiresAt == nil {
		t.Fatalf("unexpected response: %+v", resp)
	}
}

func TestWebhookRotateSecretDefaults(t *testing.T) {
	client, _ := testutil.Setup(t, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var body map[string]interface{}
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			t.Fatal(err)
		}
		if _, ok := body["expire_old_in"]; ok {
			t.Errorf("expected no expire_old_in, got %v", body)
		}
		testutil.RespondJSON(w, 200, map[string]interface{}{
			"webhook_id": "whk_1", "secret": "whsec_new", "previous_secret_expires_at": nil,
		})
	}))

	resp, err := client.Webhooks.RotateSecret(context.Background(), "whk_1", nil)
	if err != nil {
		t.Fatal(err)
	}
	if resp.PreviousSecretExpiresAt != nil {
		t.Fatalf("expected nil expiry, got %v", *resp.PreviousSecretExpiresAt)
	}
}

// --- Scheduling proposals ---

func TestProposalCreate(t *testing.T) {
	client, _ := testutil.Setup(t, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		testutil.AssertMethod(t, r, "POST")
		testutil.AssertPath(t, r, "/v1/scheduling/proposals")
		testutil.RespondJSON(w, 201, map[string]interface{}{
			"id": "prp_1", "title": "Sync", "description": nil,
			"organizer_agent_id":    "agt_1",
			"participant_agent_ids": []string{"agt_2"},
			"calendar_id":           "cal_1",
			"status":                "pending",
			"expires_at":            nil,
			"resolved_slot":         nil,
			"created_event_id":      nil,
			"metadata":              map[string]interface{}{},
			"created_at":            "2026-04-14T00:00:00Z",
			"updated_at":            "2026-04-14T00:00:00Z",
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
					"price": 0, "currency": "usd",
					"limits": map[string]interface{}{
						"agents": 3, "calendars": 10, "events": 2500,
						"api_calls": 50000, "webhook_deliveries": 5000,
						"availability_queries": 10000, "ical_subscriptions": 5,
						"proposals": 0, "webhook_endpoints": 3, "scoped_keys": 0,
					},
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
	// Guard the plan-limit fields that were previously dropped by the SDK model.
	limits := resp.Plans[0].Limits
	if limits == nil {
		t.Fatal("expected non-nil limits")
	}
	if limits.WebhookEndpoints == nil || *limits.WebhookEndpoints != 3 {
		t.Errorf("expected webhook_endpoints 3, got %v", limits.WebhookEndpoints)
	}
	if limits.ScopedKeys == nil || *limits.ScopedKeys != 0 {
		t.Errorf("expected scoped_keys 0, got %v", limits.ScopedKeys)
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
