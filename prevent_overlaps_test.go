package chronary

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

// #595: the double-booking guarantee. prevent_overlaps is sent only when set,
// decoded on every calendar, and a refusal exposes Code and ConflictingEventIDs.

func TestCalendarPreventOverlapsRequestBodies(t *testing.T) {
	var bodies []map[string]json.RawMessage
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		b, _ := io.ReadAll(r.Body)
		var m map[string]json.RawMessage
		_ = json.Unmarshal(b, &m)
		bodies = append(bodies, m)
		w.WriteHeader(200)
		w.Write([]byte(`{"id":"cal_1","name":"c","timezone":"UTC","prevent_overlaps":true}`))
	}))
	defer srv.Close()
	client, _ := NewClient(WithAPIKey("test_key"), WithBaseURL(srv.URL))
	ctx := context.Background()

	cal, err := client.Calendars.Create(ctx, &CreateCalendarParams{Name: "c", Timezone: "UTC", PreventOverlaps: Bool(true)})
	if err != nil {
		t.Fatal(err)
	}
	if !cal.PreventOverlaps {
		t.Error("expected decoded PreventOverlaps true")
	}
	if _, err := client.Calendars.Create(ctx, &CreateCalendarParams{Name: "c", Timezone: "UTC"}); err != nil {
		t.Fatal(err)
	}
	if _, err := client.Calendars.Update(ctx, "cal_1", &UpdateCalendarParams{PreventOverlaps: Bool(false)}); err != nil {
		t.Fatal(err)
	}
	if _, err := client.Calendars.Update(ctx, "cal_1", &UpdateCalendarParams{Name: String("n")}); err != nil {
		t.Fatal(err)
	}

	want := []string{"true", "", "false", ""}
	if len(bodies) != len(want) {
		t.Fatalf("expected %d requests, got %d", len(want), len(bodies))
	}
	for i, w := range want {
		raw, ok := bodies[i]["prevent_overlaps"]
		if w == "" {
			if ok {
				t.Errorf("request %d: expected prevent_overlaps omitted, got %s", i, raw)
			}
			continue
		}
		if string(raw) != w {
			t.Errorf("request %d: expected prevent_overlaps=%s, got %q", i, w, raw)
		}
	}
}

func TestCalendarPreventOverlapsDefaultsFalse(t *testing.T) {
	var cal Calendar
	if err := json.Unmarshal([]byte(`{"id":"cal_1"}`), &cal); err != nil {
		t.Fatal(err)
	}
	if cal.PreventOverlaps {
		t.Error("expected PreventOverlaps false when absent")
	}
}

func TestEventConflictErrorExposesCodeAndIDs(t *testing.T) {
	body := []byte(`{"error":{"type":"conflict","code":"event_conflict","message":"This calendar refuses overlapping events","conflicting_event_ids":["evt_1","evt_2"],"request_id":"req_1"}}`)
	apiErr := parseErrorResponse(409, body, http.Header{})
	if apiErr.StatusCode != 409 || apiErr.Type != "conflict" || apiErr.Code != "event_conflict" {
		t.Errorf("unexpected error: %+v", apiErr)
	}
	if len(apiErr.ConflictingEventIDs) != 2 || apiErr.ConflictingEventIDs[0] != "evt_1" || apiErr.ConflictingEventIDs[1] != "evt_2" {
		t.Errorf("unexpected ConflictingEventIDs: %#v", apiErr.ConflictingEventIDs)
	}
	var target *Error
	if !errors.As(error(apiErr), &target) {
		t.Fatal("expected errors.As to match *Error")
	}

	plain := parseErrorResponse(404, []byte(`{"error":{"type":"not_found","message":"gone"}}`), http.Header{})
	if plain.Code != "" || plain.ConflictingEventIDs != nil {
		t.Errorf("expected no code or ids on a plain error, got %+v", plain)
	}
}
