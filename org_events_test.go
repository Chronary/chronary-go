package chronary_test

import (
	"context"
	"net/http"
	"reflect"
	"testing"

	chronary "github.com/Chronary/chronary-go"
	"github.com/Chronary/chronary-go/internal/testutil"
)

// #1202: org-wide events list, calendar search and subscribe URL rotation.

func TestEventListAll_RepeatsCalendarID(t *testing.T) {
	client, _ := testutil.Setup(t, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		testutil.AssertMethod(t, r, "GET")
		testutil.AssertPath(t, r, "/v1/events")
		q := r.URL.Query()
		if got := q["calendar_id"]; !reflect.DeepEqual(got, []string{"cal_a", "cal_b"}) {
			t.Errorf("calendar_id = %v", got)
		}
		if q.Get("q") != "standup" || q.Get("agent_id") != "agt_1" {
			t.Errorf("unexpected query %s", r.URL.RawQuery)
		}
		testutil.RespondJSON(w, 200, map[string]interface{}{
			"data": []map[string]interface{}{{
				"id": "evt_1", "calendarId": "cal_a", "orgId": "org_1", "title": "Standup",
				"startTime": "2026-04-14T09:00:00Z", "endTime": "2026-04-14T09:30:00Z",
				"allDay": false, "status": "confirmed", "source": "internal",
				"createdAt": "2026-04-14T00:00:00Z", "updatedAt": "2026-04-14T00:00:00Z",
			}},
			"total": 1, "limit": 50, "offset": 0, "next_cursor": nil,
		})
	}))

	it := client.Events.ListAll(context.Background(), &chronary.ListOrgEventsParams{
		CalendarIDs: []string{"cal_a", "cal_b"},
		AgentID:     chronary.String("agt_1"),
		Q:           chronary.String("standup"),
		StartAfter:  chronary.String("2026-04-01T00:00:00Z"),
		StartBefore: chronary.String("2026-04-15T00:00:00Z"),
	})
	page, err := it.GetPage(context.Background(), 0)
	if err != nil {
		t.Fatal(err)
	}
	if len(page.Data) != 1 || page.Data[0].ID != "evt_1" {
		t.Errorf("unexpected page %+v", page.Data)
	}
}

func TestEventListAll_NilParams(t *testing.T) {
	client, _ := testutil.Setup(t, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		testutil.AssertPath(t, r, "/v1/events")
		testutil.RespondJSON(w, 200, map[string]interface{}{"data": []interface{}{}, "total": 0, "limit": 50, "offset": 0})
	}))

	if _, err := client.Events.ListAll(context.Background(), nil).GetPage(context.Background(), 0); err != nil {
		t.Fatal(err)
	}
}

func TestCalendarRotateICalToken(t *testing.T) {
	client, _ := testutil.Setup(t, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		testutil.AssertMethod(t, r, "POST")
		testutil.AssertPath(t, r, "/v1/calendars/cal_1/ical-token/rotate")
		testutil.RespondJSON(w, 200, map[string]interface{}{
			"id": "cal_1", "orgId": "org_1", "name": "Work", "timezone": "UTC",
			"ical_url":  "https://api.chronary.ai/ical/new.ics",
			"createdAt": "2026-04-14T00:00:00Z", "updatedAt": "2026-04-14T00:00:00Z",
		})
	}))

	cal, err := client.Calendars.RotateICalToken(context.Background(), "cal_1")
	if err != nil {
		t.Fatal(err)
	}
	if cal.ICalURL != "https://api.chronary.ai/ical/new.ics" {
		t.Errorf("ICalURL = %s", cal.ICalURL)
	}
}

func TestCalendarListSearch(t *testing.T) {
	client, _ := testutil.Setup(t, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		testutil.AssertPath(t, r, "/v1/calendars")
		if r.URL.Query().Get("q") != "Work" {
			t.Errorf("q = %q", r.URL.Query().Get("q"))
		}
		testutil.RespondJSON(w, 200, map[string]interface{}{"data": []interface{}{}, "total": 0, "limit": 50, "offset": 0})
	}))

	if _, err := client.Calendars.List(context.Background(), &chronary.ListCalendarsParams{Q: chronary.String("Work")}).GetPage(context.Background(), 0); err != nil {
		t.Fatal(err)
	}
}
