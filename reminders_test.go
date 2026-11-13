package chronary

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strings"
	"testing"
)

// #1991: an empty reminder schedule must reach the API as [] (no reminders),
// not be dropped by omitempty (which made the event/calendar inherit the
// default instead). nil = omit, [] = none, [..] = set, null = inherit.

func TestReminderHelpers(t *testing.T) {
	if got := Reminders(10, 1440); got == nil || !reflect.DeepEqual(*got, []int{10, 1440}) {
		t.Errorf("Reminders(10, 1440) = %v", got)
	}
	if got := Reminders(); got == nil || *got == nil || len(*got) != 0 {
		t.Errorf("Reminders() must be a non-nil empty slice, got %#v", got)
	}
	if got := NoReminders(); got == nil || *got == nil || len(*got) != 0 {
		t.Errorf("NoReminders() must be a non-nil empty slice, got %#v", got)
	}
	if got := InheritReminders(); got == nil || *got != nil {
		t.Errorf("InheritReminders() must point to a nil slice, got %#v", got)
	}

	// Reminders copies its arguments: changing the caller's slice afterwards
	// must not change the schedule that is sent.
	in := []int{10, 60}
	r := Reminders(in...)
	in[0] = 99
	if (*r)[0] != 10 {
		t.Errorf("Reminders must copy its input, got %v", *r)
	}
}

func TestReminderParamsMarshal(t *testing.T) {
	cases := []struct {
		name string
		val  *[]int
		want string // "" = field absent
	}{
		{"nil omits", nil, ""},
		{"empty sends []", NoReminders(), "[]"},
		{"zero-arg Reminders sends []", Reminders(), "[]"},
		{"values", Reminders(10, 1440), "[10,1440]"},
		{"inherit sends null", InheritReminders(), "null"},
	}

	params := []struct {
		name  string
		field string
		build func(*[]int) interface{}
	}{
		{"CreateCalendarParams", "default_reminders", func(v *[]int) interface{} {
			return &CreateCalendarParams{Name: "c", Timezone: "UTC", DefaultReminders: v}
		}},
		{"UpdateCalendarParams", "default_reminders", func(v *[]int) interface{} {
			return &UpdateCalendarParams{DefaultReminders: v}
		}},
		{"CreateEventParams", "reminders", func(v *[]int) interface{} {
			return &CreateEventParams{Title: "e", StartTime: "2026-01-01T09:00:00Z", EndTime: "2026-01-01T10:00:00Z", Reminders: v}
		}},
		{"UpdateEventParams", "reminders", func(v *[]int) interface{} {
			return &UpdateEventParams{Reminders: v}
		}},
	}

	for _, p := range params {
		for _, c := range cases {
			t.Run(p.name+"/"+c.name, func(t *testing.T) {
				b, err := json.Marshal(p.build(c.val))
				if err != nil {
					t.Fatal(err)
				}
				var m map[string]json.RawMessage
				if err := json.Unmarshal(b, &m); err != nil {
					t.Fatal(err)
				}
				raw, ok := m[p.field]
				if c.want == "" {
					if ok {
						t.Errorf("expected %s absent, got %s in %s", p.field, raw, b)
					}
					return
				}
				if !ok {
					t.Fatalf("expected %s=%s, field absent in %s", p.field, c.want, b)
				}
				if string(raw) != c.want {
					t.Errorf("expected %s=%s, got %s", p.field, c.want, raw)
				}
			})
		}
	}
}

func TestUpdateCalendarSendsEmptyDefaultReminders(t *testing.T) {
	var body string
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		b, _ := io.ReadAll(r.Body)
		body = string(b)
		if r.Method != http.MethodPatch || r.URL.Path != "/v1/calendars/cal_1" {
			t.Errorf("unexpected request %s %s", r.Method, r.URL.Path)
		}
		w.WriteHeader(200)
		w.Write([]byte(`{"id":"cal_1","name":"c","timezone":"UTC","default_reminders":[]}`))
	}))
	defer srv.Close()

	client, _ := NewClient(WithAPIKey("test_key"), WithBaseURL(srv.URL))
	cal, err := client.Calendars.Update(context.Background(), "cal_1", &UpdateCalendarParams{DefaultReminders: NoReminders()})
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(body, `"default_reminders":[]`) {
		t.Errorf("expected request body to carry default_reminders:[], got %s", body)
	}
	if cal.DefaultReminders == nil || len(cal.DefaultReminders) != 0 {
		t.Errorf("expected decoded empty DefaultReminders, got %#v", cal.DefaultReminders)
	}
}

func TestCreateEventSendsEmptyReminders(t *testing.T) {
	var body string
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		b, _ := io.ReadAll(r.Body)
		body = string(b)
		if r.Method != http.MethodPost || r.URL.Path != "/v1/calendars/cal_1/events" {
			t.Errorf("unexpected request %s %s", r.Method, r.URL.Path)
		}
		w.WriteHeader(201)
		w.Write([]byte(`{"id":"evt_1","calendarId":"cal_1","title":"e","startTime":"2026-01-01T09:00:00Z","endTime":"2026-01-01T10:00:00Z","reminders":[]}`))
	}))
	defer srv.Close()

	client, _ := NewClient(WithAPIKey("test_key"), WithBaseURL(srv.URL))
	evt, err := client.Events.Create(context.Background(), "cal_1", &CreateEventParams{
		Title:     "e",
		StartTime: "2026-01-01T09:00:00Z",
		EndTime:   "2026-01-01T10:00:00Z",
		Reminders: NoReminders(),
	})
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(body, `"reminders":[]`) {
		t.Errorf("expected request body to carry reminders:[], got %s", body)
	}
	if evt.Reminders == nil || len(evt.Reminders) != 0 {
		t.Errorf("expected decoded empty Reminders, got %#v", evt.Reminders)
	}
}

func TestUpdateEventSendsNullReminders(t *testing.T) {
	var body string
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		b, _ := io.ReadAll(r.Body)
		body = string(b)
		w.WriteHeader(200)
		w.Write([]byte(`{"id":"evt_1","calendarId":"cal_1","title":"e","startTime":"2026-01-01T09:00:00Z","endTime":"2026-01-01T10:00:00Z","reminders":null}`))
	}))
	defer srv.Close()

	client, _ := NewClient(WithAPIKey("test_key"), WithBaseURL(srv.URL))
	if _, err := client.Events.Update(context.Background(), "cal_1", "evt_1", &UpdateEventParams{Reminders: InheritReminders()}); err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(body, `"reminders":null`) {
		t.Errorf("expected request body to carry reminders:null, got %s", body)
	}
}
